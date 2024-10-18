package server

import (
	"context"
	"net/http"

	"github.com/gorilla/mux"
)

// Server
type Server struct {
	s  *http.Server
	cl chan error
}

// New
func New(mux *mux.Router, port string) *Server {
	return &Server{
		s: &http.Server{
			Addr:    ":" + port,
			Handler: mux,
		},
		cl: make(chan error),
	}
}

// Start
func (s *Server) Start() {
	go func() {
		err := s.s.ListenAndServe()
		s.cl <- err
	}()
}

// Stop
func (s *Server) Stop() {
	s.s.Shutdown(context.Background())
}

// Ch return error chanel
func (s *Server) Ch() <-chan error {
	return s.cl
}
