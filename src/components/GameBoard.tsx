import styles from './GameBoard.module.css';

type GameBoardProps = {
  board: number[][];
  currentTurn: 1 | 2;
  you: 1 | 2;
  onColumnClick?: (columnIndex: number) => void;
};

export function GameBoard({ board, currentTurn, you, onColumnClick }: GameBoardProps) {
  const selfLabel = 'You';
  const opponentLabel = 'Opponent';
  const ariaLabel = currentTurn === you ? 'Your move' : `${opponentLabel} move`;
  const isInteractive = typeof onColumnClick === 'function' && currentTurn === you;

  return (
    <div className={styles.wrapper}>
      <div className={styles.legend}>
        <span className={styles.legendDot} data-player="self" aria-hidden="true" />
        {selfLabel}
        <span className={styles.legendSeparator} aria-hidden="true">/</span>
        <span className={styles.legendDot} data-player="opponent" aria-hidden="true" />
        {opponentLabel}
      </div>
      <div
        className={styles.board}
        role="grid"
        aria-live="polite"
        aria-label={`Connect Four board, ${ariaLabel}`}
        data-interactive={isInteractive ? 'true' : 'false'}
      >
        {board.map((row, rowIndex) => (
          <div className={styles.row} role="row" key={`row-${rowIndex}`}>
            {row.map((cell, columnIndex) => (
              <button
                key={`cell-${rowIndex}-${columnIndex}`}
                type="button"
                className={createCellClass(cell, isInteractive)}
                onClick={() => onColumnClick?.(columnIndex)}
                disabled={!isInteractive}
                aria-label={`Drop disc in column ${columnIndex + 1}`}
              >
                <span className={createDiscClass(cell)} />
                <span className={styles.cellIndex} aria-hidden="true">
                  {columnIndex + 1}
                </span>
              </button>
            ))}
          </div>
        ))}
      </div>
      <p className={styles.helper}>Click a column to drop a disc. The grid locks whenever it is not your turn.</p>
    </div>
  );
}

function createCellClass(cell: number, isInteractive: boolean) {
  const classes = [styles.cell];
  if (cell === 1) {
    classes.push(styles.playerOne);
  }
  if (cell === 2) {
    classes.push(styles.playerTwo);
  }
  if (isInteractive) {
    classes.push(styles.actionable);
  }
  return classes.join(' ');
}

function createDiscClass(cell: number) {
  const classes = [styles.disc];
  if (cell === 1) {
    classes.push(styles.discPlayerOne);
  }
  if (cell === 2) {
    classes.push(styles.discPlayerTwo);
  }
  return classes.join(' ');
}
