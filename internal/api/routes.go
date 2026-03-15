package api

import (
	"net/http"
)

func (cfg *APIConfig) RegisterRoutes(mux *http.ServeMux) {
	// Static files route
	fs := http.FileServer(http.Dir("./static"))
	mux.Handle("/app/", cfg.MiddlewareMetricsInc(http.StripPrefix("/app/", fs)))

	cfg.RegisterAPIRoutes(mux)
	cfg.RegisterAdminRoutes(mux)
}

func (cfg *APIConfig) RegisterAPIRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/healthz", cfg.HealthHandler)
	mux.HandleFunc("POST /api/users", cfg.CreateUserHandler)
	mux.HandleFunc("GET /api/chirps", cfg.GetAllChripsHandler)
	mux.HandleFunc("POST /api/chirps", cfg.CreateChirpHandler)
}

func (cfg *APIConfig) RegisterAdminRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /admin/metrics", cfg.MetricHandlerServeHTTP)
	mux.HandleFunc("POST /admin/reset", cfg.ResetHandlerServeHTTP)
}
