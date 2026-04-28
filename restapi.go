package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
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

	cfg.fileserverHits.Store(0)
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	result := fmt.Sprintf("Reset server count back to %d", cfg.fileserverHits.Load())

	w.Write([]byte(result))
}

func handleValidateChirp(w http.ResponseWriter, r *http.Request) {
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

	chirpyWordLimit := 140

	if len(params.Body) > chirpyWordLimit {
		respondWithError(w, 400, "Chirp is to long")
		return
	}

	type returnValues struct {
		Valid bool `json:"valid"`
	}

	payload := returnValues{Valid: true}
	respondWithJson(w, 200, payload)

}

func respondWithJson(w http.ResponseWriter, code int, payload interface{}) {
	w.WriteHeader(code)
	res, err := json.Marshal(payload)
	if err != nil {
		log.Fatalf("Failed to marshall payload message: %w", err)
	}
	w.Write(res)
}

func respondWithError(w http.ResponseWriter, code int, msg string) {
	w.WriteHeader(code)
	errormsg := map[string]string{"error": msg}
	res, err := json.Marshal(errormsg)
	if err != nil {
		log.Fatalf("Failed to marshall error message: %w", err)
	}
	w.Write(res)
}
