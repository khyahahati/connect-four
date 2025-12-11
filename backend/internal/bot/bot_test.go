package bot

import (
	"testing"

	"github.com/example/connect-four/backend/internal/game"
)

func TestBotTakesWinningMove(t *testing.T) {
	gm := game.NewManager()
	g := gm.CreateGame("human", "BOT")

	board := make([][]int, game.Rows)
	for r := range board {
		board[r] = make([]int, game.Columns)
	}
	board[game.Rows-1][0] = 2
	board[game.Rows-1][1] = 2
	board[game.Rows-1][2] = 2

	g.Board = board
	g.CurrentTurn = 2

	botEngine := New(gm)
	updated, result, col, err := botEngine.TakeTurn(g.ID)
	if err != nil {
		t.Fatalf("bot move failed: %v", err)
	}
	if result != game.WIN {
		t.Fatalf("expected WIN, got %v", result)
	}
	if col != 3 {
		t.Fatalf("expected bot to play column 3, got %d", col)
	}
	if updated.Board[game.Rows-1][3] != 2 {
		t.Fatalf("expected bot disc at winning slot")
	}
	if updated.Winner == nil || *updated.Winner != "BOT" {
		t.Fatalf("expected BOT to be recorded as winner")
	}
}

func TestBotBlocksOpponentWin(t *testing.T) {
	gm := game.NewManager()
	g := gm.CreateGame("human", "BOT")

	board := make([][]int, game.Rows)
	for r := range board {
		board[r] = make([]int, game.Columns)
	}
	board[game.Rows-1][0] = 1
	board[game.Rows-1][1] = 1
	board[game.Rows-1][2] = 1

	g.Board = board
	g.CurrentTurn = 2

	botEngine := New(gm)
	updated, result, col, err := botEngine.TakeTurn(g.ID)
	if err != nil {
		t.Fatalf("bot move failed: %v", err)
	}
	if result != game.CONTINUE {
		t.Fatalf("expected CONTINUE, got %v", result)
	}
	if col != 3 {
		t.Fatalf("expected bot to block column 3, got %d", col)
	}
	if updated.Board[game.Rows-1][3] != 2 {
		t.Fatalf("expected bot disc at block position")
	}
	if updated.CurrentTurn != 1 {
		t.Fatalf("expected turn to switch back to player 1")
	}
}

func TestBotPrefersCenter(t *testing.T) {
	gm := game.NewManager()
	g := gm.CreateGame("human", "BOT")
	g.CurrentTurn = 2

	botEngine := New(gm)
	updated, result, col, err := botEngine.TakeTurn(g.ID)
	if err != nil {
		t.Fatalf("bot move failed: %v", err)
	}
	if result != game.CONTINUE {
		t.Fatalf("expected CONTINUE, got %v", result)
	}
	if col != 3 {
		t.Fatalf("expected center column 3, got %d", col)
	}
	if updated.Board[game.Rows-1][3] != 2 {
		t.Fatalf("expected bot disc at center")
	}
	if updated.CurrentTurn != 1 {
		t.Fatalf("expected turn to return to player 1")
	}
}
