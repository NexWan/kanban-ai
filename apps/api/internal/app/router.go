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
	cardHandler := handlers.NewCardHandler(db)
	columnHandler := handlers.NewColumnHandler(db)

	// Define routes

	r.Get("/health", handlers.HealthHandler)
	r.Get("/boards", boardHandler.GetBoards)
	r.Post("/boards", boardHandler.CreateBoard)
	r.Get("/cards", cardHandler.GetCards)
	r.Post("/cards", cardHandler.CreateCard)
	r.Get("/columns", columnHandler.GetColumns)
	r.Post("/columns", columnHandler.CreateColumn)

	return r
}
