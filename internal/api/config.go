package api

import (
	"sync/atomic"

	_ "github.com/lib/pq"
	"github.com/swampbear/chirpy/internal/database"
)

type ApiConfig struct {
	FileserverHits atomic.Int32
	Db             *database.Queries
	Platform       string
	TokenString    string
	PolkaKey       string
}
