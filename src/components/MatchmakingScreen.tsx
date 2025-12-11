import styles from './MatchmakingScreen.module.css';

type MatchmakingScreenProps = {
  username: string;
};

export function MatchmakingScreen({ username }: MatchmakingScreenProps) {
  return (
    <div className={styles.wrapper}>
      <header className={styles.header}>
        <h2 className={styles.title}>Matchmaking in progress</h2>
        <p className={styles.body}>
          {username ? `${username}, we're pairing you with our backend bot.` : 'Looking for an opponent.'}
        </p>
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
