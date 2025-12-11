package types

// ClientMessage represents messages sent from the frontend to the server.
type ClientMessage struct {
	Type     string `json:"type"`
	Col      *int   `json:"col,omitempty"`
	GameID   string `json:"gameId,omitempty"`
	Username string `json:"username,omitempty"`
}

// ServerMessage mirrors the frontend contract.
type ServerMessage struct {
	Type        string  `json:"type"`
	GameID      string  `json:"gameId,omitempty"`
	You         int     `json:"you,omitempty"`
	Opponent    string  `json:"opponent,omitempty"`
	Board       [][]int `json:"board,omitempty"`
	CurrentTurn int     `json:"currentTurn,omitempty"`
	Result      string  `json:"result,omitempty"`
	Message     string  `json:"message,omitempty"`
}
