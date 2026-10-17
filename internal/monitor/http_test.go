package monitor

import (
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"smartx-monitor/internal/signals"
)

func TestHTTPContract(t *testing.T) {
	service := NewService(slog.New(slog.NewTextHandler(io.Discard, nil)), 1)
	service.SetAnalyzer(signals.NewAnalyzer(1000))
	service.Publish(Event{ID: "e1", Chain: ChainSolana, Asset: "TOKEN", Type: EventTokenLaunch, ObservedAt: time.Now()})
	handler := service.Handler()

	for _, path := range []string{"/healthz", "/readyz", "/v1/signals", "/metrics"} {
		recording := httptest.NewRecorder()
		handler.ServeHTTP(recording, httptest.NewRequest(http.MethodGet, path, nil))
		if recording.Code != http.StatusOK {
			t.Fatalf("%s returned %d", path, recording.Code)
		}
	}
	metrics := httptest.NewRecorder()
	handler.ServeHTTP(metrics, httptest.NewRequest(http.MethodGet, "/metrics", nil))
	if !strings.Contains(metrics.Body.String(), "smartx_events_processed_total") {
		t.Fatal("metrics response missing processed counter")
	}
}

func TestSignalsAPIKey(t *testing.T) {
	service := NewService(slog.Default(), 1)
	service.SetAPIKey("secret")
	request := httptest.NewRequest(http.MethodGet, "/v1/signals", nil)
	response := httptest.NewRecorder()
	service.Handler().ServeHTTP(response, request)
	if response.Code != http.StatusUnauthorized {
		t.Fatalf("expected unauthorized, got %d", response.Code)
	}
}
