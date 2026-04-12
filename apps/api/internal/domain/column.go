package domain

type Domain struct {
	Id        string `json:"id"`
	BoardId   string `json:"board_id"`
	Name      string `json:"name"`
	Position  int    `json:"position"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}
