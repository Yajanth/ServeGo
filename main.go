package main

import (
	"context"
	"encoding/json"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"sync/atomic"
	"syscall"
	"time"

	"github.com/google/uuid"
)

type Response struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	TraceId string `json:"traceId"`
	Path    string `json:"path"`
	TS      string `json:"time"`
}

var (
	alive atomic.Bool
	ready atomic.Bool
)

func main() {

	mux := http.NewServeMux()
	alive.Store(true)
	ready.Store(false)

	// Routes
	mux.HandleFunc("/root", handleRoot)
	mux.HandleFunc("/health/check", handleCheck)
	mux.HandleFunc("/health/ready", handleReady)

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

	alive.Store(false)
	ready.Store(false)

	logger.Info("Shutting down ServeGO")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Printf("Graceful shutdown Failed: %v", err)
	}
	logger.Warn("Shutdown completed gracefully")

}

func handleRoot(w http.ResponseWriter, r *http.Request) {
	writeJson(w, r, http.StatusOK, "Welcome to ServeGo", "Success")

}

func handleCheck(w http.ResponseWriter, r *http.Request) {
	if alive.Load() {
		writeJson(w, r, http.StatusOK, "Instance is Alive", "Success")
		return
	}
	writeJson(w, r, http.StatusServiceUnavailable, "Instance is NOT Alive", "Fail")
}

func handleReady(w http.ResponseWriter, r *http.Request) {
	if ready.Load() {
		writeJson(w, r, http.StatusOK, "Instance is Ready", "Success")
		return
	}
	writeJson(w, r, http.StatusServiceUnavailable, "Instance is NOT Ready", "Fail")
}

func writeJson(w http.ResponseWriter, r *http.Request, code int, message string, status string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	resp := Response{
		Status:  status,
		Message: message,
		TraceId: uuid.NewString(),
		Path:    r.URL.Path,
		TS:      time.Now().Format(time.RFC3339),
	}
	json.NewEncoder(w).Encode(resp)
}
