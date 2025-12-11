package bot

import (
	"errors"

	"github.com/example/connect-four/backend/internal/game"
)

var preferenceOrder = []int{3, 2, 4, 1, 5, 0, 6}

// Bot encapsulates simple heuristics to play as the second player.
type Bot struct {
	gm *game.GameManager
}

// New creates a new bot helper bound to the provided game manager.
func New(gm *game.GameManager) *Bot {
	return &Bot{gm: gm}
}

// TakeTurn selects a column and applies the move for the bot. It returns the updated game, the move result, and the column played.
func (b *Bot) TakeTurn(gameID string) (*game.Game, game.MoveResult, int, error) {
	if b == nil || b.gm == nil {
		return nil, game.INVALID, -1, errors.New("bot not configured")
	}

	current, ok := b.gm.GetGame(gameID)
	if !ok {
		return nil, game.INVALID, -1, errors.New("game not found")
	}

	if current.Player2 != "BOT" {
		return current, game.INVALID, -1, errors.New("game is not against bot")
	}

	col, err := chooseColumn(current.Board)
	if err != nil {
		return current, game.INVALID, -1, err
	}

	updated, result, err := b.gm.ApplyMove(gameID, current.Player2, col)
	return updated, result, col, err
}

func chooseColumn(board [][]int) (int, error) {
	// Winning move.
	for col := 0; col < game.Columns; col++ {
		newBoard, _, err := game.DropDisc(board, col, 2)
		if err != nil {
			continue
		}
		if game.CheckWin(newBoard, 2) {
			return col, nil
		}
	}

	// Block opponent win.
	for col := 0; col < game.Columns; col++ {
		newBoard, _, err := game.DropDisc(board, col, 1)
		if err != nil {
			continue
		}
		if game.CheckWin(newBoard, 1) {
			return col, nil
		}
	}

	// Preference order.
	for _, col := range preferenceOrder {
		if _, _, err := game.DropDisc(board, col, 2); err == nil {
			return col, nil
		}
	}

	return 0, errors.New("bot has no valid moves")
}
