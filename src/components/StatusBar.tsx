import styles from './StatusBar.module.css';
import type { GameResult, Screen } from '../state/gameState';

type StatusBarProps = {
  screen: Screen;
  username: string;
  opponent?: string;
  currentTurn: 1 | 2;
  you: 1 | 2;
  message?: string;
  result?: GameResult;
  onRestart: () => void;
};

const resultCopy: Record<GameResult, string> = {
  WIN: 'You won the round.',
  LOSS: 'Opponent took the round.',
  DRAW: 'Round ended in a draw.'
};

export function StatusBar({ screen, username, opponent, currentTurn, you, message, result, onRestart }: StatusBarProps) {
  const engineerLabel = username || 'Pending';
  const turnLabel = currentTurn === you ? 'You' : opponent ?? 'Opponent';
  const showRestart = screen === 'GAME_OVER';
  const hintCopy = screen === 'IN_GAME' ? 'Tip: track diagonal threats early to prevent surprise wins.' : undefined;
  const resultMessage = result ? resultCopy[result] : undefined;

  return (
    <section className={styles.statusBar} aria-live="polite">
      <div className={styles.infoBlock}>
        <h2 className={styles.heading}>Match Dashboard</h2>
        <div className={styles.metaRow}>
          <span className={styles.metaLabel}>Engineer</span>
          <span className={styles.metaValue}>{engineerLabel}</span>
        </div>
        <div className={styles.metaRow}>
          <span className={styles.metaLabel}>Turn</span>
          <span className={styles.metaValue} data-active={currentTurn === you ? 'self' : 'opponent'}>
            {turnLabel}
          </span>
        </div>
        <div className={styles.metaRow}>
          <span className={styles.metaLabel}>Opponent</span>
          <span className={styles.metaValue}>{opponent ?? 'TBD'}</span>
        </div>
      </div>

      <div className={styles.statusMessage}>
        <p className={styles.messageText}>{message ?? 'Waiting for next action.'}</p>
        {resultMessage && <p className={styles.result}>{resultMessage}</p>}
        {hintCopy && <p className={styles.hint}>{hintCopy}</p>}
        <div className={styles.controls}>
          <button type="button" className={styles.restart} onClick={onRestart} disabled={!showRestart}>
            Restart match
          </button>
        </div>
      </div>
    </section>
  );
}
