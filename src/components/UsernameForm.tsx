import { FormEvent, useState } from 'react';
import styles from './UsernameForm.module.css';

type UsernameFormProps = {
  onSubmit: (username: string) => void;
};

export function UsernameForm({ onSubmit }: UsernameFormProps) {
  const [value, setValue] = useState('');
  const [error, setError] = useState('');

  const handleSubmit = (event: FormEvent<HTMLFormElement>) => {
    event.preventDefault();
    const trimmed = value.trim();
    if (!trimmed) {
      setError('Add a call sign to join the lobby.');
      return;
    }

    setError('');
    onSubmit(trimmed);
  };

  return (
    <div className={styles.wrapper}>
      <header className={styles.header}>
        <p className={styles.tag}>Phase 1 / Frontend</p>
        <h1 className={styles.title}>Connect Four - Engineering Lobby</h1>
        <p className={styles.subtitle}>
          Pick a username so we can match you with the backend bot. Dark mode is always on to keep eyes fresh during long reviews.
        </p>
      </header>
      <form className={styles.form} onSubmit={handleSubmit} noValidate>
        <label className={styles.label} htmlFor="username">
          Username
        </label>
        <input
          id="username"
          name="username"
          type="text"
          value={value}
          onChange={(event) => setValue(event.target.value)}
          className={styles.input}
          placeholder="e.g. infra_alchemist"
          maxLength={24}
          aria-describedby={error ? 'username-error' : undefined}
        />
        {error && (
          <p className={styles.error} role="alert" id="username-error">
            {error}
          </p>
        )}
        <button type="submit" className={styles.cta}>
          Enter Lobby
        </button>
      </form>
    </div>
  );
}
