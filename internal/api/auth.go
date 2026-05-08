package api

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/swampbear/chirpy/internal/auth"
	"github.com/swampbear/chirpy/internal/database"
)

func (cfg *ApiConfig) HandleLogin(w http.ResponseWriter, r *http.Request) {
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

	// get user
	dbuser, err := cfg.Db.GetUserByEmail(r.Context(), params.Email)
	if err != nil {
		respondWithError(w, 401, "Unauthorized")
		return
	}
	isHashed, err := auth.CheckPasswordHash(params.Password, dbuser.HashedPassword.String)
	if err != nil || !isHashed {
		respondWithError(w, 401, "Unauthorized")
		return
	}

	// create token
	token, err := auth.MakeJWT(dbuser.ID, cfg.TokenString, TOKEN_EXPIRATION_LIMIT_SECONDS)
	if err != nil {
		respondWithError(w, 400, err.Error())
		return
	}

	// create and save refresh token
	refToken := auth.MakeRefreshToken()
	refTokenParams := database.CreateRefreshTokenParams{Token: refToken, UserID: uuid.NullUUID{UUID: dbuser.ID, Valid: true}, ExpiresAt: time.Now().Add(REFRESH_TOKEN_EXPIRATION_LIMIT_SECONDS)}
	dbRefToken, err := cfg.Db.CreateRefreshToken(r.Context(), refTokenParams)
	if err != nil {
		respondWithError(w, 400, err.Error())
		return
	}

	// map to return user struct
	user := User{ID: dbuser.ID, CreatedAt: dbuser.CreatedAt, UpdatedAt: dbuser.UpdatedAt, Email: dbuser.Email, Token: token, RefreshToken: dbRefToken.Token}
	respondWithJson(w, 200, user)
}

func (cfg *ApiConfig) HandleRefresh(w http.ResponseWriter, r *http.Request) {
	// parse refresh token from authorization header
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, 400, err.Error())
	}
	dbToken, err := cfg.Db.GetRefreshTokenFromToken(r.Context(), token)
	if err != nil {
		respondWithError(w, 401, err.Error())
		return
	}

	// check if refresh token is valid
	currentTime := time.Now()
	if currentTime.After(dbToken.ExpiresAt) {
		respondWithError(w, 401, "refresh token is expired")
		return
	}
	if currentTime.After(dbToken.RevokedAt.Time) && !dbToken.RevokedAt.Time.IsZero() {
		respondWithError(w, 401, "refresh token has been revoked")
		return
	}
	dbuser, err := cfg.Db.GetUserFromRefreshToken(r.Context(), token)
	if err != nil {
		respondWithError(w, 401, err.Error())
		return
	}

	// create token
	newToken, err := auth.MakeJWT(dbuser.ID, cfg.TokenString, TOKEN_EXPIRATION_LIMIT_SECONDS)
	if err != nil {
		respondWithError(w, 400, err.Error())
		return
	}
	// return the token
	type ReturnParams struct {
		Token string `json:"token"`
	}
	respondWithJson(w, 200, ReturnParams{Token: newToken})
}

func (cfg *ApiConfig) HandleRevoke(w http.ResponseWriter, r *http.Request) {
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, 400, err.Error())
		return
	}

	if err = cfg.Db.RevokeRefreshToken(r.Context(), token); err != nil {
		respondWithError(w, 400, err.Error())
		return
	}
	w.WriteHeader(204)
}
