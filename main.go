package main

import (
	"encoding/json"
	"log"
	"net"
	"net/http"
	"time"
)

type Response struct {
	StatusCode int    `json:"status_code"`
	Message    string `json:"message"`
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
	// Extract client IP (accounting for possible proxies)
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		ip = r.RemoteAddr
	}

	// Log the request
	log.Printf("[%s] Request from %s %s %s",
		time.Now().Format(time.RFC3339), ip, r.Method, r.URL.Path)

	// Write response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	resp := Response{
		StatusCode: http.StatusOK,
		Message:    "Hello, World!",
	}
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func main() {
	http.HandleFunc("/", helloHandler)

	log.Println("Server is running on http://localhost:9877")
	if err := http.ListenAndServe(":9877", nil); err != nil {
		log.Fatal(err)
	}
}
