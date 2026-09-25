package server

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"nerdearla-live-subs-go/broadcast"
	"nerdearla-live-subs-go/gemini"

	"github.com/coder/websocket"
)

const (
	ModelTranscribe = "gemini-3.5-transcribe-live"
	ModelTranslate  = "gemini-3.5-live-translate-preview"
)

var langNames = map[string]string{
	"en": "English",
	"es": "Spanish",
	"pt": "Portuguese",
}

type Server struct {
	apiKey      string
	htmlPath    string
	transcribeM string
	translateM  string
}

func NewServer(apiKey, htmlPath string) *Server {
	return &Server{
		apiKey:      apiKey,
		htmlPath:    htmlPath,
		transcribeM: ModelTranscribe,
		translateM:  ModelTranslate,
	}
}

func (s *Server) Routes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/", s.handleRoot)
	mux.HandleFunc("/api/health", s.handleHealth)
	mux.HandleFunc("/api/sessions", s.handleSessions)
	mux.HandleFunc("/api/samples", s.handleListSamples)
	mux.HandleFunc("/stream/", s.handleSSE)
	mux.HandleFunc("/ws/stream/", s.handleWSStream)
	mux.HandleFunc("/ws/audio/", s.handleWSAudio)
	mux.Handle("/samples/", http.StripPrefix("/samples/", http.FileServer(http.Dir("samples"))))
	return mux
}

func (s *Server) handleRoot(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	http.ServeFile(w, r, s.htmlPath)
}

func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	resp := map[string]interface{}{
		"status": "ok",
		"models": map[string]string{
			"transcribe": s.transcribeM,
			"translate":  s.translateM,
		},
		"api_key_set": s.apiKey != "",
	}
	_ = json.NewEncoder(w).Encode(resp)
}

func (s *Server) handleSessions(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	resp := map[string]interface{}{
		"sessions": broadcast.GetSessionStats(),
	}
	_ = json.NewEncoder(w).Encode(resp)
}

func (s *Server) handleListSamples(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	files, err := os.ReadDir("samples")
	if err != nil {
		_ = json.NewEncoder(w).Encode([]map[string]string{})
		return
	}
	var samples []map[string]string
	for _, f := range files {
		if !f.IsDir() && strings.HasSuffix(strings.ToLower(f.Name()), ".wav") {
			samples = append(samples, map[string]string{
				"filename": f.Name(),
				"url":      "/samples/" + f.Name(),
			})
		}
	}
	_ = json.NewEncoder(w).Encode(samples)
}

// handleSSE handles GET /stream/{session_id}
func (s *Server) handleSSE(w http.ResponseWriter, r *http.Request) {
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming unsupported", http.StatusInternalServerError)
		return
	}

	sessionID := strings.TrimPrefix(r.URL.Path, "/stream/")
	if decoded, err := url.PathUnescape(sessionID); err == nil && decoded != "" {
		sessionID = decoded
	}
	sessionID = strings.TrimSpace(sessionID)
	if sessionID == "" {
		http.Error(w, "Missing session ID", http.StatusBadRequest)
		return
	}

	broadcaster := broadcast.GetOrCreateBroadcaster(sessionID)
	eventCh := broadcaster.Subscribe()
	defer func() {
		broadcaster.Unsubscribe(eventCh)
		log.Printf("[%s] SSE disconnected (listeners: %d)", sessionID, broadcaster.ListenerCount())
	}()

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	flusher.Flush()

	log.Printf("[%s] SSE connected (listeners: %d)", sessionID, broadcaster.ListenerCount())

	heartbeat := time.NewTicker(15 * time.Second)
	defer heartbeat.Stop()

	for {
		select {
		case <-r.Context().Done():
			return
		case <-heartbeat.C:
			fmt.Fprintf(w, ": keepalive\n\n")
			flusher.Flush()
		case evt, ok := <-eventCh:
			if !ok {
				return
			}
			if evt.Type == "__heartbeat__" {
				fmt.Fprintf(w, ": keepalive\n\n")
			} else {
				fmt.Fprintf(w, "event: %s\ndata: %s\n\n", evt.Type, evt.Text)
			}
			flusher.Flush()
		}
	}
}

