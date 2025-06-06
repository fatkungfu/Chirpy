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
	jwtSecret      string
	polkaKey       string
}

func main() {
	const filepathRoot = "."
	const port = "8080"

	godotenv.Load(".env")        // Load environment variables from .env file
	dbURL := os.Getenv("DB_URL") // Get the database URL from environment variables
	if dbURL == "" {
		log.Fatal("DB_URL must be set")
	}
	// Get the platform from environment variables
	platform := os.Getenv("PLATFORM")
	if platform == "" {
		log.Fatal("ADMIN_KEY environment variable is not set")
	}
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Fatal("JWT_SECRET environment variable is not set")
	}
	polkaKey := os.Getenv("POLKA_KEY")
	if polkaKey == "" {
		log.Fatal("POLKA_KEY environment variable must be set")
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
		jwtSecret:      jwtSecret,
		polkaKey:       polkaKey,
	}

	// This will act as our HTTP request router. Since we're not
	// registering any specific paths, it will default to a 404 for all requests.
	mux := http.NewServeMux()
	fsHandler := apiCfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot))))
	// Use middleware to wrap the file server handler
	mux.Handle("/app/", fsHandler)

	// Use mux.HandleFunc to register the healthzHandler for the /healthz path.
	mux.HandleFunc("GET /healthz", handlerReadiness)

	mux.HandleFunc("POST /api/polka/webhooks", apiCfg.handlerWebhook)

	mux.HandleFunc("POST /api/login", apiCfg.handlerLogin)
	mux.HandleFunc("POST /api/refresh", apiCfg.handlerRefresh)
	mux.HandleFunc("POST /api/revoke", apiCfg.handlerRevoke)

	// Register the users create handler
	mux.HandleFunc("POST /api/users", apiCfg.handlerUsersCreate)
	mux.HandleFunc("PUT /api/users", apiCfg.handlerUsersUpdate)

	// Register the chirps create handler
	mux.HandleFunc("POST /api/chirps", apiCfg.handlerChirpsCreate)
	// Register the chirps get by many ID handler
	mux.HandleFunc("GET /api/chirps", apiCfg.handlerChirpsRetrieve)
	mux.HandleFunc("GET /api/chirps/{chirpID}", apiCfg.handlerChirpsGet)
	mux.HandleFunc("DELETE /api/chirps/{chirpID}", apiCfg.handlerChirpsDelete)

	// Register the reset handler
	mux.HandleFunc("POST /admin/reset", apiCfg.handlerReset)
	// Register the metrics handler
	mux.HandleFunc("GET /metrics", apiCfg.handlerMetrics)

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
