package api

import (
	"fmt"
	"net/http"
)

func (cfg *ApiConfig) HandleReset(w http.ResponseWriter, r *http.Request) {
	if cfg.Platform != "dev" {
		w.WriteHeader(401)
		w.Write([]byte("FORBIDDEN"))
		return
	}

	//restet server hits
	cfg.FileserverHits.Store(0)

	//delete all users
	cfg.Db.DeleteAllUsers(r.Context())
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	result := fmt.Sprintf("Reset server count back to %d", cfg.FileserverHits.Load())

	w.Write([]byte(result))
}
