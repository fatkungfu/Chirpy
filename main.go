package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/fatkungfu/Chirpy/internal/database"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

// apiConfig struct holds any stateful, in-memory data.
// fileserverHits uses atomic.Int32 for safe concurrent access.
type apiConfig struct {
	fileserverHits atomic.Int32 // Uncomment if you want to track hits
	db             *database.Queries
	platform       string
}

func main() {
	const filepathRoot = "."
	const port = "8080"

	godotenv.Load()              // Load environment variables from .env file
	dbURL := os.Getenv("DB_URL") // Get the database URL from environment variables
	if dbURL == "" {
		log.Fatal("DB_URL must be set")
	}
	// Get the platform from environment variables
	platform := os.Getenv("PLATFORM")
	if platform == "" {
		log.Fatal("PLATFORM must be set")
	}

	// Open a connection to the database using the provided URL.
	dbConn, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("error opening database: %s", err)
	}
	dbQueries := database.New(dbConn)

	// Create an instance of apiConfig to hold our state.
	apiCfg := apiConfig{
		fileserverHits: atomic.Int32{}, // Initialize the atomic counter
		db:             dbQueries,
		platform:       platform,
	}

	// This will act as our HTTP request router. Since we're not
	// registering any specific paths, it will default to a 404 for all requests.
	mux := http.NewServeMux()
	fsHandler := apiCfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot))))
	// Use middleware to wrap the file server handler
	mux.Handle("/app/", fsHandler)
	// Register the validate_chirp handler
	mux.HandleFunc("POST /api/validate_chirp", handlerValidateChirp)
	// Use mux.HandleFunc to register the healthzHandler for the /healthz path.
	mux.HandleFunc("GET /healthz", handlerReadiness)
	// Register the metrics handler
	mux.HandleFunc("GET /metrics", apiCfg.handlerMetrics)
	// Register the reset handler
	mux.HandleFunc("POST /api/users", apiCfg.handlerUsersCreate)
	mux.HandleFunc("POST /admin/reset", apiCfg.handlerReset)

	// We'll configure its address and assign our ServeMux as its handler.
	server := &http.Server{
		Addr:    ":" + port, // Listen on port 8080
		Handler: mux,        // Use our ServeMux for routing
	}

	// Use the server's ListenAndServe method to start the server
	// This will block until the server encounters an error or is shut down.
	log.Printf("Serving on port: %s\n", port)
	log.Fatal(server.ListenAndServe())
}
