package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/swampbear/chirpy/internal/api"
	"github.com/swampbear/chirpy/internal/database"
)

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
	tokenString := os.Getenv("TOKEN_STRING")

	apiCfg := api.ApiConfig{
		FileserverHits: atomic.Int32{},
		Db:             dbQueries,
		Platform:       platform,
		TokenString:    tokenString,
	}

	// configure server
	mux := http.NewServeMux()
	port := "8080"
	server := http.Server{Addr: ":" + port, Handler: mux}

	// registering handlers to mux
	mux.Handle("/app/", http.StripPrefix("/app", http.FileServer(http.Dir("."))))
	mux.HandleFunc("GET /api/healthz", api.HandleHealth)
	mux.HandleFunc("GET /admin/metrics", apiCfg.HandleFilserverHits)
	mux.HandleFunc("POST /admin/reset", apiCfg.HandleReset)
	mux.HandleFunc("POST /api/users", apiCfg.HandleCreateUser)
	mux.HandleFunc("GET /api/chirps", apiCfg.HandleGetAllChrips)
	mux.HandleFunc("GET /api/chirps/{chirpID}", apiCfg.HandleGetChirpByID)
	mux.HandleFunc("POST /api/login", apiCfg.HandleLogin)
	mux.HandleFunc("POST /api/refresh", apiCfg.HandleRefresh)
	mux.HandleFunc("POST /api/revoke", apiCfg.HandleRevoke)

	// protected endpoints
	mux.Handle("PUT /api/users", apiCfg.MiddlewareAuth(apiCfg.HandleUpdateUser))
	mux.Handle("POST /api/chirps", apiCfg.MiddlewareAuth(apiCfg.HandleChirps))
	mux.Handle("DELETE /api/chirps/{chirpID}", apiCfg.MiddlewareAuth(apiCfg.HandleDeleteChirp))

	// add middleware
	handler := apiCfg.MiddlwareMetricsInc(mux)
	handler = api.MiddlewareLog(handler)

	// start listening
	log.Printf("Listening on port: %s", server.Addr)
	http.ListenAndServe(server.Addr, handler)

}
