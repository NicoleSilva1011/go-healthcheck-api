package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"sync/atomic"
	"time"
)

var (
	startTime    = time.Now()
	requestCount uint64
)

func main() {
	http.HandleFunc("/health", withMetrics(healthHandler))
	http.HandleFunc("/ready", withMetrics(readyHandler))
	http.HandleFunc("/metrics", metricsHandler)

	log.Println("Server running on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func withMetrics(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		atomic.AddUint64(&requestCount, 1)
		next(w, r)
	}
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	resp := map[string]string{
		"status": "ok",
		"uptime": time.Since(startTime).String(),
	}
	json.NewEncoder(w).Encode(resp)
}

func readyHandler(w http.ResponseWriter, r *http.Request) {
	dbStatus := "ok"
	cacheStatus := "ok"

	if os.Getenv("DB_DOWN") == "true" {
		dbStatus = "down"
	}
	if os.Getenv("CACHE_DOWN") == "true" {
		cacheStatus = "down"
	}

	status := "ready"
	if dbStatus == "down" || cacheStatus == "down" {
		status = "not ready"
		w.WriteHeader(http.StatusServiceUnavailable)
	}

	resp := map[string]interface{}{
		"status": status,
		"dependencies": map[string]string{
			"database": dbStatus,
			"cache":    cacheStatus,
		},
	}
	json.NewEncoder(w).Encode(resp)
}

func metricsHandler(w http.ResponseWriter, r *http.Request) {
	resp := map[string]interface{}{
		"requests":       atomic.LoadUint64(&requestCount),
		"uptime_seconds": int(time.Since(startTime).Seconds()),
	}
	json.NewEncoder(w).Encode(resp)
}
