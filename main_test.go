package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHealthHandler(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	w := httptest.NewRecorder()

	healthHandler(w, req)

	res := w.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		t.Fatalf("expected status 200, got %d", res.StatusCode)
	}

	var body map[string]string
	if err := json.NewDecoder(res.Body).Decode(&body); err != nil {
		t.Fatalf("failed to decode response body: %v", err)
	}

	if body["status"] != "ok" {
		t.Fatalf("expected status 'ok', got '%s'", body["status"])
	}
}

func TestReadyHandler_Ready(t *testing.T) {
	t.Setenv("DB_DOWN", "false")
	t.Setenv("CACHE_DOWN", "false")

	req := httptest.NewRequest(http.MethodGet, "/ready", nil)
	w := httptest.NewRecorder()

	readyHandler(w, req)

	res := w.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		t.Fatalf("expected status 200, got %d", res.StatusCode)
	}

	var body map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&body); err != nil {
		t.Fatalf("failed to decode response body: %v", err)
	}

	if body["status"] != "ready" {
		t.Fatalf("expected status 'ready', got '%v'", body["status"])
	}
}

func TestReadyHandler_NotReady_DBDown(t *testing.T) {
	t.Setenv("DB_DOWN", "true")
	t.Setenv("CACHE_DOWN", "false")

	req := httptest.NewRequest(http.MethodGet, "/ready", nil)
	w := httptest.NewRecorder()

	readyHandler(w, req)

	res := w.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusServiceUnavailable {
		t.Fatalf("expected status 503, got %d", res.StatusCode)
	}
}

func TestMetricsHandler(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/metrics", nil)
	w := httptest.NewRecorder()

	metricsHandler(w, req)

	res := w.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		t.Fatalf("expected status 200, got %d", res.StatusCode)
	}

	var body map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&body); err != nil {
		t.Fatalf("failed to decode response body: %v", err)
	}

	if _, ok := body["requests"]; !ok {
		t.Fatalf("expected 'requests' field in metrics response")
	}

	if _, ok := body["uptime_seconds"]; !ok {
		t.Fatalf("expected 'uptime_seconds' field in metrics response")
	}
}
