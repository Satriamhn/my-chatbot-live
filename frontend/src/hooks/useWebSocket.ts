import { useEffect, useRef, useCallback, useState } from "react";

export interface WSMessage {
  type: "chat_message" | "token" | "error" | "pong" | "system";
  payload: unknown;
  session_id?: string;
}

export interface StreamingToken {
  content: string;
  is_last: boolean;
}

interface UseWebSocketOptions {
  sessionId: string | null;
  orgId: string | null;
  token: string | null;
  onToken?: (chunk: StreamingToken) => void;
  onError?: (msg: string) => void;
}

const WS_BASE = import.meta.env.VITE_WS_URL || "ws://localhost:8080";
const RECONNECT_DELAYS = [1000, 2000, 4000, 8000];

export function useWebSocket({
  sessionId,
  orgId,
  token,
  onToken,
  onError,
}: UseWebSocketOptions) {
  const wsRef = useRef<WebSocket | null>(null);
  const reconnectAttempt = useRef(0);
  const reconnectTimer = useRef<ReturnType<typeof setTimeout> | null>(null);
  const [isConnected, setIsConnected] = useState(false);

  // Use refs for callbacks — prevents reconnect storm on every parent re-render
  const onTokenRef = useRef(onToken);
  const onErrorRef = useRef(onError);
  useEffect(() => { onTokenRef.current = onToken; }, [onToken]);
  useEffect(() => { onErrorRef.current = onError; }, [onError]);

  const connect = useCallback(() => {
    if (!sessionId || !orgId || !token) return;

    // Cleanup any existing connection
    if (wsRef.current) {
      wsRef.current.onclose = null; // prevent auto-reconnect from old close
      wsRef.current.close();
    }

    const url = `${WS_BASE}/api/v1/sessions/ws?token=${encodeURIComponent(token)}&org_id=${encodeURIComponent(orgId)}`;
    const ws = new WebSocket(url);
    wsRef.current = ws;

    ws.onopen = () => {
      setIsConnected(true);
      reconnectAttempt.current = 0;
    };

    ws.onmessage = (event) => {
      try {
        const msg: WSMessage = JSON.parse(event.data);
        if (msg.type === "token") {
          onTokenRef.current?.(msg.payload as StreamingToken);
        } else if (msg.type === "error") {
          const errPayload = msg.payload as { message: string };
          onErrorRef.current?.(errPayload?.message || "Unknown error");
        }
        // "system" and "pong" types are silently ignored
      } catch {
        // ignore parse errors
      }
    };

    ws.onclose = () => {
      setIsConnected(false);
      const delay = RECONNECT_DELAYS[Math.min(reconnectAttempt.current, RECONNECT_DELAYS.length - 1)];
      reconnectAttempt.current++;
      reconnectTimer.current = setTimeout(connect, delay);
    };

    ws.onerror = () => {
      ws.close();
    };
  }, [sessionId, orgId, token]); // callbacks removed from deps — using refs

  useEffect(() => {
    connect();
    return () => {
      if (reconnectTimer.current) clearTimeout(reconnectTimer.current);
      if (wsRef.current) {
        wsRef.current.onclose = null; // prevent reconnect on intentional unmount
        wsRef.current.close();
      }
    };
  }, [connect]);

  const sendMessage = useCallback(
    (content: string) => {
      if (!wsRef.current || wsRef.current.readyState !== WebSocket.OPEN) return false;
      wsRef.current.send(
        JSON.stringify({
          type: "chat_message",
          session_id: sessionId,
          payload: { content },
        })
      );
      return true;
    },
    [sessionId]
  );

  return { isConnected, sendMessage };
}
