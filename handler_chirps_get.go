package main

import (
	"chirpy/internal/database"
	"net/http"
	"sort"

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

func authorIDFromRequest(r *http.Request) (uuid.UUID, error) {
	authorIDString := r.URL.Query().Get("author_id")
	if authorIDString == "" {
		return uuid.Nil, nil
	}
	authorID, err := uuid.Parse(authorIDString)
	if err != nil {
		return uuid.Nil, err
	}
	return authorID, nil
}

func (cfg *apiConfig) handlerChirpsRetrieve(w http.ResponseWriter, r *http.Request) {
	var dbChirp []database.Chirp
	var err error

	authorID, err := authorIDFromRequest(r)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid author ID", err)
		return
	}

	if authorID != uuid.Nil {
		dbChirp, err = cfg.db.GetChirpViaAuthor(r.Context(), authorID)
	} else {
		dbChirp, err = cfg.db.GetAllChirps(r.Context())
	}
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't fetch chirps", err)
		return
	}

	jsonChirp := make([]Chirp, 0, len(dbChirp))
	for _, chirp := range dbChirp {
		jsonChirp = append(jsonChirp, Chirp{
			ID:        chirp.ID,
			CreatedAt: chirp.CreatedAt,
			UpdatedAt: chirp.UpdatedAt,
			Body:      chirp.Body,
			UserID:    chirp.UserID,
		})
	}

	sortDirection := r.URL.Query().Get("sort")

	sort.Slice(jsonChirp, func(i, j int) bool {
		if sortDirection == "desc" {
			return jsonChirp[i].CreatedAt.After(jsonChirp[j].CreatedAt)
		}
		return jsonChirp[i].CreatedAt.Before(jsonChirp[j].CreatedAt)
	})

	respondWithJSON(w, http.StatusOK, jsonChirp)
}
