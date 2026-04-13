package handlers

import (
	"agent-kanban-api/internal/domain"
	"encoding/json"
	"log"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
)

type CardHandler struct {
	DB *pgxpool.Pool
}

func NewCardHandler(db *pgxpool.Pool) *CardHandler {
	return &CardHandler{DB: db}
}

func (h *CardHandler) GetCards(w http.ResponseWriter, r *http.Request) {
	rows, err := h.DB.Query(r.Context(), "SELECT id, board_id, column_id, title, description, status, priority, position, created_at, updated_at FROM cards")
	if err != nil {
		http.Error(w, "Failed to query cards", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	cards := make([]domain.Card, 0)

	for rows.Next() {
		var card domain.Card

		err := rows.Scan(
			&card.Id,
			&card.BoardId,
			&card.ColumnId,
			&card.Title,
			&card.Description,
			&card.Status,
			&card.Priority,
			&card.Position,
			&card.CreatedAt,
			&card.UpdatedAt,
		)
		if err != nil {
			http.Error(w, "Failed to scan card", http.StatusInternalServerError)
			return
		}

		cards = append(cards, card)
	}

	writeJsonResponse(w, http.StatusOK, cards)
}

func (h *CardHandler) CreateCard(w http.ResponseWriter, r *http.Request) {
	type createCardRequest struct {
		BoardID     string `json:"board_id"`
		ColumnID    string `json:"column_id"`
		Title       string `json:"title"`
		Description string `json:"description"`
		Priority    string `json:"priority"`
		Status      string `json:"status"`
		Position    int    `json:"position"`
	}

	var req createCardRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	var cardID string
	err := h.DB.QueryRow(
		r.Context(),
		"INSERT INTO cards (board_id, column_id, title, description, status, priority, position) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id",
		req.BoardID,
		req.ColumnID,
		req.Title,
		req.Description,
		req.Status,
		req.Priority,
		req.Position,
	).Scan(&cardID)

	if err != nil {
		http.Error(w, "Failed to create card", http.StatusInternalServerError)
		log.Printf("Error inserting card: %v", err)
		return
	}

	writeJsonResponse(w, http.StatusCreated, map[string]string{"id": cardID})
}
