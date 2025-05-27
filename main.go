package main

import (
	"log"
	"net/http"
	"sync/atomic"
)

// apiConfig struct holds any stateful, in-memory data.
// fileserverHits uses atomic.Int32 for safe concurrent access.
type apiConfig struct {
	fileserverHits atomic.Int32 // Uncomment if you want to track hits
}

func main() {
	const filepathRoot = "."
	const port = "8080"

	// Create an instance of apiConfig to hold our state.
	apiCfg := &apiConfig{
		fileserverHits: atomic.Int32{}, // Initialize the atomic counter
	}

	// This will act as our HTTP request router. Since we're not
	// registering any specific paths, it will default to a 404 for all requests.
	mux := http.NewServeMux()
	fsHandler := apiCfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot))))
	mux.Handle("/app/", fsHandler)                                   // Use middleware to wrap the file server handler
	mux.HandleFunc("POST /api/validate_chirp", handlerValidateChirp) // Register the validate_chirp handler
	mux.HandleFunc("GET /healthz", handlerReadiness)                 // Use mux.HandleFunc to register the healthzHandler for the /healthz path.
	mux.HandleFunc("GET /metrics", apiCfg.handlerMetrics)            // Register the metrics handler
	mux.HandleFunc("POST /reset", apiCfg.handlerReset)               // Register the reset handler

	// We'll configure its address and assign our ServeMux as its handler.
	server := &http.Server{
		Addr:    ":" + port, // Listen on port 8080
		Handler: mux,        // Use our ServeMux for routing
	}

	// Use the server's ListenAndServe method to start the server
	// This will block until the server encounters an error or is shut down.
	log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
	log.Fatal(server.ListenAndServe())
}
