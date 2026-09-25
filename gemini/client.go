package gemini

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"sync"
	"time"

	"github.com/coder/websocket"
)

const (
	GeminiLiveHost = "generativelanguage.googleapis.com"
	GeminiLivePath = "/ws/google.ai.generativelanguage.v1beta.GenerativeService.BidiGenerateContent"
)

// Setup structures
type ClientSetupMessage struct {
	Setup SetupConfig `json:"setup"`
}

type SetupConfig struct {
	Model                    string                         `json:"model"`
	GenerationConfig         *GenerationConfig              `json:"generationConfig,omitempty"`
	SystemInstruction        *SystemInstruction             `json:"systemInstruction,omitempty"`
	InputAudioTranscription  *InputAudioTranscriptionConfig `json:"inputAudioTranscription,omitempty"`
	OutputAudioTranscription *AudioTranscriptionConfig      `json:"outputAudioTranscription,omitempty"`
	RealtimeInputConfig      *RealtimeInputConfig           `json:"realtimeInputConfig,omitempty"`
}

type GenerationConfig struct {
	ResponseModalities []string           `json:"responseModalities,omitempty"`
	TranslationConfig  *TranslationConfig `json:"translationConfig,omitempty"`
}

type TranslationConfig struct {
	TargetLanguageCode string `json:"targetLanguageCode"`
	EchoTargetLanguage bool   `json:"echoTargetLanguage,omitempty"`
}

type SystemInstruction struct {
	Parts []ContentPart `json:"parts"`
}

type ContentPart struct {
	Text string `json:"text,omitempty"`
}

// InputAudioTranscriptionConfig enables inputTranscription and interimInputTranscription
// events in server responses. Required for gemini-3.5-transcribe-live to stream
// real-time word-by-word transcriptions.
type InputAudioTranscriptionConfig struct {
	LanguageCodes []string `json:"languageCodes,omitempty"` // empty = auto-detect
}

type AudioTranscriptionConfig struct{}

type RealtimeInputConfig struct {
	ActivityHandling           string                      `json:"activityHandling,omitempty"`
	AutomaticActivityDetection *AutomaticActivityDetection `json:"automaticActivityDetection,omitempty"`
}

type AutomaticActivityDetection struct {
	Disabled                 bool   `json:"disabled"`
	StartOfSpeechSensitivity string `json:"startOfSpeechSensitivity,omitempty"`
	EndOfSpeechSensitivity   string `json:"endOfSpeechSensitivity,omitempty"`
	PrefixPaddingMs          int    `json:"prefixPaddingMs,omitempty"`
	SilenceDurationMs        int    `json:"silenceDurationMs,omitempty"`
}

// Media streaming structures
type RealtimeInputMessage struct {
	RealtimeInput RealtimeInputData `json:"realtimeInput"`
}

type RealtimeInputData struct {
	MediaChunks []MediaChunk `json:"mediaChunks"`
}

type MediaChunk struct {
	MimeType string `json:"mimeType"`
	Data     string `json:"data"` // base64 encoded
}

// Client content control messages (flush)
type ClientContentMessage struct {
	ClientContent ClientContentData `json:"clientContent"`
}

type ClientContentData struct {
	ActivityEnd   *ActivityEnd   `json:"activityEnd,omitempty"`
	ActivityStart *ActivityStart `json:"activityStart,omitempty"`
}

type ActivityEnd struct{}
type ActivityStart struct{}

// Server response structures
type ServerResponse struct {
	ServerContent *ServerContent `json:"serverContent,omitempty"`
}

type ServerContent struct {
	ModelTurn                 *ModelTurn      `json:"modelTurn,omitempty"`
	OutputTranscription       *InputTranscription `json:"outputTranscription,omitempty"`
	InputTranscription        *InputTranscription `json:"inputTranscription,omitempty"`
	InterimInputTranscription *InputTranscription `json:"interimInputTranscription,omitempty"`
	TurnComplete              bool            `json:"turnComplete,omitempty"`
}

type ModelTurn struct {
	Parts []ContentPart `json:"parts,omitempty"`
}

