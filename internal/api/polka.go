package api

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
)

func (cfg *ApiConfig) HandleWebhooks(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Event string `json:"event"`
		Data  struct {
			UserId string `json:"user_id"`
		} `json:"data"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	if err := decoder.Decode(&params); err != nil {
		respondWithError(w, 400, err.Error())
		return
	}

	// event handlers
	switch params.Event {
	case "user.upgraded":
		userId, err := uuid.Parse(params.Data.UserId)
		if err != nil {
			respondWithError(w, 400, "parse uuid")
			return
		}
		if err := cfg.Db.UpgradeToChirpyRedByID(r.Context(), userId); err != nil {
			respondWithError(w, 404, "upgrade chirpy")
			return
		}
		w.WriteHeader(http.StatusNoContent)
		return
	default:
		w.WriteHeader(http.StatusNoContent)
	}

}
