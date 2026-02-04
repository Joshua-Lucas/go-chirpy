package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync/atomic"
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

	decoder := json.NewDecoder(r.Body)
	b := body{}
	err := decoder.Decode(&b)

	if err != nil {
		log.Printf("Error decoding parameters: %s", err)
		w.WriteHeader(500)

		type error struct {
			Error string `json:"error"`
		}

		errMessage := error{
			Error: "Something went wrong",
		}

		dat, err := json.Marshal(errMessage)
		if err != nil {
			log.Printf("Error marshalling JSON: %s", err)
			w.WriteHeader(500)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(dat)
		return
	}

	CHARACTER_COUNT := 140

	// Validate length of string
	if len(b.Body) > CHARACTER_COUNT {
		type validationMsg struct {
			Error string `json:"error"`
		}

		msg := validationMsg{
			Error: "Chirp is too long",
		}

		dat, err := json.Marshal(msg)
		if err != nil {
			log.Printf("Error marshalling JSON: %s", err)
			w.WriteHeader(500)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(400)

		w.Write(dat)
		return
	}

	type returnVals struct {
		Valid bool `json:"valid"`
	}

	respBody := returnVals{
		Valid: true,
	}

	dat, err := json.Marshal(respBody)
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
		w.WriteHeader(500)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	w.Write(dat)

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
