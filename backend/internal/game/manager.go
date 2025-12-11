package game

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
)

// Game represents an active match between two players.
type Game struct {
	ID          string
	Player1     string
	Player2     string
	Board       [][]int
	CurrentTurn int
	CreatedAt   time.Time
	Winner      *string
}

// GameManager creates and stores game sessions.
type GameManager struct {
	mu    sync.RWMutex
	games map[string]*Game
}

// NewManager returns a ready-to-use GameManager.
func NewManager() *GameManager {
	return &GameManager{
		games: make(map[string]*Game),
	}
}

// CreateGame registers a new game for two participants.
func (m *GameManager) CreateGame(player1, player2 string) *Game {
	game := &Game{
		ID:          uuid.NewString(),
		Player1:     player1,
		Player2:     player2,
		Board:       newBoard(),
		CurrentTurn: 1,
		CreatedAt:   time.Now().UTC(),
	}

	m.mu.Lock()
	m.games[game.ID] = game
	m.mu.Unlock()

	return game
}

// GetGame retrieves a game by its identifier.
func (m *GameManager) GetGame(id string) (*Game, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	game, ok := m.games[id]
	return game, ok
}

// RemoveGame deletes a game from the manager.
func (m *GameManager) RemoveGame(id string) {
	m.mu.Lock()
	delete(m.games, id)
	m.mu.Unlock()
}

// MoveResult indicates the outcome of applying a move to the game state.
type MoveResult int

const (
	CONTINUE MoveResult = iota
	WIN
	DRAW
	INVALID
)

// ApplyMove validates and applies a move for the given player and column.
func (m *GameManager) ApplyMove(gameID string, player string, col int) (*Game, MoveResult, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	game, ok := m.games[gameID]
	if !ok {
		return nil, INVALID, fmt.Errorf("game %s not found", gameID)
	}

	if game.Winner != nil {
		return game, INVALID, errors.New("game already finished")
	}

	var playerNum int
	switch player {
	case game.Player1:
		playerNum = 1
	case game.Player2:
		playerNum = 2
	default:
		return game, INVALID, errors.New("player not part of this game")
	}

	if game.CurrentTurn != playerNum {
		return game, INVALID, errors.New("not your turn")
	}

	newBoard, _, err := DropDisc(game.Board, col, playerNum)
	if err != nil {
		return game, INVALID, err
	}

	game.Board = newBoard
	game.Winner = nil

	if CheckWin(newBoard, playerNum) {
		winner := player
		game.Winner = &winner
		return game, WIN, nil
	}

	if IsBoardFull(newBoard) {
		return game, DRAW, nil
	}

	if playerNum == 1 {
		game.CurrentTurn = 2
	} else {
		game.CurrentTurn = 1
	}

	return game, CONTINUE, nil
}

func newBoard() [][]int {
	board := make([][]int, Rows)
	for r := 0; r < Rows; r++ {
		board[r] = make([]int, Columns)
	}
	return board
}

// FindGameByPlayers searches for a game containing both players, regardless of order.
func (m *GameManager) FindGameByPlayers(playerA, playerB string) (*Game, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	for _, game := range m.games {
		if (game.Player1 == playerA && game.Player2 == playerB) || (game.Player1 == playerB && game.Player2 == playerA) {
			return game, true
		}
	}

	return nil, false
}
