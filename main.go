package main

import (
	"fmt"
	"net/http"
	"sync/atomic"
)

type apiConfig struct {
	fileserverHits atomic.Int32
}

func main() {
	const port = "8080"
	const fileRoot = "."

	cfg := &apiConfig{}
	serveMux := http.NewServeMux()

	fileHandler := http.StripPrefix("/app", http.FileServer(http.Dir(fileRoot)))

	serveMux.Handle("/app/", cfg.middlewareMetricsInc(fileHandler))
	serveMux.HandleFunc("GET /api/healthz", handlerReady)
	serveMux.HandleFunc("GET /admin/metrics", cfg.handlerMetrics)
	serveMux.HandleFunc("POST /admin/reset", cfg.handlerResetMetrics)

	server := &http.Server{
		Handler: serveMux,
		Addr:    ":" + port,
	}

	err := server.ListenAndServe()
	if err != nil {
		fmt.Printf("Error: %v", err)
	}
}
