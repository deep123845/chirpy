package main

import (
	"encoding/json"
	"net/http"
	"slices"
	"strings"
	"time"

	"github.com/deep123845/chirpy/internal/auth"
	"github.com/deep123845/chirpy/internal/database"
	"github.com/google/uuid"
)

type chirp struct {
	Id        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	UserId    uuid.UUID `json:"user_id"`
}

func (cfg *apiConfig) handlerCreateChirp(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	parameters := chirp{}
	err := decoder.Decode(&parameters)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	bearerToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't get authorization", err)
		return
	}

	userId, err := auth.ValidateJWT(bearerToken, cfg.jwtSecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Unauthorized", err)
		return
	}

	const maxChirpLength = 140
	if len(parameters.Body) > maxChirpLength {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long", nil)
		return
	}

	censoredBody := censorMessage(parameters.Body)

	data, err := cfg.db.CreateChirp(r.Context(), database.CreateChirpParams{
		Body:   censoredBody,
		UserID: userId,
	})

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to create chirp", err)
		return
	}

	respondWithJSON(w, http.StatusCreated, chirp{
		Id:        data.ID,
		CreatedAt: data.CreatedAt,
		UpdatedAt: data.UpdatedAt,
		Body:      data.Body,
		UserId:    data.UserID,
	})
}

func (cfg *apiConfig) handlerGetChirps(w http.ResponseWriter, r *http.Request) {
	data, err := cfg.db.GetAllChirps(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to get all chirps", err)
		return
	}

	var chirps []chirp
	for _, entry := range data {
		chirps = append(chirps, chirp{
			Id:        entry.ID,
			CreatedAt: entry.CreatedAt,
			UpdatedAt: entry.UpdatedAt,
			Body:      entry.Body,
			UserId:    entry.UserID,
		})
	}

	respondWithJSON(w, http.StatusOK, chirps)
}

func (cfg *apiConfig) handlerGetChirp(w http.ResponseWriter, r *http.Request) {
	chirpId, err := uuid.Parse(r.PathValue("chirpID"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Must Provide valid chirp id", err)
		return
	}

	data, err := cfg.db.GetChirpById(r.Context(), chirpId)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Chirp not found", err)
		return
	}

	respondWithJSON(w, http.StatusOK, chirp{
		Id:        data.ID,
		CreatedAt: data.CreatedAt,
		UpdatedAt: data.UpdatedAt,
		Body:      data.Body,
		UserId:    data.UserID,
	})
}

func censorMessage(msg string) string {
	var bannedWords = []string{"kerfuffle", "sharbert", "fornax"}

	words := strings.Split(msg, " ")

	for i, word := range words {
		word = strings.ToLower(word)

		if slices.Contains(bannedWords, word) {
			words[i] = "****"
		}
	}

	return strings.Join(words, " ")
}
