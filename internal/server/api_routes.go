package server

import (
	"net/http"

	"github.com/Joshua-Lucas/go-chirpy/internal/handlers/api"
)

func RegisterAPIRoutes(mux *http.ServeMux, cfg *APIConfig) {
	mux.HandleFunc("GET /api/healthz", api.HealthHandler)
	mux.HandleFunc("POST /api/validate_chirp", api.ValidateChirpHandler)
}
