package main

import (
	"log"
	"net/http"
	"sync/atomic"
)

type apiConfig struct {
	fileserverHits atomic.Int32
}



func main() {
	const port = "8080"
	const filepathRoot = "."

	apiCfg := apiConfig{
		fileserverHits: atomic.Int32{},
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
	// register the handler that logs the server hits on /metrics
	mux.HandleFunc("GET /api/metrics", apiCfg.handlerMetrics)
	// register the handler that resets the counter to 0 on /reset
	mux.HandleFunc("POST /api/reset", apiCfg.handlerReset)

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


