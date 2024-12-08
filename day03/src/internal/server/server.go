package server

import (
	"context"
	"net/http"
)

type Server struct {
	server http.Server
}

func New(addr string, h http.Handler) *Server {
	return &Server{
		server: http.Server{
			Addr:    addr,
			Handler: h,
		},
	}
}

func (s *Server) Start() error {
	return s.server.ListenAndServe()
}

func (s *Server) Stop() error {
	return s.server.Shutdown(context.Background())
}
