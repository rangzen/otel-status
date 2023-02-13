// Simple HTTP server that responds with different status codes and delays.
package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

func main() {
	fmt.Println("Starting HTTP servers...")
	fmt.Println("8080: HTML normal response")
	fmt.Println("8081: JSON slow response")
	fmt.Println("8082: 401 no body")
	start()
}

// start starts 3 http servers:
// * one that respond normally (HTML),
// * one that responds in 1 second (JSON),
// * and one that always responds 401 (no body).
func start() {
	errChan := make(chan error)
	go func() {
		errChan <- http.ListenAndServe(":8080", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "text/html")
			_, err := w.Write([]byte("<h1>Hello World</h1>"))
			if err != nil {
				errChan <- fmt.Errorf("writing response: %w", err)
				return
			}
		}))
	}()

	go func() {
		errChan <- http.ListenAndServe(":8081", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			time.Sleep(1 * time.Second)
			w.Header().Set("Content-Type", "application/json")
			_, err := w.Write([]byte(`{"message": "slow response"}`))
			if err != nil {
				errChan <- fmt.Errorf("writing response: %w", err)
				return
			}
		}))
	}()

	go func() {
		errChan <- http.ListenAndServe(":8082", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(401)
		}))
	}()

	log.Fatal(<-errChan)
}
