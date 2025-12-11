import { useEffect, useState } from 'react';
import styles from './Leaderboard.module.css';
import { API_BASE } from '../config';

type LeaderboardEntry = {
  username: string;
  wins: number;
  losses?: number;
};

type LeaderboardProps = {
  highlightUsername?: string;
};

export function Leaderboard({ highlightUsername }: LeaderboardProps) {
  const [entries, setEntries] = useState<LeaderboardEntry[]>([]);

  useEffect(() => {
    fetch(`${API_BASE}/leaderboard`)
      .then((res) => res.json())
      .then((data) => setEntries(data))
      .catch((err) => console.error('Leaderboard fetch failed:', err));
  }, []);

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
            {entries.map((entry, index) => {
              const isCurrentUser = highlightUsername && entry.username.toLowerCase() === highlightUsername.toLowerCase();
              return (
                <tr key={entry.username} data-active={isCurrentUser ? 'true' : undefined}>
                  <td>{index + 1}</td>
                  <td>{entry.username}</td>
                  <td>{entry.wins}</td>
                  <td>{entry.losses ?? 0}</td>
                </tr>
              );
            })}
          </tbody>
        </table>
      </div>
    </section>
  );
}
