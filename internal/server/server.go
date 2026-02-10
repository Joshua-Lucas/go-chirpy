package server

import (
	"net/http"
)

type Server struct {
	mux *http.ServeMux
}

func NewServer(cfg *APIConfig) *Server {
	s := &Server{mux: http.NewServeMux()}
	s.routes(cfg)
	return s
}

func (s Server) routes(cfg *APIConfig) {
	 // Static files route
    fs := http.FileServer(http.Dir("./static")) 
    s.mux.Handle("/app/", cfg.MiddlewareMetricsInc(http.StripPrefix("/app/", fs)))

	RegisterAPIRoutes(s.mux, cfg)
	
  //RegisterAdminRoutes(s.mux, cfg)	
}


func (s *Server) ListenAndServe(addr string) error {
    srv := &http.Server{
        Addr:    addr,
        Handler: s.mux,
    }
    return srv.ListenAndServe()
}
