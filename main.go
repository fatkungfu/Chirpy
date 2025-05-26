package main

import (
	"log"
	"net/http"
)

func main() {
	const filepathRoot = "."
	const port = "8080"

	// This will act as our HTTP request router. Since we're not
	// registering any specific paths, it will default to a 404 for all requests.
	mux := http.NewServeMux()
	mux.Handle("/app/", http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot))))
	// Use mux.HandleFunc to register the healthzHandler for the /healthz path.
	mux.HandleFunc("/healthz", healthzHandler)

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

// healthzHandler is a handler for the /healthz endpoint.
// It returns a 200 OK status with "OK" in the body,
// indicating that the server is ready to receive traffic.
func healthzHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)                    // Set the status code to 200 OK
	w.Write([]byte(http.StatusText(http.StatusOK))) // Write "OK" to the response body
}
