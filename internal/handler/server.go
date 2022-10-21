package handler

import (
	"context"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/Astemirdum/logs/internal/config"
)

type Server struct {
	s *http.Server
}

func NewServer(cfg config.Server, h *logHandler) *Server {
	s := &http.Server{
		Addr:        net.JoinHostPort(cfg.Host, cfg.Port),
		Handler:     h.newRouter(),
		ReadTimeout: 30 * time.Second,
	}
	return &Server{s: s}
}

func (s *Server) Run() {
	go func() {
		if err := s.s.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()
}

func (s *Server) Stop(ctx context.Context) error {
	return s.s.Shutdown(ctx)
}
