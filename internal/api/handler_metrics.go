package api

import (
	"fmt"
	"net/http"
)

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
