package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

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
		return
	}
	w.Write(res)
}
