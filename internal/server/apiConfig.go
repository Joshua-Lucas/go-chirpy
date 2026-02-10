package server

import (
	"fmt"
	"net/http"
	"sync/atomic"

	"github.com/Joshua-Lucas/go-chirpy/internal/database"
)

type APIConfig struct {
	FileserverHits atomic.Int32
	DBQueries      *database.Queries
}

func (c *APIConfig) MiddlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c.FileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

func (c *APIConfig) MetricHandlerServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	fmtMessage := fmt.Sprintf(`<html>
  <body>
    <h1>Welcome, Chirpy Admin</h1>
    <p>Chirpy has been visited %d times!</p>
  </body>
</html>`, c.FileserverHits.Load())
	w.Write([]byte(fmtMessage))
}

func (c *APIConfig) ResetMetricsHandlerServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	c.FileserverHits.Store(0)
	w.Write([]byte(("")))
}
