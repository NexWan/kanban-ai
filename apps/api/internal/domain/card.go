package domain

import "time"

type Card struct {
	Id             string    `json:"id"`
	BoardId        string    `json:"board_id"`
	ColumnId       string    `json:"column_id"`
	ParentCardId   *string   `json:"parent_card_id,omitempty"`
	Title          string    `json:"title"`
	Description    string    `json:"description"`
	Priority       string    `json:"priority"`
	Status         string    `json:"status"`
	Position       int       `json:"position"`
	AssignedUserId *string   `json:"assigned_user_id,omitempty"`
	AgentOwnerId   *string   `json:"agent_owner_id,omitempty"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}
