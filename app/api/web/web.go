package web

import (
	_ "embed"
	"net/http"
)

type route string

type endpoint struct {
	route       route
	serve       func(w http.ResponseWriter, r *http.Request)
	server      *Server
	Description string
}

func NewEndpoint(s *Server, route route, Description string, serve func(w http.ResponseWriter, r *http.Request)) {
	e := &endpoint{}
	e.route = route
	e.serve = serve
	e.server = s
	e.Description = Description
	s.router[route] = e
}

func (e *endpoint) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if err := e.Check404(w, r); err != nil {
		return
	}
	e.serve(w, r)
}

func InternalServerErrorHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte("500 Internal Server Error"))
}
