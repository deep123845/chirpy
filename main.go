package main

import (
	"fmt"
	"net/http"
)

func main() {
	serveMux := http.NewServeMux()
	handler := http.FileServer(http.Dir("."))
	serveMux.Handle("/", handler)

	server := http.Server{
		Handler: serveMux,
		Addr:    ":8080",
	}

	err := server.ListenAndServe()
	if err != nil {
		fmt.Printf("Error: %v", err)
	}
}
