package api

import (
	"context"
	"log"
	"net/http"

	"github.com/swampbear/chirpy/internal/auth"
)

// middleware for loggging server requests
func (cfg *ApiConfig) MiddlwareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.FileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

// logging middleware
func MiddlewareLog(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	})

}

// auth middleware that check bearer token
func (cfg *ApiConfig) MiddlewareAuth(next http.HandlerFunc) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, err := auth.GetBearerToken(r.Header)
		if err != nil {
			respondWithError(w, 401, err.Error())
			return
		}
		// check if jwt uuid matches returned user
		userId, err := auth.ValidateJWT(token, cfg.TokenString)
		if err != nil {
			respondWithError(w, 401, err.Error())
			return
		}
		// add userId to context so it is available for the next handler function
		ctx := context.WithValue(r.Context(), "userId", userId)
		next(w, r.WithContext(ctx))
	})
}
