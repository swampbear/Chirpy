package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/google/uuid"
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

	// save chirp to database
	chirpdb, err := cfg.db.CreateChirp(r.Context(), chirpParams)

	chirp := Chirp{ID: chirpdb.ID, CreatedAt: chirpdb.CreatedAt, UpdatedAt: chirpdb.UpdatedAt, Body: chirpdb.Body, UserId: chirpdb.UserID.UUID}

	respondWithJson(w, 201, chirp)

}

func (cfg *apiConfig) handleCreateUser(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email string `json:"email"`
	}

	decoder := json.NewDecoder(r.Body)

	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		w.WriteHeader(400)
		log.Println("ERROR: failed to decode parameters")
		return
	}

	dbuser, err := cfg.db.CreateUsers(r.Context(), params.Email)
	if err != nil {
		w.WriteHeader(500)
		log.Printf("ERROR: failed to create user: %w", err)
		return
	}
	user := User{ID: dbuser.ID, CreatedAt: dbuser.CreatedAt, UpdatedAt: dbuser.UpdatedAt, Email: dbuser.Email}
	respondWithJson(w, 201, user)

}

func cleanChirp(text string) string {
	badWords := map[string]struct{}{
		"kerfuffle": {},
		"sharbert":  {},
		"fornax":    {},
	}
	textSlice := strings.Split(text, " ")
	for i, word := range textSlice {
		if _, ok := badWords[strings.ToLower(word)]; ok {
			textSlice[i] = "****"
			continue
		}
	}
	filteredText := strings.Join(textSlice, " ")
	return filteredText

}

func respondWithJson(w http.ResponseWriter, code int, payload any) {
	w.WriteHeader(code)
	res, err := json.Marshal(payload)
	if err != nil {
		_ = fmt.Errorf("Failed to marshall payload message: %w", err)
	}
	w.Write(res)
}

func respondWithError(w http.ResponseWriter, code int, msg string) {
	w.WriteHeader(code)
	errormsg := map[string]string{"error": msg}
	res, err := json.Marshal(errormsg)
	if err != nil {
		log.Printf("Failed to marshall error message: %w", err)
	}
	w.Write(res)
}
