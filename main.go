package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"sync/atomic"

	httputil "github.com/Joshua-Lucas/go-chirpy/internal"
	"github.com/Joshua-Lucas/go-chirpy/internal/database"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type apiConfig struct {
	fileserverHits atomic.Int32
	dbQueries      *database.Queries
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
		CleanedBody string `json:"cleaned_body"`
	}

	decoder := json.NewDecoder(r.Body)
	b := body{}
	err := decoder.Decode(&b)

	if err != nil {
		httputil.RespondWithError(w, http.StatusBadRequest, "Something went wrong")
		return
	}

	max_characters := 140

	// Validate length of string
	if len(b.Body) > max_characters {
		httputil.RespondWithError(w, http.StatusBadRequest, "Chirp is too long")
		return
	}

	// Validate profane words
	var profaneWords = map[string]struct{}{
		"kerfuffle": {},
		"sharbert":  {},
		"fornax":    {},
	}

	words := strings.Fields(b.Body)
	for i := range words {
		if _, ok := profaneWords[strings.ToLower(words[i])]; ok {
			words[i] = "****"
		}
	}

	parsedWords := strings.Join(words, " ")

	respBody := returnVals{
		CleanedBody: parsedWords,
	}

	httputil.RespondWithJSON(w, http.StatusOK, respBody)

}

func main() {
	godotenv.Load()

	// DB Connection
	dbURL := os.Getenv("DB_URL")

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal(err)
	}

	dbQueries := database.New(db)

	apiCfg := apiConfig{
		fileserverHits: atomic.Int32{},
		dbQueries:      dbQueries,
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

	log.Fatal(server.ListenAndServe())

}