// InputTranscription is used for inputTranscription, interimInputTranscription,
// and outputTranscription server events.
type InputTranscription struct {
	Text string `json:"text,omitempty"`
}

// LiveSession represents a persistent bidirectional streaming connection to Gemini
type LiveSession struct {
	conn   *websocket.Conn
	ctx    context.Context
	cancel context.CancelFunc
	mu     sync.Mutex
}

func ConnectLive(ctx context.Context, apiKey string, setup SetupConfig) (*LiveSession, error) {
	u := url.URL{
		Scheme: "wss",
		Host:   GeminiLiveHost,
		Path:   GeminiLivePath,
	}
	q := u.Query()
	q.Set("key", apiKey)
	u.RawQuery = q.Encode()

	sessCtx, cancel := context.WithCancel(ctx)

	conn, resp, err := websocket.Dial(sessCtx, u.String(), &websocket.DialOptions{
		HTTPHeader: http.Header{
			"Content-Type": []string{"application/json"},
		},
	})
	if err != nil {
		cancel()
		if resp != nil {
			return nil, fmt.Errorf("gemini websocket dial failed (HTTP %d): %w", resp.StatusCode, err)
		}
		return nil, fmt.Errorf("gemini websocket dial failed: %w", err)
	}
	conn.SetReadLimit(16 * 1024 * 1024)

	// Send initial setup frame
	setupMsg := ClientSetupMessage{Setup: setup}
	data, err := json.Marshal(setupMsg)
	if err != nil {
		conn.Close(websocket.StatusInternalError, "setup marshal error")
		cancel()
		return nil, fmt.Errorf("failed to marshal setup: %w", err)
	}

	writeCtx, writeCancel := context.WithTimeout(sessCtx, 10*time.Second)
	defer writeCancel()

	if err := conn.Write(writeCtx, websocket.MessageText, data); err != nil {
		conn.Close(websocket.StatusInternalError, "setup send error")
		cancel()
		return nil, fmt.Errorf("failed to send setup frame: %w", err)
	}

	return &LiveSession{
		conn:   conn,
		ctx:    sessCtx,
		cancel: cancel,
	}, nil
}

func (s *LiveSession) SendAudio(pcmData []byte) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	encoded := base64.StdEncoding.EncodeToString(pcmData)
	msg := RealtimeInputMessage{
		RealtimeInput: RealtimeInputData{
			MediaChunks: []MediaChunk{
				{
					MimeType: "audio/pcm;rate=16000",
					Data:     encoded,
				},
			},
		},
	}

	data, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	writeCtx, cancel := context.WithTimeout(s.ctx, 5*time.Second)
	defer cancel()

	return s.conn.Write(writeCtx, websocket.MessageText, data)
}

func (s *LiveSession) SendActivityEnd() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	msg := ClientContentMessage{
		ClientContent: ClientContentData{
			ActivityEnd: &ActivityEnd{},
		},
	}
	data, _ := json.Marshal(msg)
	writeCtx, cancel := context.WithTimeout(s.ctx, 5*time.Second)
	defer cancel()
	return s.conn.Write(writeCtx, websocket.MessageText, data)
}

func (s *LiveSession) SendActivityStart() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	msg := ClientContentMessage{
		ClientContent: ClientContentData{
			ActivityStart: &ActivityStart{},
		},
	}
	data, _ := json.Marshal(msg)
	writeCtx, cancel := context.WithTimeout(s.ctx, 5*time.Second)
	defer cancel()
	return s.conn.Write(writeCtx, websocket.MessageText, data)
}

func (s *LiveSession) Receive() (*ServerResponse, error) {
	_, data, err := s.conn.Read(s.ctx)
	if err != nil {
		return nil, err
	}

	var resp ServerResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		log.Printf("Gemini unmarshal warning: %v (raw: %s)", err, string(data))
		return nil, err
	}
	return &resp, nil
}

func (s *LiveSession) Close() error {
	s.cancel()
	return s.conn.Close(websocket.StatusNormalClosure, "session closed")
}
