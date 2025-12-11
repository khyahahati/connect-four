import { useCallback, useEffect, useReducer } from 'react';
import { gameActions, gameReducer, initialGameState } from './gameState';
import { evaluatePlayerMove } from './playerMoves';
import {
  MOCK_BOT_ENABLED,
  MOCK_BOT_NAME,
  MOCK_BOT_PLAYER,
  cancelMockBotTurn,
  evaluateMockBotMove,
  scheduleMockBotTurn
} from '../utils/mockBot';

const MATCHMAKING_DELAY_MS = 1600;

const columnFullMessage = 'Column full - pick another lane.';
const drawMessage = 'Dead heat - the grid is full.';
const botBlockedMessage = 'Bot cannot find a valid column. Restart to continue.';
const playerWinMessage = 'You connected four! Well played.';

const turnMessage = (turn: 1 | 2, you: 1 | 2, opponent?: string): string =>
  turn === you ? 'Your move - drop a disc to get four in a row.' : `${opponent ?? 'Opponent'} thinking - stay sharp.`;

export function useGameController() {
  const [state, dispatch] = useReducer(gameReducer, initialGameState);

  useEffect(() => {
    if (state.screen !== 'MATCHMAKING') {
      return;
    }

    const opponentName = MOCK_BOT_ENABLED ? MOCK_BOT_NAME : 'Queued Opponent';
    const timerId = window.setTimeout(() => {
      dispatch(
        gameActions.startGame({
          opponent: opponentName,
          you: 1,
          firstTurn: 1,
          message: turnMessage(1, 1, opponentName)
        })
      );
    }, MATCHMAKING_DELAY_MS);

    return () => window.clearTimeout(timerId);
  }, [dispatch, state.screen]);

  useEffect(() => {
    if (!MOCK_BOT_ENABLED || state.opponent !== MOCK_BOT_NAME) {
      return;
    }

    if (state.screen !== 'IN_GAME' || state.currentTurn === state.you || Boolean(state.result)) {
      return;
    }

    const timerId = scheduleMockBotTurn(() => {
      const evaluation = evaluateMockBotMove(state.board);

      if (evaluation.status === 'BLOCKED') {
        dispatch(gameActions.endGame({ result: 'DRAW', message: botBlockedMessage, board: state.board }));
        return;
      }

      if (evaluation.status === 'LOSS') {
        const lossMessage = `${state.opponent ?? 'Opponent'} connected four.`;
        dispatch(gameActions.updateBoard({ board: evaluation.board, currentTurn: MOCK_BOT_PLAYER, message: lossMessage }));
        dispatch(gameActions.endGame({ result: 'LOSS', message: lossMessage, board: evaluation.board }));
        return;
      }

      if (evaluation.status === 'DRAW') {
        dispatch(gameActions.updateBoard({ board: evaluation.board, currentTurn: state.you, message: drawMessage }));
        dispatch(gameActions.endGame({ result: 'DRAW', message: drawMessage, board: evaluation.board }));
        return;
      }

      dispatch(
        gameActions.updateBoard({
          board: evaluation.board,
          currentTurn: state.you,
          message: turnMessage(state.you, state.you, state.opponent)
        })
      );
    });

    return () => cancelMockBotTurn(timerId);
  }, [dispatch, state.board, state.currentTurn, state.opponent, state.result, state.screen, state.you]);

  const submitUsername = useCallback(
    (value: string) => {
      const trimmed = value.trim();
      if (!trimmed) {
        return;
      }
      dispatch(gameActions.setUsername(trimmed));
      dispatch(gameActions.changeScreen('MATCHMAKING', 'Pairing you with an opponent...'));
    },
    [dispatch]
  );

  const handleColumnClick = useCallback(
    (columnIndex: number) => {
      if (state.screen !== 'IN_GAME' || state.currentTurn !== state.you || Boolean(state.result)) {
        return;
      }

      const evaluation = evaluatePlayerMove(state.board, columnIndex, state.you);

      if (evaluation.type === 'COLUMN_FULL') {
        dispatch(gameActions.setMessage(columnFullMessage));
        return;
      }

      if (evaluation.type === 'WIN') {
        dispatch(gameActions.updateBoard({ board: evaluation.board, currentTurn: state.you, message: playerWinMessage }));
        dispatch(gameActions.endGame({ result: 'WIN', message: playerWinMessage, board: evaluation.board }));
        return;
      }

      if (evaluation.type === 'DRAW') {
        dispatch(gameActions.updateBoard({ board: evaluation.board, currentTurn: state.you, message: drawMessage }));
        dispatch(gameActions.endGame({ result: 'DRAW', message: drawMessage, board: evaluation.board }));
        return;
      }

      const nextTurn: 1 | 2 = state.you === 1 ? 2 : 1;
      dispatch(
        gameActions.updateBoard({
          board: evaluation.board,
          currentTurn: nextTurn,
          message: turnMessage(nextTurn, state.you, state.opponent)
        })
      );
    },
    [dispatch, state.board, state.currentTurn, state.opponent, state.result, state.screen, state.you]
  );

  const handleRestart = useCallback(() => {
    dispatch(gameActions.resetGame());
  }, [dispatch]);

  return {
    state,
    actions: {
      submitUsername,
      handleColumnClick,
      handleRestart
    }
  };
}
