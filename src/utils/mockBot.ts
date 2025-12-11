import { checkWin, cloneBoard, dropDisc, findAvailableRow, isBoardFull } from './board';

export const MOCK_BOT_ENABLED = true;
export const MOCK_BOT_PLAYER: 2 = 2;
export const MOCK_BOT_DELAY_MS = 900;
export const MOCK_BOT_NAME = 'Backend Bot';

export type MockBotMoveStatus = 'WIN' | 'LOSS' | 'DRAW' | 'CONTINUE' | 'BLOCKED';

export interface MockBotMoveResult {
  status: MockBotMoveStatus;
  board: number[][];
}

export function selectMockBotColumn(board: number[][]): number | null {
  const availableColumns: number[] = [];
  for (let col = 0; col < board[0].length; col += 1) {
    if (findAvailableRow(board, col) !== -1) {
      availableColumns.push(col);
    }
  }

  if (availableColumns.length === 0) {
    return null;
  }

  for (const column of availableColumns) {
    const simulated = cloneBoard(board);
    const row = findAvailableRow(simulated, column);
    if (row === -1) {
      continue;
    }
    simulated[row][column] = MOCK_BOT_PLAYER;
    if (checkWin(simulated, MOCK_BOT_PLAYER)) {
      return column;
    }
  }

  const opponent = 1;
  for (const column of availableColumns) {
    const simulated = cloneBoard(board);
    const row = findAvailableRow(simulated, column);
    if (row === -1) {
      continue;
    }
    simulated[row][column] = opponent;
    if (checkWin(simulated, opponent)) {
      return column;
    }
  }

  const priorityOrder = [3, 2, 4, 1, 5, 0, 6];
  for (const preferred of priorityOrder) {
    if (availableColumns.includes(preferred)) {
      return preferred;
    }
  }

  return availableColumns[0];
}

export function evaluateMockBotMove(board: number[][]): MockBotMoveResult {
  const column = selectMockBotColumn(board);
  if (column === null) {
    return { status: 'BLOCKED', board };
  }

  const drop = dropDisc(board, column, MOCK_BOT_PLAYER);
  if (!drop) {
    return { status: 'BLOCKED', board };
  }

  const updatedBoard = drop.nextBoard;

  if (checkWin(updatedBoard, MOCK_BOT_PLAYER)) {
    return { status: 'LOSS', board: updatedBoard };
  }

  if (isBoardFull(updatedBoard)) {
    return { status: 'DRAW', board: updatedBoard };
  }

  return { status: 'CONTINUE', board: updatedBoard };
}

export function scheduleMockBotTurn(callback: () => void): number {
  return window.setTimeout(callback, MOCK_BOT_DELAY_MS);
}

export function cancelMockBotTurn(timerId: number) {
  window.clearTimeout(timerId);
}
