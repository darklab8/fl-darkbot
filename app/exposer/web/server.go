package web

import (
	"net/http"

	"github.com/darklab8/fl-darkbot/app/settings/logus"
)

type Server struct {
	router map[route]*endpoint
}

func (s *Server) GetRouter() map[route]*endpoint {
	return s.router
}

func NewServer() *Server {
	s := &Server{}
	s.router = make(map[route]*endpoint)

	return s
}

func (s *Server) RegisterEndpoint(route route, Description string, handler func(w http.ResponseWriter, r *http.Request)) *Server {

	NewEndpoint(s, route, Description, handler)
	return s
}

func (s *Server) GetMux() *http.ServeMux {
	mux := http.NewServeMux()
	for route, endpoint := range s.router {
		mux.Handle(string(route), endpoint)
	}
	logus.Log.Info("started server")
	return mux
}
