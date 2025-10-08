package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"log/slog"
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
	mux.HandleFunc("/health/check", healthCheckFunc)
	mux.HandleFunc("/health/alive", IsReady)

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	logger.Info("ServeGo is Starting", "Port", "8082")

	server := &http.Server{
		Addr:              ":8082",
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       15 * time.Second,
		WriteTimeout:      15 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	go func() {
		logger.Info("ServeGo is Running", "Port", server.Addr)
		// If the server stopped normally, don’t panic.
		// But if it failed unexpectedly, log the reason and exit.
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("Listen Error", "Error", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	sig := <-quit
	logger.Warn("Shutdown Signal Recieved", "Signal", sig.String())

	fmt.Printf("\n")
	logger.Info("Shutting down ServeGO")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Printf("Graceful shutdown Failed: %v", err)
	}
	log.Printf("Server Stopped")

}

func handleRoot(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	resp := Response{
		Status:  "success",
		Message: "Welcome to ServeGo",
		Time:    time.Now().Format(time.RFC3339),
	}
	json.NewEncoder(w).Encode(resp)

}

func healthCheckFunc(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	resp := Response{
		Status:  "Success",
		Message: "Instance to Healthy",
		Time:    time.Now().Format(time.RFC3339),
	}
	json.NewEncoder(w).Encode(resp)
}

func IsReady(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	resp := Response{
		Status:  "Success",
		Message: "Instance is Alive",
		Time:    time.Now().Format(time.RFC3339),
	}
	json.NewEncoder(w).Encode(resp)
}

type Response struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Time    string `json:"time"`
}
