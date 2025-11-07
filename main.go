package main

import (
	"log"
	"net/http"
	"sync/atomic"
	"fmt"
)

type apiConfig struct {
	fileserverHits atomic.Int32
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

func (cfg *apiConfig) handleFileServerHits(w http.ResponseWriter, r *http.Request) {
	log.Printf("Hits: %d\n", cfg.fileserverHits.Load())
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("Hits: %d", cfg.fileserverHits.Load())))

}


func main() {
	const port = "8080"
	const filepathRoot = "."

	apiCfg := apiConfig{
		fileserverHits: atomic.Int32{},
	}

	// Create the multiplexer
	mux := http.NewServeMux()
	mux.Handle("/app/", 
		apiCfg.middlewareMetricsInc(
			http.StripPrefix("/app", 
				http.FileServer(http.Dir(filepathRoot)),
			),
		),
	)
	mux.HandleFunc("/healthz", handleReadiness)
	// register the handler that logs the server hits on /metrics
	mux.HandleFunc("/metrics", apiCfg.handleFileServerHits)
	// register the handler that resets the counter to 0 on /reset
	mux.HandleFunc("/reset", apiCfg.handleFileServerHits)

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


