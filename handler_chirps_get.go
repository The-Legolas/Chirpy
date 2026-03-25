package main

import (
	"net/http"
)

func (cfg *apiConfig) handlerChirpsRetrieve(w http.ResponseWriter, r *http.Request) {
	DBChirp, err := cfg.db.GetAllChirps(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't fetch chirps", err)
		return
	}

	jsonChirp := make([]Chirp, 0, len(DBChirp))

	for _, chirp := range DBChirp {
		jsonChirp = append(jsonChirp, Chirp{
			ID:        chirp.ID,
			CreatedAt: chirp.CreatedAt,
			UpdatedAt: chirp.UpdatedAt,
			Body:      chirp.Body,
			UserID:    chirp.UserID,
		})
	}

	respondWithJSON(w, http.StatusOK, jsonChirp)
}
