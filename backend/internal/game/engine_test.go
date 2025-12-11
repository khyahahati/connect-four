package game

import "testing"

func TestDropDiscPlacesDisc(t *testing.T) {
	board := newBoard()

	updated, row, err := DropDisc(board, 3, 1)
	if err != nil {
		t.Fatalf("drop disc failed: %v", err)
	}
	if row != Rows-1 {
		t.Fatalf("expected row %d, got %d", Rows-1, row)
	}
	if updated[Rows-1][3] != 1 {
		t.Fatalf("expected player disc at bottom")
	}
	if board[Rows-1][3] != 0 {
		t.Fatalf("original board should remain unchanged")
	}
}

func TestDropDiscFullColumnError(t *testing.T) {
	board := newBoard()
	for i := 0; i < Rows; i++ {
		var err error
		board, _, err = DropDisc(board, 0, 1+(i%2))
		if err != nil {
			t.Fatalf("unexpected error filling column: %v", err)
		}
	}
	if _, _, err := DropDisc(board, 0, 1); err == nil {
		t.Fatalf("expected error when dropping into full column")
	}
}

func TestCheckWinDetections(t *testing.T) {
	board := newBoard()

	// Horizontal win.
	board[Rows-1][0], board[Rows-1][1], board[Rows-1][2], board[Rows-1][3] = 1, 1, 1, 1
	if !CheckWin(board, 1) {
		t.Fatalf("expected horizontal win")
	}

	board = newBoard()
	// Vertical win.
	for i := 0; i < 4; i++ {
		board[Rows-1-i][2] = 2
	}
	if !CheckWin(board, 2) {
		t.Fatalf("expected vertical win")
	}

	board = newBoard()
	// Diagonal ↘ win.
	for i := 0; i < 4; i++ {
		board[Rows-1-i][i] = 1
	}
	if !CheckWin(board, 1) {
		t.Fatalf("expected diagonal ↘ win")
	}

	board = newBoard()
	// Diagonal ↗ win.
	for i := 0; i < 4; i++ {
		board[Rows-1-i][3-i] = 2
	}
	if !CheckWin(board, 2) {
		t.Fatalf("expected diagonal ↗ win")
	}

	if CheckWin(newBoard(), 1) {
		t.Fatalf("did not expect win on empty board")
	}
}

func TestIsBoardFull(t *testing.T) {
	if IsBoardFull(newBoard()) {
		t.Fatalf("empty board should not be full")
	}

	board := newBoard()
	for r := 0; r < Rows; r++ {
		for c := 0; c < Columns; c++ {
			board[r][c] = 1
		}
	}
	if !IsBoardFull(board) {
		t.Fatalf("expected board to be full")
	}
}
