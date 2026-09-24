package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"nerdearla-live-subs-go/server"

	"github.com/joho/godotenv"
)

func main() {
	// Try loading .env from current dir or sibling directory
	_ = godotenv.Load(".env")
	_ = godotenv.Load("../nerdearla-live-subs/.env")

	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		log.Println("WARNING: GEMINI_API_KEY is not set.")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	htmlPath := server.FindHTMLPath()
	log.Printf("Using index.html from: %s", htmlPath)

	srv := server.NewServer(apiKey, htmlPath)
	httpServer := &http.Server{
		Addr:         ":" + port,
		Handler:      srv.Routes(),
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 0, // Streaming responses (SSE / WebSockets) need unlimited write timeout
		IdleTimeout:  60 * time.Second,
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	go func() {
		log.Printf("🚀 Nerdearla Live Subs (Go Edition) running on http://localhost:%s", port)
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server listen error: %v", err)
		}
	}()

	<-stop
	log.Println("Shutting down server gracefully...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := httpServer.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exiting")
}
