import { useCallback, useEffect, useRef, useState } from 'react';
import type { ClientMessage, ServerMessage } from '../types/network';

type MessageCallback = (message: ServerMessage) => void;

const RECONNECT_DELAY_MS = 1500;

export function useWebSocketClient() {
  const socketRef = useRef<WebSocket | null>(null);
  const reconnectTimerRef = useRef<number | null>(null);
  const lastUsernameRef = useRef<string | null>(null);
  const shouldReconnectRef = useRef(false);
  const messageHandlersRef = useRef(new Set<MessageCallback>());

  const [connected, setConnected] = useState(false);
  const [socketError, setSocketError] = useState<string | null>(null);

  const clearReconnectTimer = useCallback(() => {
    if (reconnectTimerRef.current !== null) {
      window.clearTimeout(reconnectTimerRef.current);
      reconnectTimerRef.current = null;
    }
  }, []);

  const disconnect = useCallback(() => {
    shouldReconnectRef.current = false;
    clearReconnectTimer();

    const socket = socketRef.current;
    if (socket) {
      socketRef.current = null;
      socket.close();
    }

    setConnected(false);
  }, [clearReconnectTimer]);

  const connect = useCallback(
    (username: string) => {
      const baseUrl = import.meta.env.VITE_API_WS_URL;
      if (!baseUrl) {
        setSocketError('Missing VITE_API_WS_URL environment variable.');
        return;
      }

      lastUsernameRef.current = username;
      shouldReconnectRef.current = true;

      const existing = socketRef.current;
      if (existing && (existing.readyState === WebSocket.OPEN || existing.readyState === WebSocket.CONNECTING)) {
        return;
      }

      clearReconnectTimer();

      const normalized = baseUrl.endsWith('/') ? baseUrl.slice(0, -1) : baseUrl;

      try {
        const socket = new WebSocket(`${normalized}/ws?username=${encodeURIComponent(username)}`);
        socketRef.current = socket;

        socket.onopen = () => {
          setSocketError(null);
          setConnected(true);
          clearReconnectTimer();
        };

        socket.onmessage = (event: MessageEvent<string>) => {
          if (typeof event.data !== 'string') {
            return;
          }

          try {
            const parsed = JSON.parse(event.data) as ServerMessage;
            messageHandlersRef.current.forEach((handler) => handler(parsed));
          } catch (error) {
            const message = error instanceof Error ? error.message : 'Unknown parse error';
            setSocketError(`Failed to parse server message: ${message}`);
          }
        };

        socket.onclose = () => {
          setConnected(false);
          socketRef.current = null;

          if (!shouldReconnectRef.current) {
            return;
          }

          const usernameSnapshot = lastUsernameRef.current;
          if (!usernameSnapshot) {
            return;
          }

          if (reconnectTimerRef.current !== null) {
            return;
          }

          reconnectTimerRef.current = window.setTimeout(() => {
            reconnectTimerRef.current = null;
            connect(usernameSnapshot);
          }, RECONNECT_DELAY_MS);
        };

        socket.onerror = (event) => {
          const description = event instanceof ErrorEvent ? event.message : 'WebSocket error';
          setSocketError(description);
        };
      } catch (error) {
        const message = error instanceof Error ? error.message : 'Unknown connection error';
        setSocketError(message);
        setConnected(false);

        if (!shouldReconnectRef.current) {
          return;
        }

        if (reconnectTimerRef.current === null) {
          reconnectTimerRef.current = window.setTimeout(() => {
            reconnectTimerRef.current = null;
            const usernameSnapshot = lastUsernameRef.current;
            if (usernameSnapshot) {
              connect(usernameSnapshot);
            }
          }, RECONNECT_DELAY_MS);
        }
      }
    },
    [clearReconnectTimer]
  );

  const send = useCallback((message: ClientMessage) => {
    const socket = socketRef.current;
    if (!socket || socket.readyState !== WebSocket.OPEN) {
      return;
    }

    socket.send(JSON.stringify(message));
  }, []);

  const onMessage = useCallback((callback: MessageCallback) => {
    messageHandlersRef.current.add(callback);
    return () => {
      messageHandlersRef.current.delete(callback);
    };
  }, []);

  useEffect(() => {
    return () => {
      disconnect();
      messageHandlersRef.current.clear();
    };
  }, [disconnect]);

  return {
    connect,
    disconnect,
    send,
    onMessage,
    connected,
    socketError
  };
}
