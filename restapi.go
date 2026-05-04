package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/swampbear/chirpy/internal/auth"
	"github.com/swampbear/chirpy/internal/database"
)

func handleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func (cfg *apiConfig) handleFilserverHits(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	result := fmt.Sprintf(`<html>
  <body>
    <h1>Welcome, Chirpy Admin</h1>
    <p>Chirpy has been visited %d times!</p>
  </body>
</html>`,
		cfg.fileserverHits.Load())
	w.Write([]byte(result))

}

func (cfg *apiConfig) handleReset(w http.ResponseWriter, r *http.Request) {
	if cfg.platform != "dev" {
		w.WriteHeader(401)
		w.Write([]byte("FORBIDDEN"))
		return
	}

	//restet server hits
	cfg.fileserverHits.Store(0)

	//delete all users
	cfg.db.DeleteAllUsers(r.Context())
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	result := fmt.Sprintf("Reset server count back to %d", cfg.fileserverHits.Load())

	w.Write([]byte(result))
}
func (cfg *apiConfig) handleChirps(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		UserId string `json:"user_id"`
		Body   string `json:"body"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, 500, "something went wrong")
		return
	}

	// check limit
	chirpyWordLimit := 140

	if len(params.Body) > chirpyWordLimit {
		respondWithError(w, 400, "Chirp is to long")
		return
	}

	cleanedText := cleanChirp(params.Body)
	userId, err := uuid.Parse(params.UserId)
	if err != nil {
		fmt.Printf("failed to parse user id %w", err)
	}

	chirpParams := database.CreateChirpParams{Body: cleanedText, UserID: uuid.NullUUID{UUID: userId, Valid: true}}

	chirpdb, err := cfg.db.CreateChirp(r.Context(), chirpParams)

	chirp := dbChirpToModelChirp(chirpdb)

	respondWithJson(w, 201, chirp)

}

func (cfg *apiConfig) handleGetAllChrips(w http.ResponseWriter, r *http.Request) {
	chirpsDB, err := cfg.db.GetAllChirps(r.Context())
	if err != nil {
		log.Println(err)
		respondWithError(w, 500, "Error: failed to connect to database: %w")
		return
	}
	chirps := []Chirp{}

	for _, chirp := range chirpsDB {
		chirps = append(chirps, dbChirpToModelChirp(chirp))
	}
	respondWithJson(w, 200, chirps)
}

func (cfg *apiConfig) handleGetChirpByID(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(r.PathValue("chirpID"))
	if err != nil {
		respondWithError(w, 404, "Failed to parse id")
		return
	}
	chirpDB, err := cfg.db.GetChirpByID(r.Context(), id)
	if err != nil {
		respondWithError(w, 404, "Error: chirp does not exist")
		return
	}
	chirp := dbChirpToModelChirp(chirpDB)
	respondWithJson(w, 200, chirp)

}

func (cfg *apiConfig) handleCreateUser(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	decoder := json.NewDecoder(r.Body)

	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		w.WriteHeader(400)
		log.Println("ERROR: failed to decode parameters")
		return
	}
	hashedPass, err := auth.HashPassword(params.Password)

	if err != nil {
		respondWithError(w, 400, "error: failed to hash password %w")
		return

	}

	userDBParams := database.CreateUsersParams{Email: params.Email, HashedPassword: sql.NullString{String: hashedPass, Valid: true}}

	dbuser, err := cfg.db.CreateUsers(r.Context(), userDBParams)
	if err != nil {
		w.WriteHeader(500)
		log.Printf("ERROR: failed to create user: %w", err)
		return
	}
	user := User{ID: dbuser.ID, CreatedAt: dbuser.CreatedAt, UpdatedAt: dbuser.UpdatedAt, Email: dbuser.Email}
	respondWithJson(w, 201, user)
}

func (cfg *apiConfig) handleLogin(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	if err := decoder.Decode(&params); err != nil {
		respondWithError(w, 401, "error: failed to decode parameters")
		return
	}

	dbuser, err := cfg.db.GetUserByEmail(r.Context(), params.Email)
	if err != nil {
		respondWithError(w, 401, "Unauthorized")
		return
	}
	isHashed, err := auth.CheckPasswordHash(params.Password, dbuser.HashedPassword.String)
	if err != nil || !isHashed {
		respondWithError(w, 401, "Unauthorized")
		return
	}

	user := User{ID: dbuser.ID, CreatedAt: dbuser.CreatedAt, UpdatedAt: dbuser.UpdatedAt, Email: dbuser.Email}
	respondWithJson(w, 200, user)
}
