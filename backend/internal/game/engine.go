package game

import "errors"

const (
	Rows    = 6
	Columns = 7
)

// DropDisc returns a new board state after placing a disc for player in the specified column.
func DropDisc(board [][]int, col int, player int) ([][]int, int, error) {
	if col < 0 || col >= Columns {
		return nil, -1, errors.New("invalid column")
	}
	if player != 1 && player != 2 {
		return nil, -1, errors.New("invalid player")
	}

	newBoard := make([][]int, Rows)
	for r := 0; r < Rows; r++ {
		if len(board) > r {
			newBoard[r] = append([]int(nil), board[r]...)
		} else {
			newBoard[r] = make([]int, Columns)
		}
	}

	for row := Rows - 1; row >= 0; row-- {
		if newBoard[row][col] == 0 {
			newBoard[row][col] = player
			return newBoard, row, nil
		}
	}

	return nil, -1, errors.New("column is full")
}

// CheckWin determines whether the specified player has a connect-four on the board.
func CheckWin(board [][]int, player int) bool {
	// Horizontal
	for r := 0; r < Rows; r++ {
		for c := 0; c <= Columns-4; c++ {
			if board[r][c] == player && board[r][c+1] == player && board[r][c+2] == player && board[r][c+3] == player {
				return true
			}
		}
	}

	// Vertical
	for c := 0; c < Columns; c++ {
		for r := 0; r <= Rows-4; r++ {
			if board[r][c] == player && board[r+1][c] == player && board[r+2][c] == player && board[r+3][c] == player {
				return true
			}
		}
	}

	// Diagonal ↘
	for r := 0; r <= Rows-4; r++ {
		for c := 0; c <= Columns-4; c++ {
			if board[r][c] == player && board[r+1][c+1] == player && board[r+2][c+2] == player && board[r+3][c+3] == player {
				return true
			}
		}
	}

	// Diagonal ↗
	for r := 3; r < Rows; r++ {
		for c := 0; c <= Columns-4; c++ {
			if board[r][c] == player && board[r-1][c+1] == player && board[r-2][c+2] == player && board[r-3][c+3] == player {
				return true
			}
		}
	}

	return false
}

// IsBoardFull reports whether all board slots are occupied.
func IsBoardFull(board [][]int) bool {
	for r := 0; r < Rows; r++ {
		if len(board) <= r {
			return false
		}
		for c := 0; c < Columns; c++ {
			if board[r][c] == 0 {
				return false
			}
		}
	}
	return true
}
