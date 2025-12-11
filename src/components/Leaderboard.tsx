import { useEffect, useState } from 'react';
import styles from './Leaderboard.module.css';

type LeaderboardEntry = {
  username: string;
  wins: number;
  losses: number;
};

type LeaderboardProps = {
  highlightUsername?: string;
};

export function Leaderboard({ highlightUsername }: LeaderboardProps) {
  const [entries, setEntries] = useState<LeaderboardEntry[]>([]);

  useEffect(() => {
    let isMounted = true;

    fetch('http://localhost:8080/leaderboard')
      .then((res) => res.json())
      .then((data: Array<{ username: string; wins: number }>) => {
        if (!isMounted) return;
        const withComputed = data.map((entry) => ({
          username: entry.username,
          wins: entry.wins,
          losses: 0
        }));
        setEntries(withComputed);
      })
      .catch(() => {
        if (!isMounted) return;
        setEntries([]);
      });

    return () => {
      isMounted = false;
    };
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
