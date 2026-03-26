package main

import (
	"encoding/json"
	"net/http"
	"slices"
	"strings"
	"time"

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

	const maxChirpLength = 140
	if len(parameters.Body) > maxChirpLength {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long", nil)
		return
	}

	censoredBody := censorMessage(parameters.Body)

	data, err := cfg.db.CreateChirp(r.Context(), database.CreateChirpParams{
		Body:   censoredBody,
		UserID: parameters.UserId,
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
