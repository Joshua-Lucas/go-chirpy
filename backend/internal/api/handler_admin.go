package api

import (
	"fmt"
	"net/http"

	httputil "github.com/Joshua-Lucas/go-chirpy/internal"
)

func (cfg *APIConfig) MetricHandlerServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	fmtMessage := fmt.Sprintf(`<html>
  <body>
    <h1>Welcome, Chirpy Admin</h1>
    <p>Chirpy has been visited %d times!</p>
  </body>
</html>`, cfg.FileserverHits.Load())
	w.Write([]byte(fmtMessage))
}

func (cfg *APIConfig) ResetHandlerServeHTTP(w http.ResponseWriter, r *http.Request) {
	if cfg.Platform != "dev" {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	err := cfg.DBQueries.DeleteAllUsers(r.Context())
	if err != nil {

		httputil.RespondWithError(w, http.StatusInternalServerError, "Something went wrong deleting ")
	}

	w.WriteHeader(http.StatusOK)

}
