package main

import (
	"fmt"
	"net/http"
	"sync/atomic"
)

type apiConfig struct {
	fileserverHits atomic.Int32
}

func handlerReady(w http.ResponseWriter, _ *http.Request) {
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)

	w.Write([]byte("OK"))
}

func (cfg *apiConfig) handlerMetrics(w http.ResponseWriter, _ *http.Request) {
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)

	hits := cfg.fileserverHits.Load()
	hitsString := fmt.Sprintf("Hits: %v", hits)
	w.Write([]byte(hitsString))
}

func (cfg *apiConfig) handlerResetMetrics(w http.ResponseWriter, _ *http.Request) {
	cfg.fileserverHits.Store(0)

	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

func main() {
	const port = "8080"
	const fileRoot = "."

	cfg := &apiConfig{}
	serveMux := http.NewServeMux()

	fileHandler := http.StripPrefix("/app", http.FileServer(http.Dir(fileRoot)))

	serveMux.Handle("/app/", cfg.middlewareMetricsInc(fileHandler))
	serveMux.HandleFunc("/healthz", handlerReady)
	serveMux.HandleFunc("/metrics", cfg.handlerMetrics)
	serveMux.HandleFunc("/reset", cfg.handlerResetMetrics)

	server := &http.Server{
		Handler: serveMux,
		Addr:    ":" + port,
	}

	err := server.ListenAndServe()
	if err != nil {
		fmt.Printf("Error: %v", err)
	}
}
