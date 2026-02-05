package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync/atomic"

	httputil "github.com/Joshua-Lucas/go-chirpy/internal"
)

type apiConfig struct {
	fileserverHits atomic.Int32
}

func (c *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

func (c *apiConfig) metricHandlerServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	fmtMessage := fmt.Sprintf(`<html>
  <body>
    <h1>Welcome, Chirpy Admin</h1>
    <p>Chirpy has been visited %d times!</p>
  </body>
</html>`, c.fileserverHits.Load())
	w.Write([]byte(fmtMessage))
}

func (c *apiConfig) resetMetricsHandlerServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	c.fileserverHits.Store(0)
	w.Write([]byte(("")))
}

func validateChirpHandler(w http.ResponseWriter, r *http.Request) {
	type body struct {
		Body string `json:"body"`
	}

	type returnVals struct {
		Valid bool `json:"valid"`
	}

	decoder := json.NewDecoder(r.Body)
	b := body{}
	err := decoder.Decode(&b)

	if err != nil {
		httputil.RespondWithError(w, http.StatusBadRequest, "Something went wrong")
		return
	}

	CHARACTER_COUNT := 140

	// Validate length of string
	if len(b.Body) > CHARACTER_COUNT {
		httputil.RespondWithError(w, http.StatusBadRequest, "Chirp is too long")
		return
	}

	respBody := returnVals{
		Valid: true,
	}

	httputil.RespondWithJSON(w, http.StatusOK, respBody)

}

func main() {

	apiCfg := apiConfig{
		fileserverHits: atomic.Int32{},
	}
	mux := http.NewServeMux()

	server := http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	// Api Routes
	mux.HandleFunc("GET /api/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))

	})

	mux.HandleFunc("POST /api/validate_chirp", validateChirpHandler)

	// App routes
	fs := http.FileServer(http.Dir("."))
	mux.Handle("/app/", apiCfg.middlewareMetricsInc(http.StripPrefix("/app", fs)))

	// Admin Routes
	mux.HandleFunc("GET /admin/metrics", apiCfg.metricHandlerServeHTTP)
	mux.HandleFunc("POST /admin/reset", apiCfg.resetMetricsHandlerServeHTTP)

	server.ListenAndServe()

}
