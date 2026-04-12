package handlers

import (
	"encoding/json"
	"net/http"

	"agent-kanban-api/internal/domain"

	"github.com/jackc/pgx/v5/pgxpool"
)

type BoardHandler struct {
	DB *pgxpool.Pool
}

func NewBoardHandler(db *pgxpool.Pool) *BoardHandler {
	return &BoardHandler{DB: db}
}

func (h *BoardHandler) GetBoards(w http.ResponseWriter, r *http.Request) {
	rows, err := h.DB.Query(r.Context(), "SELECT id, name, description, created_at, updated_at FROM boards")
	if err != nil {
		http.Error(w, "Failed to query boards", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	boards := make([]domain.Board, 0)

	for rows.Next() {
		var board domain.Board

		err := rows.Scan(
			&board.ID,
			&board.Name,
			&board.Description,
			&board.CreatedAt,
			&board.UpdatedAt,
		)
		if err != nil {
			http.Error(w, "Failed to scan board", http.StatusInternalServerError)
			return
		}

		boards = append(boards, board)
	}

	writeJsonResponse(w, http.StatusOK, boards)
}

func (h *BoardHandler) CreateBoard(w http.ResponseWriter, r *http.Request) {
	type createBoardRequest struct {
		Name        string `json:"name"`
		Description string `json:"description"`
	}

	var req createBoardRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	var board domain.Board
	err := h.DB.QueryRow(
		r.Context(),
		"INSERT INTO boards (name, description) VALUES ($1, $2) RETURNING id, name, description, created_at, updated_at",
		req.Name,
		req.Description,
	).Scan(&board.ID, &board.Name, &board.Description, &board.CreatedAt, &board.UpdatedAt)

	if err != nil {
		http.Error(w, "Failed to create board", http.StatusInternalServerError)
		return
	}

	writeJsonResponse(w, http.StatusCreated, board)
}

func writeJsonResponse(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(data)
}
