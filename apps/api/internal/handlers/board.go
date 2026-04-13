package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"agent-kanban-api/internal/domain"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type BoardHandler struct {
	DB *pgxpool.Pool
}

func NewBoardHandler(db *pgxpool.Pool) *BoardHandler {
	return &BoardHandler{DB: db}
}

type cardResponse struct {
	ID           string  `json:"id"`
	Title        string  `json:"title"`
	Description  string  `json:"description"`
	Priority     string  `json:"priority"`
	Status       string  `json:"status"`
	Position     int     `json:"position"`
	ParentCardID *string `json:"parent_card_id"`
}

type columnResponse struct {
	ID       string         `json:"id"`
	Name     string         `json:"name"`
	Position int            `json:"position"`
	Cards    []cardResponse `json:"cards"`
}

type boardResponse struct {
	ID          string           `json:"id"`
	Name        string           `json:"name"`
	Description string           `json:"description"`
	Columns     []columnResponse `json:"columns"`
}

func (h *BoardHandler) GetBoards(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// 1. Fetch all boards
	boardRows, err := h.DB.Query(ctx, "SELECT id, name, description FROM boards ORDER BY created_at")
	if err != nil {
		http.Error(w, "Failed to query boards", http.StatusInternalServerError)
		return
	}

	boardMap := map[string]*boardResponse{}
	boardOrder := []string{}

	for boardRows.Next() {
		b := new(boardResponse)
		if err := boardRows.Scan(&b.ID, &b.Name, &b.Description); err != nil {
			boardRows.Close()
			http.Error(w, "Failed to scan board", http.StatusInternalServerError)
			return
		}
		boardMap[b.ID] = b
		boardOrder = append(boardOrder, b.ID)
	}
	boardRows.Close()

	// 2. Fetch all columns
	colRows, err := h.DB.Query(ctx, "SELECT id, board_id, name, position FROM columns ORDER BY board_id, position")
	if err != nil {
		http.Error(w, "Failed to query columns", http.StatusInternalServerError)
		return
	}

	colMap := map[string]*columnResponse{}
	colOrder := map[string][]string{} // boardID -> ordered colIDs

	for colRows.Next() {
		col := &columnResponse{Cards: make([]cardResponse, 0)}
		var boardID string
		if err := colRows.Scan(&col.ID, &boardID, &col.Name, &col.Position); err != nil {
			colRows.Close()
			http.Error(w, "Failed to scan column", http.StatusInternalServerError)
			return
		}
		colMap[col.ID] = col
		colOrder[boardID] = append(colOrder[boardID], col.ID)
	}
	colRows.Close()

	// 3. Fetch all cards
	cardRows, err := h.DB.Query(ctx, "SELECT id, column_id, parent_card_id, title, description, priority, status, position FROM cards ORDER BY column_id, position")
	if err != nil {
		http.Error(w, "Failed to query cards", http.StatusInternalServerError)
		return
	}

	for cardRows.Next() {
		var card cardResponse
		var columnID string
		if err := cardRows.Scan(&card.ID, &columnID, &card.ParentCardID, &card.Title, &card.Description, &card.Priority, &card.Status, &card.Position); err != nil {
			cardRows.Close()
			http.Error(w, "Failed to scan card", http.StatusInternalServerError)
			return
		}
		if col, ok := colMap[columnID]; ok {
			col.Cards = append(col.Cards, card)
		}
	}
	cardRows.Close()

	// 4. Assemble final response in original board order
	boards := make([]boardResponse, 0, len(boardOrder))
	for _, boardID := range boardOrder {
		b := boardMap[boardID]
		b.Columns = make([]columnResponse, 0, len(colOrder[boardID]))
		for _, colID := range colOrder[boardID] {
			b.Columns = append(b.Columns, *colMap[colID])
		}
		boards = append(boards, *b)
	}

	writeJsonResponse(w, http.StatusOK, boards)
}

func (h *BoardHandler) GetBoardByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		http.Error(w, "Missing board ID", http.StatusBadRequest)
		return
	}

	ctx := r.Context()

	// 1. Fetch the board
	b := &boardResponse{}
	err := h.DB.QueryRow(ctx,
		"SELECT id, name, description FROM boards WHERE id = $1",
		id,
	).Scan(&b.ID, &b.Name, &b.Description)
	if errors.Is(err, pgx.ErrNoRows) {
		http.Error(w, "Board not found", http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, "Failed to fetch board", http.StatusInternalServerError)
		return
	}

	// 2. Fetch columns for this board
	colRows, err := h.DB.Query(ctx,
		"SELECT id, name, position FROM columns WHERE board_id = $1 ORDER BY position",
		id,
	)
	if err != nil {
		http.Error(w, "Failed to query columns", http.StatusInternalServerError)
		return
	}

	colMap := map[string]*columnResponse{}
	colOrder := []string{}

	for colRows.Next() {
		col := &columnResponse{Cards: make([]cardResponse, 0)}
		if err := colRows.Scan(&col.ID, &col.Name, &col.Position); err != nil {
			colRows.Close()
			http.Error(w, "Failed to scan column", http.StatusInternalServerError)
			return
		}
		colMap[col.ID] = col
		colOrder = append(colOrder, col.ID)
	}
	colRows.Close()

	// 3. Fetch cards for this board
	cardRows, err := h.DB.Query(ctx,
		"SELECT id, column_id, parent_card_id, title, description, priority, status, position FROM cards WHERE board_id = $1 ORDER BY column_id, position",
		id,
	)
	if err != nil {
		http.Error(w, "Failed to query cards", http.StatusInternalServerError)
		return
	}

	for cardRows.Next() {
		var card cardResponse
		var columnID string
		if err := cardRows.Scan(&card.ID, &columnID, &card.ParentCardID, &card.Title, &card.Description, &card.Priority, &card.Status, &card.Position); err != nil {
			cardRows.Close()
			http.Error(w, "Failed to scan card", http.StatusInternalServerError)
			return
		}
		if col, ok := colMap[columnID]; ok {
			col.Cards = append(col.Cards, card)
		}
	}
	cardRows.Close()

	// 4. Assemble columns in order
	b.Columns = make([]columnResponse, 0, len(colOrder))
	for _, colID := range colOrder {
		b.Columns = append(b.Columns, *colMap[colID])
	}

	writeJsonResponse(w, http.StatusOK, b)
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

func (h *BoardHandler) DeleteBoard(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		http.Error(w, "Missing board ID", http.StatusBadRequest)
		return
	}

	commandTag, err := h.DB.Exec(
		r.Context(),
		"DELETE FROM boards WHERE id = $1",
		id,
	)
	if err != nil {
		http.Error(w, "Failed to delete board", http.StatusInternalServerError)
		return
	}

	if commandTag.RowsAffected() == 0 {
		http.Error(w, "Board not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *BoardHandler) UpdateBoard(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		http.Error(w, "Missing board ID", http.StatusBadRequest)
		return
	}

	type updateBoardRequest struct {
		Name        string `json:"name"`
		Description string `json:"description"`
	}

	var req updateBoardRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	commandTag, err := h.DB.Exec(
		r.Context(),
		"UPDATE boards SET name = $1, description = $2, updated_at = NOW() WHERE id = $3",
		req.Name,
		req.Description,
		id,
	)
	if err != nil {
		http.Error(w, "Failed to update board", http.StatusInternalServerError)
		return
	}

	if commandTag.RowsAffected() == 0 {
		http.Error(w, "Board not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func writeJsonResponse(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(data)
}
