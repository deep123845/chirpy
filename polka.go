package main

import (
	"encoding/json"
	"net/http"

	"github.com/deep123845/chirpy/internal/database"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerPolkaWebhook(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Event string `json:"event"`
		Data  struct {
			UserID uuid.UUID `json:"user_id"`
		} `json:"data"`
	}
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Coundn't decode parameters", err)
		return
	}

	const userUpgradeEvent = "user.upgraded"
	if params.Event != userUpgradeEvent {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	_, err = cfg.db.UpdateUserRedStatus(r.Context(), database.UpdateUserRedStatusParams{
		ID:          params.Data.UserID,
		IsChirpyRed: true,
	})
	if err != nil {
		respondWithError(w, http.StatusNotFound, "User not Upgraded", err)
	}

	w.WriteHeader(http.StatusNoContent)
}
