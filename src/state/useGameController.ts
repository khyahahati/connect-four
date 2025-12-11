import { useCallback, useEffect, useReducer } from 'react';
import {
  gameActions,
  gameReducer,
  initialGameState,
  type ServerBoardUpdatePayload,
  type ServerGameOverPayload,
  type ServerGameStartPayload
} from './gameState';
import { evaluatePlayerMove } from './playerMoves';
import { useWebSocketClient } from '../network/useWebSocketClient';
import type { ServerMessage } from '../types/network';
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

  const { connect, disconnect, send, onMessage, connected, socketError } = useWebSocketClient();

  const handleServerGameStart = useCallback((payload: ServerGameStartPayload) => {
    const nextTurn = payload.firstTurn ?? payload.you;
    dispatch({
      type: 'SERVER_GAME_START',
      payload: {
        ...payload,
        firstTurn: nextTurn,
        message: payload.message ?? turnMessage(nextTurn, payload.you, payload.opponent)
      }
    });
  }, [dispatch]);

  const handleServerBoardUpdate = useCallback((payload: ServerBoardUpdatePayload) => {
    dispatch({
      type: 'SERVER_BOARD_UPDATE',
      payload: {
        ...payload,
        message: payload.message ?? turnMessage(payload.currentTurn, state.you, state.opponent)
      }
    });
  }, [dispatch, state.opponent, state.you]);

  const handleServerGameOver = useCallback((payload: ServerGameOverPayload) => {
    const fallbackMessage =
      payload.result === 'WIN'
        ? playerWinMessage
        : payload.result === 'LOSS'
          ? `${state.opponent ?? 'Opponent'} connected four.`
          : drawMessage;

    dispatch({
      type: 'SERVER_GAME_OVER',
      payload: {
        ...payload,
        message: payload.message ?? fallbackMessage
      }
    });
  }, [dispatch, state.opponent]);

  const sendMakeMoveToServer = useCallback((columnIndex: number) => {
    if (!state.gameId || !state.username || !connected) {
      return;
    }

    send({
      type: 'MAKE_MOVE',
      col: columnIndex,
      gameId: state.gameId
    });
  }, [connected, send, state.gameId, state.username]);

  useEffect(() => {
    if (state.gameMode !== 'LOCAL' || state.screen !== 'MATCHMAKING') {
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
  }, [dispatch, state.gameMode, state.screen]);

  useEffect(() => {
    if (state.gameMode !== 'LOCAL') {
      return;
    }

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
  }, [dispatch, state.board, state.currentTurn, state.gameMode, state.opponent, state.result, state.screen, state.you]);

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

      if (state.gameMode === 'LOCAL') {
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
        return;
      }

      sendMakeMoveToServer(columnIndex);
    },
    [
      dispatch,
      sendMakeMoveToServer,
      state.board,
      state.currentTurn,
      state.gameMode,
      state.opponent,
      state.result,
      state.screen,
      state.you
    ]
  );

  const handleRestart = useCallback(() => {
    if (state.gameMode === 'LOCAL') {
      dispatch(gameActions.resetGame());
      return;
    }

    if (!state.username) {
      return;
    }

    send({
      type: 'RECONNECT',
      username: state.username,
      gameId: state.gameId
    });
  }, [dispatch, send, state.gameId, state.gameMode, state.username]);

  useEffect(() => {
    const shouldConnect = state.gameMode === 'ONLINE' && Boolean(state.username) && state.screen !== 'ENTER_NAME';

    if (shouldConnect) {
      connect(state.username);
    } else {
      disconnect();
    }
  }, [connect, disconnect, state.gameMode, state.screen, state.username]);

  useEffect(() => {
    if (state.gameMode !== 'ONLINE') {
      return;
    }

    if (!connected || !state.username) {
      return;
    }

    send({
      type: 'RECONNECT',
      username: state.username,
      gameId: state.gameId
    });
  }, [connected, send, state.gameId, state.gameMode, state.username]);

  useEffect(() => {
    if (state.gameMode !== 'ONLINE') {
      return;
    }

    if (!socketError) {
      return;
    }

    dispatch(gameActions.setMessage(socketError));
  }, [dispatch, socketError, state.gameMode]);

  useEffect(() => {
    if (state.gameMode !== 'ONLINE') {
      return;
    }

    const unsubscribe = onMessage((message: ServerMessage) => {
      switch (message.type) {
        case 'GAME_START':
          handleServerGameStart({
            gameId: message.gameId,
            you: message.you,
            opponent: message.opponent
          });
          break;
        case 'BOARD_UPDATE':
          handleServerBoardUpdate({
            board: message.board,
            currentTurn: message.currentTurn
          });
          break;
        case 'GAME_OVER':
          handleServerGameOver({
            result: message.result,
            board: message.board
          });
          break;
        case 'INFO':
          dispatch(gameActions.setMessage(message.message));
          break;
        default:
          break;
      }
    });

    return unsubscribe;
  }, [dispatch, handleServerBoardUpdate, handleServerGameOver, handleServerGameStart, onMessage, state.gameMode]);

  return {
    state,
    actions: {
      submitUsername,
      handleColumnClick,
      handleRestart
    }
  };
}
