package main

import (
	"chirpy/internal/auth"
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
)

type PolkaWebhooks struct {
	Event string `json:"event"`
	Data  struct {
		UserID string `json:"user_id"`
	} `json:"data"`
}

func (cfg *apiConfig) handlerWebhooks(w http.ResponseWriter, r *http.Request) {
	_, err := auth.GetAPIKey(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't find polka key", err)
		return
	}

	decoder := json.NewDecoder(r.Body)
	params := PolkaWebhooks{}
	err = decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Couldn't decode parameters", err)
		return
	}

	if params.Event != "user.upgraded" {
		respondWithJSON(w, http.StatusNoContent, nil)
		return
	}

	userID, err := uuid.Parse(params.Data.UserID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't convert string to uuid", err)
		return
	}

	_, err = cfg.db.UpgradeUser(r.Context(), userID)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Error upgrading the user", err)
		return
	}

	respondWithJSON(w, http.StatusNoContent, nil)
}
