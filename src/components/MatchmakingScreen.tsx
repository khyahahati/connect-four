import styles from './MatchmakingScreen.module.css';
import type { GameMode } from '../state/gameState';

type MatchmakingScreenProps = {
  username: string;
  mode: GameMode;
};

export function MatchmakingScreen({ username, mode }: MatchmakingScreenProps) {
  const localCopy = username ? `${username}, we're pairing you with our backend bot.` : 'Looking for an opponent.';
  const onlineCopy = username
    ? `${username}, waiting for the server to confirm your match.`
    : 'Connecting to the game service.';

  const bodyCopy = mode === 'LOCAL' ? localCopy : onlineCopy;

  return (
    <div className={styles.wrapper}>
      <header className={styles.header}>
        <h2 className={styles.title}>Matchmaking in progress</h2>
        <p className={styles.body}>{bodyCopy}</p>
      </header>
      <div className={styles.panel}>
        <div className={styles.spinner} aria-hidden="true" />
        <div className={styles.statusCopy}>
          <p className={styles.statusHeadline}>Analyzing open matches...</p>
          <p className={styles.statusSubline}>Expect a short wait while we sync the board state.</p>
        </div>
      </div>
    </div>
  );
}
