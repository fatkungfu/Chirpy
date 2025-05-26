package main

import (
	"fmt"
	"net/http"
)

// metricsHandler is a method on *apiConfig that writes the number of requests
// to the HTTP response in the format "Hits: x".
func (cfg *apiConfig) handlerMetrics(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	// Write the number of hits to the response
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("Hits: %d", cfg.fileserverHits.Load())))
}

// middlewareMetricsInc is a middleware method on *apiConfig
// that increments the fileserverHits counter every time it's called.
func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Increment the fileserverHits counter atomically
		cfg.fileserverHits.Add(1)
		// Call the next handler in the chain
		next.ServeHTTP(w, r)
	})
}
