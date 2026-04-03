package main

import (
	"chirpy/internal/database"
	"net/http"

	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerChirpsGet(w http.ResponseWriter, r *http.Request) {

	chirpID := r.PathValue("chirpID")

	parsedID, err := uuid.Parse(chirpID)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid chirp ID", err)
		return
	}

	DBChirp, err := cfg.db.GetChirp(r.Context(), parsedID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Couldn't fetch chirp", err)
		return
	}

	jsonChirp := Chirp{
		ID:        DBChirp.ID,
		CreatedAt: DBChirp.CreatedAt,
		UpdatedAt: DBChirp.UpdatedAt,
		Body:      DBChirp.Body,
		UserID:    DBChirp.UserID,
	}

	respondWithJSON(w, http.StatusOK, jsonChirp)
}

func (cfg *apiConfig) handlerChirpsRetrieve(w http.ResponseWriter, r *http.Request) {
	var DBChirp []database.Chirp
	var err error

	author_id_string := r.URL.Query().Get("author_id")

	if author_id_string != "" {

		author_uuid, err := uuid.Parse(author_id_string)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "Error converting to uuid", err)
			return
		}

		DBChirp, err = cfg.db.GetChirpViaAuthor(r.Context(), author_uuid)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Couldn't fetch chirps", err)
			return
		}
	} else {
		DBChirp, err = cfg.db.GetAllChirps(r.Context())
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Couldn't fetch chirps", err)
			return
		}
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
