package app

import (
	"net/http"

	"agent-kanban-api/internal/handlers"

	"github.com/go-chi/chi/v5"
)

func NewRouter() http.Handler {
	r := chi.NewRouter()

	r.Get("/health", handlers.HealthHandler)
	r.Get("/boards", handlers.ListBoardsHandler)

	return r
}
