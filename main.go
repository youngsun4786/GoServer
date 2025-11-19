package main

import (
	"log"
	"net/http"
	"sync/atomic"
	"goserver/internal/database"
	"database/sql"
	"os"
	_ "github.com/lib/pq" // PostgreSQL driver
	_ "github.com/joho/godotenv"
)

type apiConfig struct {
	fileserverHits atomic.Int32
	db *database.Queries
}



func main() {
	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("DB_URL must be set")
	}

	dbConn, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal("Error opening database: %s", err)
	}
	dbQueries := database.New(dbConn)

	const port = "8080"
	const filepathRoot = "."

	apiCfg := apiConfig{
		fileserverHits: atomic.Int32{},
		db: dbQueries,
	}

	// Create the multiplexer
	mux := http.NewServeMux()
	fsHandler := apiCfg.middlewareMetricsInc(
			http.StripPrefix("/app", 
				http.FileServer(http.Dir(filepathRoot)),
			),
		)
	mux.Handle("/app/", fsHandler)
	mux.HandleFunc("GET /api/healthz", handleReadiness)
	mux.HandleFunc("POST /api/validate_chirp", handlerValidate)
	// register the handler that logs the server hits on /metrics
	mux.HandleFunc("GET /admin/metrics", apiCfg.handlerMetrics)
	// register the handler that resets the counter to 0 on /reset
	mux.HandleFunc("POST /admin/reset", apiCfg.handlerReset)

	// Create the server with configuration
	server := &http.Server {
		Addr: ":" + port,
		Handler: mux,
	}

	// Start the server
	log.Printf("Serving files from %s on port %s\n", filepathRoot, port)
	// ListenAndServe returns error
	log.Fatal(server.ListenAndServe())
}


