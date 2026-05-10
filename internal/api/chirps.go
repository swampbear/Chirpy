package api

import (
	"encoding/json"
	"github.com/google/uuid"
	"github.com/swampbear/chirpy/internal/database"
	"log"
	"net/http"
)

func (cfg *ApiConfig) HandleChirps(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, 500, "something went wrong")
		return
	}

	// check limit
	if len(params.Body) > CHIRPY_LIMIT {
		respondWithError(w, 400, "Chirp is to long")
		return
	}

	cleanedText := cleanChirp(params.Body)

	userId := r.Context().Value("userId").(uuid.UUID)
	// retrieve chirp from database
	chirpParams := database.CreateChirpParams{Body: cleanedText, UserID: uuid.NullUUID{UUID: userId, Valid: true}}
	chirpdb, err := cfg.Db.CreateChirp(r.Context(), chirpParams)
	chirp := parseChirp(chirpdb)

	respondWithJson(w, 201, chirp)

}

func (cfg *ApiConfig) HandleGetAllChrips(w http.ResponseWriter, r *http.Request) {
	chirpsDB, err := cfg.Db.GetAllChirps(r.Context())
	if err != nil {
		log.Println(err)
		respondWithError(w, 500, "failed to connect to database: %w")
		return
	}
	chirps := []Chirp{}

	for _, chirp := range chirpsDB {
		chirps = append(chirps, parseChirp(chirp))
	}
	respondWithJson(w, 200, chirps)
}

func (cfg *ApiConfig) HandleGetChirpByID(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(r.PathValue("chirpID"))
	if err != nil {
		respondWithError(w, 404, "parse chirp id")
		return
	}
	chirpDB, err := cfg.Db.GetChirpByID(r.Context(), id)
	if err != nil {
		respondWithError(w, 404, "chirp does not exist")
		return
	}
	chirp := parseChirp(chirpDB)
	respondWithJson(w, 200, chirp)
}

func (cfg *ApiConfig) HandleDeleteChirp(w http.ResponseWriter, r *http.Request) {
	// get values from request
	id, err := uuid.Parse(r.PathValue("chirpID"))

	if err != nil {
		respondWithError(w, 404, "parse chirp id")
		return
	}
	chirpDB, err := cfg.Db.GetChirpByID(r.Context(), id)
	if err != nil {
		respondWithError(w, 404, "chirp does not exist")
		return
	}

	//get user id from context and check if it matches the foreign key in chirpDB
	userId := r.Context().Value("userId").(uuid.UUID)
	if chirpDB.UserID.UUID != userId {
		respondWithError(w, 403, "unauthorized")
		return
	}
	if err = cfg.Db.DeleteByID(r.Context(), id); err != nil {
		respondWithError(w, 404, "failed to delete chirp")
		return
	}

	// return status for deletion was successfull
	w.WriteHeader(204)

}
