package handlers

import (
	"agent-kanban-api/internal/domain"
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type CardHandler struct {
	DB *pgxpool.Pool
}

func NewCardHandler(db *pgxpool.Pool) *CardHandler {
	return &CardHandler{DB: db}
}

func (h *CardHandler) GetCards(w http.ResponseWriter, r *http.Request) {
	rows, err := h.DB.Query(r.Context(), "SELECT id, board_id, column_id, parent_card_id, title, description, status, priority, position, assigned_user_id, agent_owner_id, created_at, updated_at FROM cards")
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
			&card.ParentCardId,
			&card.Title,
			&card.Description,
			&card.Status,
			&card.Priority,
			&card.Position,
			&card.AssignedUserId,
			&card.AgentOwnerId,
			&card.CreatedAt,
			&card.UpdatedAt,
		)
		if err != nil {
			http.Error(w, "Failed to scan card", http.StatusInternalServerError)
			return
		}

		cards = append(cards, card)
	}
	if err := rows.Err(); err != nil {
		http.Error(w, "Failed to read cards", http.StatusInternalServerError)
		return
	}

	writeJsonResponse(w, http.StatusOK, cards)
}

func (h *CardHandler) GetCardByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		http.Error(w, "Missing card ID", http.StatusBadRequest)
		return
	}

	var card domain.Card
	err := h.DB.QueryRow(r.Context(),
		"SELECT id, board_id, column_id, parent_card_id, title, description, status, priority, position, assigned_user_id, agent_owner_id, created_at, updated_at FROM cards WHERE id = $1",
		id,
	).Scan(
		&card.Id,
		&card.BoardId,
		&card.ColumnId,
		&card.ParentCardId,
		&card.Title,
		&card.Description,
		&card.Status,
		&card.Priority,
		&card.Position,
		&card.AssignedUserId,
		&card.AgentOwnerId,
		&card.CreatedAt,
		&card.UpdatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		http.Error(w, "Card not found", http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, "Failed to fetch card", http.StatusInternalServerError)
		return
	}

	writeJsonResponse(w, http.StatusOK, card)
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

func (h *CardHandler) UpdateCard(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		http.Error(w, "Missing card ID", http.StatusBadRequest)
		return
	}

	type updateCardRequest struct {
		ColumnID    string `json:"column_id"`
		Title       string `json:"title"`
		Description string `json:"description"`
		Priority    string `json:"priority"`
		Status      string `json:"status"`
		Position    int    `json:"position"`
	}

	var req updateCardRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	commandTag, err := h.DB.Exec(r.Context(),
		"UPDATE cards SET column_id = $1, title = $2, description = $3, priority = $4, status = $5, position = $6, updated_at = NOW() WHERE id = $7",
		req.ColumnID,
		req.Title,
		req.Description,
		req.Priority,
		req.Status,
		req.Position,
		id,
	)
	if err != nil {
		http.Error(w, "Failed to update card", http.StatusInternalServerError)
		return
	}

	if commandTag.RowsAffected() == 0 {
		http.Error(w, "Card not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *CardHandler) DeleteCard(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		http.Error(w, "Missing card ID", http.StatusBadRequest)
		return
	}

	commandTag, err := h.DB.Exec(r.Context(),
		"DELETE FROM cards WHERE id = $1",
		id,
	)
	if err != nil {
		http.Error(w, "Failed to delete card", http.StatusInternalServerError)
		return
	}

	if commandTag.RowsAffected() == 0 {
		http.Error(w, "Card not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
