package main

import (
	"encoding/json"
	"net/http"
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

	type valid struct {
		Valid bool `json:"valid"`
	}
	validRes := valid{Valid: true}
	respondWithJSON(w, http.StatusOK, validRes)
}
