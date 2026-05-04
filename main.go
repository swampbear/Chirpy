package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/swampbear/chirpy/internal/database"
)

type apiConfig struct {
	fileserverHits atomic.Int32
	db             *database.Queries
	platform       string
}

func main() {
	//load .env file
	godotenv.Load()

	// setup database
	dbUrl := os.Getenv("DB_URL")
	log.Println(dbUrl)
	db, err := sql.Open("postgres", dbUrl)
	if err != nil {
		log.Println("ERROR: failed to connect to database")
	}
	dbQueries := database.New(db)

	// setup api config

	platform := os.Getenv("PLATFORM")
	apiCfg := apiConfig{
		fileserverHits: atomic.Int32{},
		db:             dbQueries,
		platform:       platform,
	}

	mux := http.NewServeMux()

	port := "8080"

	server := http.Server{Addr: ":" + port, Handler: mux}

	// registering handlers to mux
	mux.Handle("/app/", http.StripPrefix("/app", http.FileServer(http.Dir("."))))
	mux.HandleFunc("GET /api/healthz", handleHealth)
	mux.HandleFunc("GET /admin/metrics", apiCfg.handleFilserverHits)
	mux.HandleFunc("POST /admin/reset", apiCfg.handleReset)
	mux.HandleFunc("POST /api/users", apiCfg.handleCreateUser)
	mux.HandleFunc("POST /api/chirps", apiCfg.handleChirps)
	mux.HandleFunc("GET /api/chirps", apiCfg.handleGetAllChrips)
	mux.HandleFunc("GET /api/chirps/{chirpID}", apiCfg.handleGetChirpByID)
	mux.HandleFunc("POST /api/login", apiCfg.handleLogin)

	// add middleware
	handler := apiCfg.middlwareMetricsInc(mux)
	handler = middlewareLog(handler)

	log.Printf("Listening on port: %s", server.Addr)
	http.ListenAndServe(server.Addr, handler)

}
