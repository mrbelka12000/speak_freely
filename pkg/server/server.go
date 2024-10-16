package server

import (
	"context"
	"net/http"
)

type Server struct {
	s  *http.Server
	cl chan error
}

func New(mux *http.ServeMux, port string) *Server {
	return &Server{
		s: &http.Server{
			Addr:    ":" + port,
			Handler: mux,
		},
		cl: make(chan error),
	}
}

func (s *Server) Start() {
	go func() {
		err := s.s.ListenAndServe()
		s.cl <- err
	}()
}

func (s *Server) Stop() {
	s.s.Shutdown(context.Background())
}

func (s *Server) Ch() <-chan error {
	return s.cl
}
