import styles from './App.module.css';
import { UsernameForm } from './components/UsernameForm';
import { MatchmakingScreen } from './components/MatchmakingScreen';
import { GameBoard } from './components/GameBoard';
import { StatusBar } from './components/StatusBar';
import { Leaderboard } from './components/Leaderboard';
import { useGameController } from './state/useGameController';

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
              <Leaderboard highlightUsername={username} />
            </aside>
          </div>
        ) : null}
      </main>
    </div>
  );
}
 
