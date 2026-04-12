package handlers

import (
	"encoding/json"
	"net/http"

	"agent-kanban-api/internal/domain"
)

func ListBoardsHandler(w http.ResponseWriter, r *http.Request) {
	// For demonstration, we'll return a static list of boards.
	boards := []domain.Board{
		{
			ID:          "1",
			Name:        "Project Alpha",
			Description: "This is the first project board.",
			CreatedAt:   "2024-01-01T12:00:00Z",
			UpdatedAt:   "2024-01-02T12:00:00Z",
			AssignedTo:  "Leo",
		},
		{
			ID:          "2",
			Name:        "Project Beta",
			Description: "This is the second project board.",
			CreatedAt:   "2024-01-03T12:00:00Z",
			UpdatedAt:   "2024-01-04T12:00:00Z",
			AssignedTo:  "Diego",
		},
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(boards)
}