// handleWSStream handles WS /ws/stream/{session_id}
func (s *Server) handleWSStream(w http.ResponseWriter, r *http.Request) {
	sessionID := strings.TrimPrefix(r.URL.Path, "/ws/stream/")
	if decoded, err := url.PathUnescape(sessionID); err == nil && decoded != "" {
		sessionID = decoded
	}
	sessionID = strings.TrimSpace(sessionID)
	if sessionID == "" {
		http.Error(w, "Missing session ID", http.StatusBadRequest)
		return
	}

	conn, err := websocket.Accept(w, r, &websocket.AcceptOptions{
		OriginPatterns: []string{"*"},
	})
	if err != nil {
		log.Printf("WS Stream accept error: %v", err)
		return
	}
	defer conn.Close(websocket.StatusNormalClosure, "done")

	broadcaster := broadcast.GetOrCreateBroadcaster(sessionID)
	eventCh := broadcaster.Subscribe()
	defer func() {
		broadcaster.Unsubscribe(eventCh)
		log.Printf("[%s] WS stream disconnected (listeners: %d)", sessionID, broadcaster.ListenerCount())
	}()

	log.Printf("[%s] WS stream connected (listeners: %d)", sessionID, broadcaster.ListenerCount())

	for {
		select {
		case <-r.Context().Done():
			return
		case evt, ok := <-eventCh:
			if !ok {
				return
			}
			if evt.Type == "__heartbeat__" {
				continue
			}
			data, _ := json.Marshal(evt)
			writeCtx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
			err := conn.Write(writeCtx, websocket.MessageText, data)
			cancel()
			if err != nil {
				return
			}
		}
	}
}

// Audio config message from browser
type AudioConfigMessage struct {
	Source string `json:"source"`
	Target string `json:"target"`
}

// handleWSAudio handles WS /ws/audio/{session_id}
func (s *Server) handleWSAudio(w http.ResponseWriter, r *http.Request) {
	sessionID := strings.TrimPrefix(r.URL.Path, "/ws/audio/")
	if decoded, err := url.PathUnescape(sessionID); err == nil && decoded != "" {
		sessionID = decoded
	}
	sessionID = strings.TrimSpace(sessionID)
	if sessionID == "" {
		http.Error(w, "Missing session ID", http.StatusBadRequest)
		return
	}

	conn, err := websocket.Accept(w, r, &websocket.AcceptOptions{
		OriginPatterns: []string{"*"},
	})
	if err != nil {
		log.Printf("WS Audio accept error: %v", err)
		return
	}
	defer conn.Close(websocket.StatusNormalClosure, "done")

	log.Printf("[%s] Audio WebSocket connected", sessionID)
	broadcaster := broadcast.GetOrCreateBroadcaster(sessionID)
	defer func() {
		if broadcaster.ListenerCount() == 0 {
			broadcast.CleanupSession(sessionID)
		}
	}()

	// 1. Read configuration message
	readCtx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	msgType, data, err := conn.Read(readCtx)
	cancel()
	if err != nil {
		log.Printf("[%s] Failed to read audio config: %v", sessionID, err)
		return
	}

	sourceLang := "es"
	targetLang := "en"
	if msgType == websocket.MessageText {
		var cfg AudioConfigMessage
		if err := json.Unmarshal(data, &cfg); err == nil {
			if cfg.Source != "" {
				sourceLang = cfg.Source
			}
			if cfg.Target != "" {
				targetLang = cfg.Target
			}
		}
	}
	log.Printf("[%s] Language: %s -> %s", sessionID, sourceLang, targetLang)

	// Keep reconnect loop like Python
	reconnectCount := 0
	for {
		if r.Context().Err() != nil {
			break
		}
		shouldReconnect := s.runDualSession(r.Context(), conn, sessionID, sourceLang, targetLang, broadcaster)
		if !shouldReconnect || reconnectCount >= 5 {
			break
		}
		reconnectCount++
		log.Printf("[%s] Reconnecting (#%d)...", sessionID, reconnectCount)
		broadcaster.Broadcast("status", fmt.Sprintf("🔄 Reconnecting... (#%d)", reconnectCount))
		time.Sleep(1 * time.Second)
	}
}

