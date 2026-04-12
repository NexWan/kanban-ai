package app

import (
	"net/http"

	"agent-kanban-api/internal/handlers"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func NewRouter(db *pgxpool.Pool) http.Handler {
	r := chi.NewRouter()

	boardHandler := handlers.NewBoardHandler(db)

	// Define routes

	r.Get("/health", handlers.HealthHandler)
	r.Get("/boards", boardHandler.GetBoards)
	r.Post("/boards", boardHandler.CreateBoard)

	return r
}
