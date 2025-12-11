export type PlayerId = 1 | 2;

export function cloneBoard(board: number[][]): number[][] {
  return board.map((row) => [...row]);
}

export function findAvailableRow(board: number[][], columnIndex: number): number {
  for (let row = board.length - 1; row >= 0; row -= 1) {
    if (board[row][columnIndex] === 0) {
      return row;
    }
  }
  return -1;
}

export function dropDisc(board: number[][], columnIndex: number, player: PlayerId): { nextBoard: number[][]; row: number } | null {
  const targetRow = findAvailableRow(board, columnIndex);
  if (targetRow === -1) {
    return null;
  }

  const nextBoard = cloneBoard(board);
  nextBoard[targetRow][columnIndex] = player;
  return { nextBoard, row: targetRow };
}

export function checkWin(board: number[][], player: PlayerId): boolean {
  const rows = board.length;
  const columns = board[0].length;

  for (let row = 0; row < rows; row += 1) {
    for (let col = 0; col <= columns - 4; col += 1) {
      if (
        board[row][col] === player &&
        board[row][col + 1] === player &&
        board[row][col + 2] === player &&
        board[row][col + 3] === player
      ) {
        return true;
      }
    }
  }

  for (let col = 0; col < columns; col += 1) {
    for (let row = 0; row <= rows - 4; row += 1) {
      if (
        board[row][col] === player &&
        board[row + 1][col] === player &&
        board[row + 2][col] === player &&
        board[row + 3][col] === player
      ) {
        return true;
      }
    }
  }

  for (let row = 0; row <= rows - 4; row += 1) {
    for (let col = 0; col <= columns - 4; col += 1) {
      if (
        board[row][col] === player &&
        board[row + 1][col + 1] === player &&
        board[row + 2][col + 2] === player &&
        board[row + 3][col + 3] === player
      ) {
        return true;
      }
    }
  }

  for (let row = 3; row < rows; row += 1) {
    for (let col = 0; col <= columns - 4; col += 1) {
      if (
        board[row][col] === player &&
        board[row - 1][col + 1] === player &&
        board[row - 2][col + 2] === player &&
        board[row - 3][col + 3] === player
      ) {
        return true;
      }
    }
  }

  return false;
}

export function isBoardFull(board: number[][]): boolean {
  return board.every((row) => row.every((cell) => cell !== 0));
}
