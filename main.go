package main

import (
	"flag"
	"fmt"
	"log"
	httpserver "mock-http-server/http"
	"net/http"
)

var (
	port = flag.String("port", "8080", "http port for server")
)

func main() {
	flag.Parse()
	// Create a new mock server with the provided YAML configuration
	mockServer, err := httpserver.NewMockServer()
	if err != nil {
		log.Fatalf("Failed to create mock server: %v", err)
	}
	http.HandleFunc("/", mockServer.Handler)
	http.HandleFunc("/metrics", httpserver.MetricsHandler().ServeHTTP)
	// Start the mock server
	log.Printf("Mock server is running on port %s...", *port)
	ln, err := httpserver.GetListener(fmt.Sprintf("0.0.0.0:%s", *port))
	if err != nil {
		log.Fatal(err)
	}
	http.Serve(*ln, nil)
}
