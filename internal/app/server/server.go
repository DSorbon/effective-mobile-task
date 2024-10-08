package server

import (
	"context"
	"net/http"
	"time"

	"github.com/DSorbon/effective-mobile-task/internal/config"
)

type Server struct {
	httpServer *http.Server
}

func NewServer() *Server {
	return &Server{}
}

func (s *Server) Run(handler http.Handler) error {
	s.httpServer = &http.Server{
		Addr:           ":" + config.Values.APIPort,
		Handler:        handler,
		MaxHeaderBytes: config.Values.MaxHeaderBytes << 20, // MB
		ReadTimeout:    time.Duration(config.Values.ReadTimeout) * time.Second,
		WriteTimeout:   time.Duration(config.Values.WriteTimeout) * time.Second,
	}

	return s.httpServer.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}
