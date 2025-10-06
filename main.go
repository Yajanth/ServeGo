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

	mux.HandleFunc("/root", handleRoot)
	mux.HandleFunc("/health", healthCheckFunc)

	fmt.Println("ServeGo is starting.....")

	server := &http.Server{
		Addr:    ":8082",
		Handler: mux}

	go func() {
		log.Printf("ServeGo is listening on %s", server.Addr)
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

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	server.Shutdown(ctx)

}

func handleRoot(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "This is root HandlerFunction")
}

func healthCheckFunc(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Instance is Healthy...")
}
