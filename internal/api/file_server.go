package api

import (
	"fmt"
	"net/http"
)

func (cfg *ApiConfig) HandleFilserverHits(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	result := fmt.Sprintf(`<html>
  <body>
    <h1>Welcome, Chirpy Admin</h1>
    <p>Chirpy has been visited %d times!</p>
  </body>
</html>`,
		cfg.FileserverHits.Load())
	w.Write([]byte(result))

}
