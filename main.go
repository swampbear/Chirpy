package main

import (
	"log"
	"net/http"
	"sync/atomic"
)

func main() {
	mux := http.NewServeMux()

	port := "8080"

	server := http.Server{Addr: ":" + port, Handler: mux}
	apiCfg := apiConfig{
		fileserverHits: atomic.Int32{},
	}

	// registering handlers to mux
	mux.Handle("/app/", http.StripPrefix("/app", apiCfg.middlwareMetricsInc(middlewareLog(http.FileServer(http.Dir("."))))))
	mux.HandleFunc("GET /api/healthz", handleHealth)
	mux.HandleFunc("GET /admin/metrics", apiCfg.handleFilserverHits)
	mux.HandleFunc("POST /admin/reset", apiCfg.handleReset)

	mux.HandleFunc("POST /api/validate_chirp", handleValidateChirp)

	log.Printf("Listening on port: %s", server.Addr)
	http.ListenAndServe(server.Addr, server.Handler)

}
