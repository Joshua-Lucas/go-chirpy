package server

import "net/http"

func RegisterAdminRoutes(mux *http.ServeMux, cfg *APIConfig) {
	mux.HandleFunc("GET /admin/metrics", cfg.MetricHandlerServeHTTP)
	mux.HandleFunc("POST /admin/reset", cfg.ResetMetricsHandlerServeHTTP)
}
