package handlers

import (
	"agent-kanban-api/internal/domain"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5"
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
	if err := rows.Err(); err != nil {
		http.Error(w, "Failed to read columns", http.StatusInternalServerError)
		return
	}

	writeJsonResponse(w, http.StatusOK, columns)
}

func (h *ColumnHandler) GetColumnByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		http.Error(w, "Missing column ID", http.StatusBadRequest)
		return
	}

	var column domain.Column
	err := h.DB.QueryRow(r.Context(),
		"SELECT id, board_id, name, position, created_at, updated_at FROM columns WHERE id = $1",
		id,
	).Scan(
		&column.Id,
		&column.BoardId,
		&column.Name,
		&column.Position,
		&column.CreatedAt,
		&column.UpdatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		http.Error(w, "Column not found", http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, "Failed to fetch column", http.StatusInternalServerError)
		return
	}

	writeJsonResponse(w, http.StatusOK, column)
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

func (h *ColumnHandler) UpdateColumn(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		http.Error(w, "Missing column ID", http.StatusBadRequest)
		return
	}

	type updateColumnRequest struct {
		Name     string `json:"name"`
		Position int    `json:"position"`
	}

	var req updateColumnRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	commandTag, err := h.DB.Exec(r.Context(),
		"UPDATE columns SET name = $1, position = $2, updated_at = NOW() WHERE id = $3",
		req.Name,
		req.Position,
		id,
	)
	if err != nil {
		http.Error(w, "Failed to update column", http.StatusInternalServerError)
		return
	}

	if commandTag.RowsAffected() == 0 {
		http.Error(w, "Column not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *ColumnHandler) DeleteColumn(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		http.Error(w, "Missing column ID", http.StatusBadRequest)
		return
	}

	commandTag, err := h.DB.Exec(r.Context(),
		"DELETE FROM columns WHERE id = $1",
		id,
	)
	if err != nil {
		http.Error(w, "Failed to delete column", http.StatusInternalServerError)
		return
	}

	if commandTag.RowsAffected() == 0 {
		http.Error(w, "Column not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
