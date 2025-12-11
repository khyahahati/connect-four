import styles from './Leaderboard.module.css';

export type LeaderboardEntry = {
  rank: number;
  username: string;
  wins: number;
  losses: number;
};

type LeaderboardProps = {
  entries: LeaderboardEntry[];
  highlightUsername?: string;
};

export function Leaderboard({ entries, highlightUsername }: LeaderboardProps) {
  return (
    <section className={styles.wrapper} aria-labelledby="leaderboard-title">
      <header className={styles.header}>
        <h2 className={styles.title} id="leaderboard-title">
          Leaderboard
        </h2>
        <p className={styles.subtitle}>Backend bot stays active - see who is keeping pace.</p>
      </header>
      <div className={styles.tableShell} role="region" aria-live="polite">
        <table className={styles.table}>
          <thead>
            <tr>
              <th scope="col">Rank</th>
              <th scope="col">Username</th>
              <th scope="col">Wins</th>
              <th scope="col">Losses</th>
            </tr>
          </thead>
          <tbody>
            {entries.map((entry) => {
              const isCurrentUser = highlightUsername && entry.username.toLowerCase() === highlightUsername.toLowerCase();
              return (
                <tr key={entry.rank} data-active={isCurrentUser ? 'true' : undefined}>
                  <td>{entry.rank}</td>
                  <td>{entry.username}</td>
                  <td>{entry.wins}</td>
                  <td>{entry.losses}</td>
                </tr>
              );
            })}
          </tbody>
        </table>
      </div>
    </section>
  );
}
