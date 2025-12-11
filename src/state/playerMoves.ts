import { checkWin, dropDisc, isBoardFull } from '../utils/board';
import type { GameResult } from './gameState';

export type PlayerMoveEvaluation =
  | { type: 'COLUMN_FULL' }
  | { type: 'CONTINUE'; board: number[][] }
  | { type: GameResult; board: number[][] };

export function evaluatePlayerMove(board: number[][], columnIndex: number, player: 1 | 2): PlayerMoveEvaluation {
  const drop = dropDisc(board, columnIndex, player);
  if (!drop) {
    return { type: 'COLUMN_FULL' };
  }

  const updatedBoard = drop.nextBoard;

  if (checkWin(updatedBoard, player)) {
    return { type: 'WIN', board: updatedBoard };
  }

  if (isBoardFull(updatedBoard)) {
    return { type: 'DRAW', board: updatedBoard };
  }

  return { type: 'CONTINUE', board: updatedBoard };
}
