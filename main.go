package main

import (
	"log"
	"net/http"
)

func main() {
	addr := ":8080"
	baseURL := "http://localhost:8080"

	store := NewURLStore()
	handler := NewURLHandler(store, baseURL)

	mux := http.NewServeMux()
	mux.HandleFunc("/shorten", handler.Shorten)
	mux.HandleFunc("/", handler.Redirect)

	log.Printf("server starting on %s", addr)
	if err := http.ListenAndServe(addr, loggingMiddleware(mux)); err != nil {
		log.Fatal(err)
	}
}
