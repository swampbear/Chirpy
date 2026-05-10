package api

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/swampbear/chirpy/internal/auth"
	"github.com/swampbear/chirpy/internal/database"
)

func (cfg *ApiConfig) HandleCreateUser(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	decoder := json.NewDecoder(r.Body)

	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		w.WriteHeader(400)
		log.Println("failed to decode parameters")
		return
	}
	hashedPass, err := auth.HashPassword(params.Password)

	if err != nil {
		respondWithError(w, 400, "failed to hash password %w")
		return

	}

	userDBParams := database.CreateUsersParams{Email: params.Email, HashedPassword: sql.NullString{String: hashedPass, Valid: true}}

	dbuser, err := cfg.Db.CreateUsers(r.Context(), userDBParams)
	if err != nil {
		w.WriteHeader(500)
		log.Printf("failed to create user: %s", err)
		return
	}
	user := parseUser(dbuser)
	respondWithJson(w, 201, user)
}

func (cfg *ApiConfig) HandleUpdateUser(w http.ResponseWriter, r *http.Request) {
	// set up parameters struct and decode body
	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	if err := decoder.Decode(&params); err != nil {
		respondWithError(w, 401, "decode parameters")
		return
	}
	hashedPassword, err := auth.HashPassword(params.Password)
	if err != nil {
		respondWithError(w, 401, "hash password")
		return
	}

	userId := r.Context().Value("userId").(uuid.UUID)
	updateParams := database.UpdateUserByIdParams{Email: params.Email, HashedPassword: sql.NullString{String: hashedPassword, Valid: true}, ID: userId}

	dbuser, err := cfg.Db.UpdateUserById(r.Context(), updateParams)
	if err != nil {
		respondWithError(w, 500, "update user")
		return
	}
	user := parseUser(dbuser)
	// respond with 200
	respondWithJson(w, 200, user)

}
