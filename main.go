package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	username = flag.String("username", "", "Username for basic authentication")
	password = flag.String("password", "", "Password for basic authentication")
	port     = flag.Int("port", 8080, "http port for server")
)

func main() {
	flag.Parse()
	var handler, metricsHandler func(w http.ResponseWriter, r *http.Request)
	if *username != "" && *password != "" {
		log.Println("Started Http Server with basic auth enabled")
		handler = handleWithBasicAuth(serverHandler, *username, *password)
		metricsHandler = handleWithBasicAuth(promhttp.Handler().ServeHTTP, *username, *password)
	} else {
		log.Println("Started Http Server with basic auth disabled")
		handler = serverHandler
		metricsHandler = promhttp.Handler().ServeHTTP
	}
	http.HandleFunc("/metrics", metricsHandler)
	http.HandleFunc("/", handler)
	port := *port
	addr := fmt.Sprintf(":%d", port)
	fmt.Printf("Server started at http://localhost:%d\n", port)
	log.Fatal(http.ListenAndServe(addr, nil))
}

func serverHandler(w http.ResponseWriter, r *http.Request) {
	// Log request details
	log.Printf("Received %s request from %s with payload: %s\n", r.Method, r.RemoteAddr, extractPayload(r))

	// Dummy response
	response := "Hello, this is a dummy response!"
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(response))
}

func handleWithBasicAuth(handler http.HandlerFunc, username, password string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, pass, ok := r.BasicAuth()

		if !ok || user != username || pass != password {
			log.Printf("Unauthorized access attempt from %s with incorrect credentials\n", r.RemoteAddr)
			w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		handler(w, r)
	}
}

func extractPayload(r *http.Request) string {
	if r.Method == http.MethodPost {
		body := make([]byte, r.ContentLength)
		r.Body.Read(body)
		return string(body)
	}
	return ""
}
