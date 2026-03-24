package main

import (
	"encoding/json"
	"net/http"
	"slices"
	"strings"
)

type chirp struct {
	Body string `json:"body"`
}

func handlerValidateChirp(w http.ResponseWriter, r *http.Request) {
	const maxChirpLength = 140

	decoder := json.NewDecoder(r.Body)
	currChirp := chirp{}
	err := decoder.Decode(&currChirp)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	if len(currChirp.Body) > maxChirpLength {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long", nil)
		return
	}

	type response struct {
		CleanedBody string `json:"cleaned_body"`
	}
	validRes := response{CleanedBody: censorMessage(currChirp.Body)}
	respondWithJSON(w, http.StatusOK, validRes)
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
