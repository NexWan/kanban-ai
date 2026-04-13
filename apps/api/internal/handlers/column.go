package handlers

import (
	"agent-kanban-api/internal/domain"
	"encoding/json"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
)

type ColumnHandler struct {
	DB *pgxpool.Pool
}

func NewColumnHandler(db *pgxpool.Pool) *ColumnHandler {
	return &ColumnHandler{DB: db}
}

func (h *ColumnHandler) GetColumns(w http.ResponseWriter, r *http.Request) {
	rows, err := h.DB.Query(r.Context(), "SELECT id, board_id, name, position, created_at, updated_at FROM columns")
	if err != nil {
		http.Error(w, "Failed to query columns", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	columns := make([]domain.Column, 0)

	for rows.Next() {
		var column domain.Column

		err := rows.Scan(
			&column.Id,
			&column.BoardId,
			&column.Name,
			&column.Position,
			&column.CreatedAt,
			&column.UpdatedAt,
		)
		if err != nil {
			http.Error(w, "Failed to scan column", http.StatusInternalServerError)
			return
		}

		columns = append(columns, column)
	}

	writeJsonResponse(w, http.StatusOK, columns)
}

func (h *ColumnHandler) CreateColumn(w http.ResponseWriter, r *http.Request) {
	type createColumnRequest struct {
		BoardID  string `json:"board_id"`
		Name     string `json:"name"`
		Position int    `json:"position"`
	}

	var req createColumnRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	var columnId string
	err := h.DB.QueryRow(
		r.Context(),
		"INSERT INTO columns (board_id, name, position) VALUES ($1, $2, $3) RETURNING id",
		req.BoardID,
		req.Name,
		req.Position,
	).Scan(&columnId)

	if err != nil {
		http.Error(w, "Failed to create column", http.StatusInternalServerError)
		return
	}

	writeJsonResponse(w, http.StatusCreated, map[string]string{"id": columnId})
}
