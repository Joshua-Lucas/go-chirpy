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
	RegisterAPIRoutes(s.mux, cfg)
	//RegisterAppRoutes(s.mux)
  //RegisterAdminRoutes(s.mux, cfg)	
}


func (s *Server) ListenAndServe(addr string) error {
    srv := &http.Server{
        Addr:    addr,
        Handler: s.mux,
    }
    return srv.ListenAndServe()
}
