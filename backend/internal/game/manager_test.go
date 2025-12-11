package game

import "testing"

func TestApplyMoveSwitchesTurn(t *testing.T) {
	gm := NewManager()
	g := gm.CreateGame("alice", "bob")

	updated, result, err := gm.ApplyMove(g.ID, "alice", 0)
	if err != nil {
		t.Fatalf("apply move failed: %v", err)
	}
	if result != CONTINUE {
		t.Fatalf("expected CONTINUE, got %v", result)
	}
	if updated.CurrentTurn != 2 {
		t.Fatalf("expected current turn 2, got %d", updated.CurrentTurn)
	}
	if updated.Board[Rows-1][0] != 1 {
		t.Fatalf("expected disc in column 0")
	}
}

func TestApplyMoveRejectsWrongTurn(t *testing.T) {
	gm := NewManager()
	g := gm.CreateGame("alice", "bob")

	_, result, err := gm.ApplyMove(g.ID, "bob", 0)
	if err == nil {
		t.Fatalf("expected error for wrong turn")
	}
	if result != INVALID {
		t.Fatalf("expected INVALID result, got %v", result)
	}
}

func TestApplyMoveDetectsWin(t *testing.T) {
	gm := NewManager()
	g := gm.CreateGame("alice", "bob")

	sequences := []struct {
		player string
		col    int
	}{
		{"alice", 0},
		{"bob", 1},
		{"alice", 0},
		{"bob", 1},
		{"alice", 0},
		{"bob", 1},
	}

	for _, move := range sequences {
		if _, _, err := gm.ApplyMove(g.ID, move.player, move.col); err != nil {
			t.Fatalf("setup move failed: %v", err)
		}
	}

	updated, result, err := gm.ApplyMove(g.ID, "alice", 0)
	if err != nil {
		t.Fatalf("winning move failed: %v", err)
	}
	if result != WIN {
		t.Fatalf("expected WIN, got %v", result)
	}
	if updated.Winner == nil || *updated.Winner != "alice" {
		t.Fatalf("expected winner alice")
	}
}

func TestApplyMoveDraw(t *testing.T) {
	gm := NewManager()
	g := gm.CreateGame("alice", "bob")

	pattern := [][]int{
		{0, 1, 2, 2, 1, 1, 2},
		{2, 2, 1, 1, 2, 2, 1},
		{1, 1, 2, 2, 1, 1, 2},
		{2, 2, 1, 1, 2, 2, 1},
		{1, 1, 2, 2, 1, 1, 2},
		{2, 2, 1, 1, 2, 2, 1},
	}

	clone := newBoard()
	for r := 0; r < Rows; r++ {
		copy(clone[r], pattern[r])
	}

	g.Board = clone
	g.CurrentTurn = 1

	updated, result, err := gm.ApplyMove(g.ID, "alice", 0)
	if err != nil {
		t.Fatalf("draw move failed: %v", err)
	}
	if result != DRAW {
		t.Fatalf("expected DRAW, got %v", result)
	}
	if !IsBoardFull(updated.Board) {
		t.Fatalf("expected board to be full after draw")
	}
	if updated.Winner != nil {
		t.Fatalf("draw should not set winner")
	}
}
