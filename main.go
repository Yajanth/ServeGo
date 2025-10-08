package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {

	mux := http.NewServeMux()

	// Routes
	mux.HandleFunc("/root", handleRoot)
	mux.HandleFunc("/health", healthCheckFunc)
	mux.HandleFunc("/health/alive", IsReady)

	fmt.Println("ServeGo is starting.....")

	server := &http.Server{
		Addr:              ":8082",
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       15 * time.Second,
		WriteTimeout:      15 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	go func() {
		log.Printf("ServeGo is listening on %s ", server.Addr)
		// If the server stopped normally, don’t panic.
		// But if it failed unexpectedly, log the reason and exit.
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Listen Error: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down ServeGO...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Printf("Graceful shutdown Failed: %v", err)
	}
	log.Printf("Server Stopped")

}

func handleRoot(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to ServeGO...")
}

func healthCheckFunc(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Instance is Healthy...")
}

func IsReady(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Instance is ready to use")
}

type Post struct {
	ID   int    `json:"id"`
	body string `json:"body"`
}
