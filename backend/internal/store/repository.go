package store

import (
	"database/sql"
	"encoding/json"
	"errors"
	"time"
)

// Repository exposes persistence helpers for leaderboards and completed games.
type Repository struct {
	db *sql.DB
}

// LeaderboardEntry represents a leaderboard row.
type LeaderboardEntry struct {
	Username string `json:"username"`
	Wins     int    `json:"wins"`
}

// CompletedMove describes a single move in a completed game.
type CompletedMove struct {
	Player     string `json:"player"`
	Column     int    `json:"column"`
	MoveNumber int    `json:"moveNumber"`
}

// CompletedGame captures the data required to persist a finished match.
type CompletedGame struct {
	ID        string
	Player1   string
	Player2   string
	Winner    *string
	IsDraw    bool
	Moves     []CompletedMove
	StartedAt time.Time
	EndedAt   time.Time
}

// NewRepository constructs a Repository using an existing sql.DB connection.
func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

// EnsurePlayer guarantees that a player row exists.
func (r *Repository) EnsurePlayer(username string) error {
	if username == "" {
		return errors.New("username is required")
	}

	_, err := r.db.Exec(`INSERT INTO players (username) VALUES ($1) ON CONFLICT (username) DO NOTHING`, username)
	return err
}

// IncrementWin increments the win counter for a player.
func (r *Repository) IncrementWin(username string) error {
	if username == "" {
		return errors.New("username is required")
	}

	res, err := r.db.Exec(`UPDATE players SET wins = wins + 1 WHERE username = $1`, username)
	if err != nil {
		return err
	}
	if affected, err := res.RowsAffected(); err == nil && affected == 0 {
		// Player row missing, ensure and retry.
		if err := r.EnsurePlayer(username); err != nil {
			return err
		}
		_, err = r.db.Exec(`UPDATE players SET wins = wins + 1 WHERE username = $1`, username)
	}
	return err
}

// SaveCompletedGame inserts a completed game record.
func (r *Repository) SaveCompletedGame(game *CompletedGame) error {
	if game == nil {
		return errors.New("game is required")
	}
	if game.ID == "" {
		return errors.New("game id is required")
	}
	if game.Player1 == "" || game.Player2 == "" {
		return errors.New("player names are required")
	}

	movesJSON, err := json.Marshal(game.Moves)
	if err != nil {
		return err
	}

	winner := sql.NullString{}
	if game.Winner != nil && *game.Winner != "" {
		winner.Valid = true
		winner.String = *game.Winner
	}

	startedAt := game.StartedAt
	if startedAt.IsZero() {
		startedAt = time.Now().UTC()
	}

	endedAt := game.EndedAt
	if endedAt.IsZero() {
		endedAt = time.Now().UTC()
	}

	_, err = r.db.Exec(
		`INSERT INTO games (id, player1, player2, winner, is_draw, moves, started_at, ended_at)
         VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`,
		game.ID,
		game.Player1,
		game.Player2,
		winner,
		game.IsDraw,
		movesJSON,
		startedAt,
		endedAt,
	)
	return err
}

// GetLeaderboard returns the top players ordered by wins.
func (r *Repository) GetLeaderboard(limit int) ([]LeaderboardEntry, error) {
	if limit <= 0 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}

	rows, err := r.db.Query(`SELECT username, wins FROM players ORDER BY wins DESC, username ASC LIMIT $1`, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	entries := make([]LeaderboardEntry, 0, limit)
	for rows.Next() {
		var entry LeaderboardEntry
		if err := rows.Scan(&entry.Username, &entry.Wins); err != nil {
			return nil, err
		}
		entries = append(entries, entry)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return entries, nil
}
