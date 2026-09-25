package server

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
)

func TestServerEndpoints(t *testing.T) {
	wd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get working directory: %v", err)
	}
	rootDir := filepath.Dir(wd)
	htmlPath := filepath.Join(rootDir, "index.html")

	s := NewServer("test-api-key", htmlPath)
	handler := s.Routes()

	// 1. Health check
	t.Run("HealthCheck", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/health", nil)
		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d", rec.Code)
		}

		var body map[string]interface{}
		if err := json.NewDecoder(rec.Body).Decode(&body); err != nil {
			t.Fatalf("Failed to decode health response: %v", err)
		}
		if body["status"] != "ok" {
			t.Errorf("Expected status ok, got %v", body["status"])
		}
	})

	// 2. Samples list
	t.Run("ListSamples", func(t *testing.T) {
		// Change working directory to root dir for the samples endpoint to locate "samples"
		origWd, _ := os.Getwd()
		_ = os.Chdir(rootDir)
		defer func() { _ = os.Chdir(origWd) }()

		req := httptest.NewRequest("GET", "/api/samples", nil)
		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d", rec.Code)
		}

		var samples []map[string]string
		if err := json.NewDecoder(rec.Body).Decode(&samples); err != nil {
			t.Fatalf("Failed to decode samples: %v", err)
		}
		if len(samples) < 2 {
			t.Errorf("Expected at least 2 samples, got %d", len(samples))
		}
	})

	// 3. Static Root (index.html)
	t.Run("RootHtml", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/", nil)
		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d", rec.Code)
		}
	})
}
