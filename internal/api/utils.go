package api

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/swampbear/chirpy/internal/database"
)

// takes in database chirp and parses to Chirp struct to control fields for json responses
func parseChirp(dbchirp database.Chirp) Chirp {
	chirp := Chirp{ID: dbchirp.ID, CreatedAt: dbchirp.CreatedAt, UpdatedAt: dbchirp.UpdatedAt, Body: dbchirp.Body, UserId: dbchirp.UserID.UUID}
	return chirp
}

// takes in database user and parses to User struct to control fields for json responses
func parseUser(dbuser database.User) User {
	user := User{ID: dbuser.ID, CreatedAt: dbuser.CreatedAt, UpdatedAt: dbuser.UpdatedAt, IsChirpyRed: dbuser.IsChirpyRed.Bool, Email: dbuser.Email}
	return user
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
		log.Printf("marshall: %s", err)
		return
	}
	w.Write(res)
}

func respondWithError(w http.ResponseWriter, code int, msg string) {
	w.WriteHeader(code)
	errormsg := map[string]string{"error": msg}
	res, err := json.Marshal(errormsg)
	if err != nil {
		log.Printf("marshall: %s", err)
		return
	}
	w.Write(res)
}
