package api

import (
	"github.com/go-chi/chi"

	"github.com/Quantum12k/healthcheck-service/internal/api/handlers"
)

func (s *Server) getRouter() *chi.Mux {
	r := chi.NewRouter()

	r.Get("/get_status", handlers.GetStatus(s.appCache.LastChecks))

	return r
}


