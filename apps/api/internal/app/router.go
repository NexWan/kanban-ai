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
	// Board routes
	r.Get("/boards", boardHandler.GetBoards)
	r.Get("/boards/{id}", boardHandler.GetBoardByID)
	r.Post("/boards", boardHandler.CreateBoard)
	r.Put("/boards/{id}", boardHandler.UpdateBoard)
	r.Delete("/boards/{id}", boardHandler.DeleteBoard)
	// Card routes
	r.Get("/cards", cardHandler.GetCards)
	r.Get("/cards/{id}", cardHandler.GetCardByID)
	r.Post("/cards", cardHandler.CreateCard)
	r.Put("/cards/{id}", cardHandler.UpdateCard)
	r.Delete("/cards/{id}", cardHandler.DeleteCard)
	// Column routes
	r.Get("/columns", columnHandler.GetColumns)
	r.Get("/columns/{id}", columnHandler.GetColumnByID)
	r.Post("/columns", columnHandler.CreateColumn)
	r.Put("/columns/{id}", columnHandler.UpdateColumn)
	r.Delete("/columns/{id}", columnHandler.DeleteColumn)

	return r
}
