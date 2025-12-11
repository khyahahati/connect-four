import styles from './App.module.css';
import { UsernameForm } from './components/UsernameForm';
import { MatchmakingScreen } from './components/MatchmakingScreen';
import { GameBoard } from './components/GameBoard';
import { StatusBar } from './components/StatusBar';
import { Leaderboard, LeaderboardEntry } from './components/Leaderboard';
import { useGameController } from './state/useGameController';

const leaderboardSeed: LeaderboardEntry[] = [
  { rank: 1, username: 'backend_bot', wins: 18, losses: 3 },
  { rank: 2, username: 'ops_guru', wins: 14, losses: 6 },
  { rank: 3, username: 'api_architect', wins: 11, losses: 8 },
  { rank: 4, username: 'queue_master', wins: 9, losses: 9 },
  { rank: 5, username: 'cache_hit', wins: 7, losses: 10 }
];

export default function App() {
  const {
    state,
    actions: { submitUsername, handleColumnClick, handleRestart }
  } = useGameController();

  const { screen, username, opponent, board, currentTurn, you, result, message, gameMode } = state;
  const isBoardInteractive = screen === 'IN_GAME' && currentTurn === you && !result;

  return (
    <div className={styles.appShell}>
      <main className={styles.mainPanel}>
        {screen === 'ENTER_NAME' && (
          <section className={styles.sectionCard}>
            <UsernameForm onSubmit={submitUsername} mode={gameMode} />
          </section>
        )}

        {screen === 'MATCHMAKING' && (
          <section className={styles.sectionCard}>
            <MatchmakingScreen username={username} mode={gameMode} />
          </section>
        )}

        {screen === 'IN_GAME' || screen === 'GAME_OVER' ? (
          <div className={styles.gameLayout}>
            <section className={styles.gamePanel}>
              <StatusBar
                screen={screen}
                username={username}
                opponent={opponent}
                currentTurn={currentTurn}
                you={you}
                message={message}
                result={result}
                onRestart={handleRestart}
              />
              <GameBoard
                board={board}
                currentTurn={currentTurn}
                you={you}
                onColumnClick={isBoardInteractive ? handleColumnClick : undefined}
              />
            </section>
            <aside className={styles.leaderboardPanel}>
              <Leaderboard entries={leaderboardSeed} highlightUsername={username} />
            </aside>
          </div>
        ) : null}
      </main>
    </div>
  );
}
 
