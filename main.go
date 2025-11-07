package main

import (
	"log"
	"net/http"
)

func main() {
	const port = "8080"
	const root = "/"
	const filepathRoot = "."
	const readinessEndpoint = "/healthz"

	// Create the multiplexer
	mux := http.NewServeMux()
	mux.Handle("/app/", http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot))))
	mux.HandleFunc(readinessEndpoint, handleReadiness)

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

func handleReadiness(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}
