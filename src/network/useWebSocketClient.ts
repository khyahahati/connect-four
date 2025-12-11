import type { ClientMessage } from '../types/network';

type MessageCallback = (message: unknown) => void;

export function useWebSocketClient() {
  let onMessageCallback: MessageCallback | null = null;

  const connect = (_username: string) => {
    // TODO: wire real WebSocket connection in online mode.
  };

  const disconnect = () => {
    // TODO: close active WebSocket connection when implemented.
  };

  const send = (_message: ClientMessage) => {
    // TODO: forward client message to the backend WebSocket once available.
  };

  const onMessage = (callback: MessageCallback) => {
    onMessageCallback = callback;
    return () => {
      onMessageCallback = null;
    };
  };

  return {
    connect,
    disconnect,
    send,
    onMessage
  };
}
