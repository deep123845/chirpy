package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/deep123845/chirpy/internal/database"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type apiConfig struct {
	fileserverHits atomic.Int32
	db             *database.Queries
	dev            bool
	jwtSecret      string
}

func main() {
	godotenv.Load()

	const port = "8080"
	const fileRoot = "."

	dbURL := os.Getenv("DB_URL")
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("Could not open DB, %v", err)
	}

	dev := os.Getenv("PLATFORM") == "dev"
	jwtSecret := os.Getenv("JWT_SECRET")

	dbQueries := database.New(db)

	cfg := &apiConfig{
		db:        dbQueries,
		dev:       dev,
		jwtSecret: jwtSecret,
	}
	serveMux := http.NewServeMux()

	fileHandler := http.StripPrefix("/app", http.FileServer(http.Dir(fileRoot)))

	serveMux.Handle("/app/", cfg.middlewareMetricsInc(fileHandler))

	serveMux.HandleFunc("GET /api/healthz", handlerReady)
	serveMux.HandleFunc("POST /api/users", cfg.handlerCreateUser)
	serveMux.HandleFunc("POST /api/login", cfg.handlerLogin)
	serveMux.HandleFunc("POST /api/chirps", cfg.handlerCreateChirp)
	serveMux.HandleFunc("GET /api/chirps", cfg.handlerGetChirps)
	serveMux.HandleFunc("GET /api/chirps/{chirpID}", cfg.handlerGetChirp)

	serveMux.HandleFunc("GET /admin/metrics", cfg.handlerMetrics)
	serveMux.HandleFunc("POST /admin/reset", cfg.handlerReset)

	server := &http.Server{
		Handler: serveMux,
		Addr:    ":" + port,
	}

	err = server.ListenAndServe()
	if err != nil {
		fmt.Printf("Error: %v", err)
	}
}
