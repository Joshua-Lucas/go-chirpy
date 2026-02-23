package api

import (
	"sync/atomic"

	"github.com/Joshua-Lucas/go-chirpy/internal/database"
)

type APIConfig struct {
	FileserverHits atomic.Int32
	DBQueries      *database.Queries
}
