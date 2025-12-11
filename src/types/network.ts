export type ClientMessage =
  | { type: 'MAKE_MOVE'; col: number; gameId: string }
  | { type: 'RECONNECT'; username: string; gameId?: string };

export type ServerMessage =
  | { type: 'GAME_START'; gameId: string; you: 1 | 2; opponent: string }
  | { type: 'BOARD_UPDATE'; board: number[][]; currentTurn: 1 | 2 }
  | { type: 'GAME_OVER'; result: 'WIN' | 'LOSS' | 'DRAW'; board: number[][] }
  | { type: 'INFO'; message: string };
