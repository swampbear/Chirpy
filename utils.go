package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/swampbear/chirpy/internal/auth"
	"github.com/swampbear/chirpy/internal/database"
)

func dbChirpToModelChirp(chirpdb database.Chirp) Chirp {
	chirp := Chirp{ID: chirpdb.ID, CreatedAt: chirpdb.CreatedAt, UpdatedAt: chirpdb.UpdatedAt, Body: chirpdb.Body, UserId: chirpdb.UserID.UUID}
	return chirp

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

func getExepirationTime(seconds int) time.Duration {
	// configure expiration if unconfigured or zero set to 1 hour, or if over 1 hour limit it to one hour

	if seconds == 0 || seconds > TOKEN_EXPIRATION_LIMIT_SECONDS {
		seconds = TOKEN_EXPIRATION_LIMIT_SECONDS
	}
	// convert to duration
	expirationTime := time.Second * time.Duration(seconds)

	return expirationTime

}

func respondWithJson(w http.ResponseWriter, code int, payload any) {
	w.WriteHeader(code)
	res, err := json.Marshal(payload)
	if err != nil {
		_ = fmt.Errorf("marshall payload message: %s", err)
	}
	w.Write(res)
}

func respondWithError(w http.ResponseWriter, code int, msg string) {
	w.WriteHeader(code)
	errormsg := map[string]string{"error": msg}
	res, err := json.Marshal(errormsg)
	if err != nil {
		log.Printf("marshall error message: %s", err)
	}
	w.Write(res)
}