func (s *Server) runDualSession(
	ctx context.Context,
	wsConn *websocket.Conn,
	sessionID, sourceLang, targetLang string,
	broadcaster *broadcast.SessionBroadcaster,
) bool {
	if s.apiKey == "" {
		broadcaster.Broadcast("status", "GEMINI_API_KEY is not set.")
		log.Printf("[%s] GEMINI_API_KEY not set", sessionID)
		return false
	}

	sessCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	// 1. Make configurations
	var langCodes []string
	if sourceLang != "" {
		langCodes = []string{sourceLang}
	}
	transcribeCfg := gemini.SetupConfig{
		Model: "models/" + s.transcribeM,
		GenerationConfig: &gemini.GenerationConfig{
			ResponseModalities: []string{"TEXT"},
		},
		InputAudioTranscription: &gemini.InputAudioTranscriptionConfig{
			LanguageCodes: langCodes,
		},
		RealtimeInputConfig: &gemini.RealtimeInputConfig{
			ActivityHandling: "NO_INTERRUPTION",
			AutomaticActivityDetection: &gemini.AutomaticActivityDetection{
				Disabled:                 false,
				StartOfSpeechSensitivity: "START_SENSITIVITY_HIGH",
				EndOfSpeechSensitivity:   "END_SENSITIVITY_LOW",
				PrefixPaddingMs:          100,
				SilenceDurationMs:        800,
			},
		},
	}

	targetCode := targetLang
	if targetCode == "" {
		targetCode = "es"
	}

	// Translation model requires TranslationConfig to translate into target language (e.g. "pt", "es", "en")
	// and responseModalities=["AUDIO"] with outputAudioTranscription.
	translateCfg := gemini.SetupConfig{
		Model: "models/" + s.translateM,
		GenerationConfig: &gemini.GenerationConfig{
			ResponseModalities: []string{"AUDIO"},
			TranslationConfig: &gemini.TranslationConfig{
				TargetLanguageCode: targetCode,
				EchoTargetLanguage: true,
			},
		},
		OutputAudioTranscription: &gemini.AudioTranscriptionConfig{},
		RealtimeInputConfig: &gemini.RealtimeInputConfig{
			ActivityHandling: "NO_INTERRUPTION",
		},
	}

	// 2. Connect to both sessions concurrently
	type connectResult struct {
		session *gemini.LiveSession
		err     error
	}
	tCh := make(chan connectResult, 1)
	xCh := make(chan connectResult, 1)

	go func() {
		s, err := gemini.ConnectLive(sessCtx, s.apiKey, transcribeCfg)
		tCh <- connectResult{session: s, err: err}
	}()

	go func() {
		s, err := gemini.ConnectLive(sessCtx, s.apiKey, translateCfg)
		xCh <- connectResult{session: s, err: err}
	}()

	tRes := <-tCh
	xRes := <-xCh

	if tRes.err != nil || xRes.err != nil {
		if tRes.session != nil {
			tRes.session.Close()
		}
		if xRes.session != nil {
			xRes.session.Close()
		}
		log.Printf("[%s] Gemini connection error: transcribe err=%v, translate err=%v", sessionID, tRes.err, xRes.err)
		if ctx.Err() != nil {
			return false
		}
		return true // retry
	}

	transcribeSess := tRes.session
	translateSess := xRes.session
	defer transcribeSess.Close()
	defer translateSess.Close()

	log.Printf("[%s] Both Gemini sessions established (%s + %s)", sessionID, s.transcribeM, s.translateM)
	broadcaster.Broadcast("status", fmt.Sprintf("Connected — translating to %s...", strings.ToUpper(targetCode)))

	// 3. Receive loops
	go s.transcribeReceiveLoop(sessCtx, cancel, transcribeSess, broadcaster, sessionID)
	go s.translateReceiveLoop(sessCtx, cancel, translateSess, broadcaster, sessionID)

	// 4. Sequential FIFO audio workers (guarantees in-order PCM delivery to Gemini)
	transcribeAudioCh := make(chan []byte, 200)
	translateAudioCh := make(chan []byte, 200)

	go func() {
		for {
			select {
			case <-sessCtx.Done():
				return
			case chunk, ok := <-transcribeAudioCh:
				if !ok {
					return
				}
				if err := transcribeSess.SendAudio(chunk); err != nil {
					log.Printf("[%s][transcribe] SendAudio error: %v", sessionID, err)
					cancel()
					return
				}
			}
		}
	}()

	go func() {
		for {
			select {
			case <-sessCtx.Done():
				return
			case chunk, ok := <-translateAudioCh:
				if !ok {
					return
				}
				if err := translateSess.SendAudio(chunk); err != nil {
					log.Printf("[%s][translate] SendAudio error: %v", sessionID, err)
					cancel()
					return
				}
			}
		}
	}()

	// 5. Ingest loop: Read PCM from audio WebSocket and forward sequentially to both workers
	for {
		msgType, pcmData, err := wsConn.Read(sessCtx)
		if err != nil {
			if errors.Is(err, context.Canceled) || websocket.CloseStatus(err) == websocket.StatusNormalClosure {
				log.Printf("[%s] Audio WebSocket disconnected cleanly", sessionID)
				return false
			}
			log.Printf("[%s] Audio WebSocket read error: %v", sessionID, err)
			return false
		}

		if msgType != websocket.MessageBinary || len(pcmData) == 0 {
			continue
		}

		select {
		case transcribeAudioCh <- pcmData:
		default:
		}

		select {
		case translateAudioCh <- pcmData:
		default:
		}
	}
}

