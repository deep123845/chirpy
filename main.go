package main

import (
	"fmt"
	"net/http"
)

func handlerReady(w http.ResponseWriter, _ *http.Request) {
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func main() {
	const port = "8080"
	const fileRoot = "."

	serveMux := http.NewServeMux()

	fileHandler := http.FileServer(http.Dir(fileRoot))

	serveMux.Handle("/app/", http.StripPrefix("/app", fileHandler))
	serveMux.HandleFunc("/healthz", handlerReady)

	server := &http.Server{
		Handler: serveMux,
		Addr:    ":" + port,
	}

	err := server.ListenAndServe()
	if err != nil {
		fmt.Printf("Error: %v", err)
	}
}
