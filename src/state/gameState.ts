import { cloneBoard } from '../utils/board';

export type Screen = 'ENTER_NAME' | 'MATCHMAKING' | 'IN_GAME' | 'GAME_OVER';

export type GameResult = 'WIN' | 'LOSS' | 'DRAW';

export type GameMode = 'LOCAL' | 'ONLINE';

export interface GameState {
  screen: Screen;
  username: string;
  opponent?: string;
  board: number[][];
  currentTurn: 1 | 2;
  you: 1 | 2;
  gameId?: string;
  result?: GameResult;
  message?: string;
  gameMode: GameMode;
}

export interface StartGamePayload {
  opponent?: string;
  you: 1 | 2;
  firstTurn?: 1 | 2;
  board?: number[][];
  gameId?: string;
  message?: string;
}

export interface UpdateBoardPayload {
  board: number[][];
  currentTurn: 1 | 2;
  message?: string;
}

export interface EndGamePayload {
  result: GameResult;
  message?: string;
  board?: number[][];
}

export interface ServerGameStartPayload {
  gameId: string;
  you: 1 | 2;
  opponent: string;
  firstTurn?: 1 | 2;
  message?: string;
  board?: number[][];
}

export interface ServerBoardUpdatePayload {
  board: number[][];
  currentTurn: 1 | 2;
  message?: string;
}

export interface ServerGameOverPayload {
  result: GameResult;
  board: number[][];
  message?: string;
}

export type GameAction =
  | { type: 'SET_USERNAME'; username: string }
  | { type: 'CHANGE_SCREEN'; screen: Screen; message?: string }
  | { type: 'START_GAME'; payload: StartGamePayload }
  | { type: 'UPDATE_BOARD'; payload: UpdateBoardPayload }
  | { type: 'END_GAME'; payload: EndGamePayload }
  | { type: 'RESET_GAME' }
  | { type: 'SET_MESSAGE'; message: string }
  | { type: 'SERVER_GAME_START'; payload: ServerGameStartPayload }
  | { type: 'SERVER_BOARD_UPDATE'; payload: ServerBoardUpdatePayload }
  | { type: 'SERVER_GAME_OVER'; payload: ServerGameOverPayload };

const ROWS = 6;
const COLUMNS = 7;

export const BOARD_ROWS = ROWS;
export const BOARD_COLUMNS = COLUMNS;

export function createEmptyBoard(rows: number = ROWS, columns: number = COLUMNS): number[][] {
  return Array.from({ length: rows }, () => Array.from({ length: columns }, () => 0));
}

export const initialGameState: GameState = {
  screen: 'ENTER_NAME',
  username: '',
  opponent: undefined,
  board: createEmptyBoard(),
  currentTurn: 1,
  you: 1,
  gameId: undefined,
  result: undefined,
  message: 'Enter a username to start a match.',
  gameMode: 'ONLINE'
};

export function setUsername(state: GameState, username: string): GameState {
  return {
    ...state,
    username
  };
}

export function changeScreen(state: GameState, screen: Screen, message?: string): GameState {
  return {
    ...state,
    screen,
    message: message ?? state.message
  };
}

export function startGame(state: GameState, payload: StartGamePayload): GameState {
  const nextBoard = payload.board ? cloneBoard(payload.board) : createEmptyBoard();
  const nextTurn = payload.firstTurn ?? payload.you;
  return {
    ...state,
    screen: 'IN_GAME',
    opponent: payload.opponent ?? state.opponent,
    board: nextBoard,
    currentTurn: nextTurn,
    you: payload.you,
    gameId: payload.gameId,
    result: undefined,
    message: payload.message ?? state.message
  };
}

export function updateBoard(state: GameState, payload: UpdateBoardPayload): GameState {
  return {
    ...state,
    board: cloneBoard(payload.board),
    currentTurn: payload.currentTurn,
    message: payload.message ?? state.message
  };
}

export function endGame(state: GameState, payload: EndGamePayload): GameState {
  return {
    ...state,
    screen: 'GAME_OVER',
    result: payload.result,
    board: payload.board ? cloneBoard(payload.board) : state.board,
    message: payload.message ?? state.message
  };
}

export function resetGame(state: GameState): GameState {
  return {
    ...state,
    screen: 'IN_GAME',
    board: createEmptyBoard(),
    currentTurn: state.you,
    result: undefined,
    message: 'Your move - drop a disc to get four in a row.'
  };
}

export function setMessage(state: GameState, message: string): GameState {
  return {
    ...state,
    message
  };
}

export function gameReducer(state: GameState, action: GameAction): GameState {
  switch (action.type) {
    case 'SET_USERNAME':
      return setUsername(state, action.username);
    case 'CHANGE_SCREEN':
      return changeScreen(state, action.screen, action.message);
    case 'START_GAME':
      return startGame(state, action.payload);
    case 'UPDATE_BOARD':
      return updateBoard(state, action.payload);
    case 'END_GAME':
      return endGame(state, action.payload);
    case 'RESET_GAME':
      return resetGame(state);
    case 'SET_MESSAGE':
      return setMessage(state, action.message);
    case 'SERVER_GAME_START': {
      return startGame(state, {
        opponent: action.payload.opponent,
        you: action.payload.you,
        firstTurn: action.payload.firstTurn,
        board: action.payload.board,
        gameId: action.payload.gameId,
        message: action.payload.message
      });
    }
    case 'SERVER_BOARD_UPDATE':
      return updateBoard(state, {
        board: action.payload.board,
        currentTurn: action.payload.currentTurn,
        message: action.payload.message
      });
    case 'SERVER_GAME_OVER':
      return endGame(state, {
        result: action.payload.result,
        board: action.payload.board,
        message: action.payload.message
      });
    default:
      return state;
  }
}

export const gameActions = {
  setUsername: (username: string): GameAction => ({ type: 'SET_USERNAME', username }),
  changeScreen: (screen: Screen, message?: string): GameAction => ({ type: 'CHANGE_SCREEN', screen, message }),
  startGame: (payload: StartGamePayload): GameAction => ({ type: 'START_GAME', payload }),
  updateBoard: (payload: UpdateBoardPayload): GameAction => ({ type: 'UPDATE_BOARD', payload }),
  endGame: (payload: EndGamePayload): GameAction => ({ type: 'END_GAME', payload }),
  resetGame: (): GameAction => ({ type: 'RESET_GAME' }),
  setMessage: (message: string): GameAction => ({ type: 'SET_MESSAGE', message })
};