func (s *Server) transcribeReceiveLoop(
	ctx context.Context,
	cancel context.CancelFunc,
	session *gemini.LiveSession,
	broadcaster *broadcast.SessionBroadcaster,
	sessionID string,
) {
	defer cancel()
	for {
		select {
		case <-ctx.Done():
			return
		default:
		}

		resp, err := session.Receive()
		if err != nil {
			log.Printf("[%s][transcribe] Receive ended: %v", sessionID, err)
			return
		}

		sc := resp.ServerContent
		if sc == nil {
			continue
		}

		// Interim: low-latency speculative partial hypothesis updated while speaker talks.
		// Sent as "transcription_interim" so the UI can REPLACE (not append) the current line.
		if sc.InterimInputTranscription != nil && sc.InterimInputTranscription.Text != "" {
			log.Printf("[%s][transcribe] ⟳ interim: %q", sessionID, sc.InterimInputTranscription.Text)
			broadcaster.Broadcast("transcription_interim", sc.InterimInputTranscription.Text)
		}

		// Final: authoritative committed transcript emitted when speech is finalized.
		// Sent as "transcription" — UI commits this text permanently.
		if sc.InputTranscription != nil && sc.InputTranscription.Text != "" {
			log.Printf("[%s][transcribe] ✅ final: %q", sessionID, sc.InputTranscription.Text)
			broadcaster.Broadcast("transcription", sc.InputTranscription.Text)
		}

		// Fallback: model_turn text parts (only present when inputAudioTranscription is NOT set)
		// Kept here as a safety net but should no longer fire with the correct config.
		if sc.ModelTurn != nil {
			for _, part := range sc.ModelTurn.Parts {
				if part.Text != "" {
					log.Printf("[%s][transcribe] ⚠ model_turn fallback: %q", sessionID, part.Text)
					broadcaster.Broadcast("transcription", part.Text)
				}
			}
		}

		if sc.TurnComplete {
			log.Printf("[%s][transcribe] turn_complete", sessionID)
			broadcaster.Broadcast("transcription_done", "")
		}
	}
}

func (s *Server) translateReceiveLoop(
	ctx context.Context,
	cancel context.CancelFunc,
	session *gemini.LiveSession,
	broadcaster *broadcast.SessionBroadcaster,
	sessionID string,
) {
	defer cancel()
	for {
		select {
		case <-ctx.Done():
			return
		default:
		}

		resp, err := session.Receive()
		if err != nil {
			log.Printf("[%s][translate] Receive ended: %v", sessionID, err)
			return
		}

		sc := resp.ServerContent
		if sc == nil {
			continue
		}

		if sc.OutputTranscription != nil && sc.OutputTranscription.Text != "" {
			log.Printf("[%s][translate] ✅ output_transcription: %q", sessionID, sc.OutputTranscription.Text)
			broadcaster.Broadcast("translation", sc.OutputTranscription.Text)
		}

		if sc.ModelTurn != nil {
			for _, part := range sc.ModelTurn.Parts {
				if part.Text != "" {
					log.Printf("[%s][translate] ✅ text part: %q", sessionID, part.Text)
					broadcaster.Broadcast("translation", part.Text)
				}
			}
		}

		if sc.TurnComplete {
			log.Printf("[%s][translate] turn_complete", sessionID)
			broadcaster.Broadcast("translation_done", "")
		}
	}
}

// FindHTMLPath searches for index.html in standard relative locations
func FindHTMLPath() string {
	candidates := []string{
		"index.html",
		"../nerdearla-live-subs/index.html",
		"nerdearla-live-subs/index.html",
		filepath.Join(os.Getenv("HOME"), "Documents/Nerdearla Hackathon/nerdearla-live-subs/index.html"),
	}
	for _, c := range candidates {
		if _, err := os.Stat(c); err == nil {
			return c
		}
	}
	return "index.html"
}
