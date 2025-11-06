package main

import (
	"log"
	"net/http"
)

func main() {
	const port = "8080"
	const root = "/"
	const filepathRoot = "."

	// Create the multiplexer
	mux := http.NewServeMux()
	mux.Handle(root, http.FileServer(http.Dir(filepathRoot)))

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
