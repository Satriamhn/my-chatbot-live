import api from "./api";

// ─── Types ──────────────────────────────────────────────────────────────────

export interface Message {
  id: string;
  sender: "user" | "bot" | "human_agent" | "system";
  text: string;
  timestamp?: string;
}

export interface Conversation {
  id: string;
  user: string;
  avatar?: string;
  status?: "online" | "offline" | "busy" | "away";
  lastMessage: string;
  unread: number;
  sessionStatus: "bot_handling" | "human_assigned" | "closed";
  messages: Message[];
}

// ─── Backend response mappers ────────────────────────────────────────────────

interface BackendSession {
  id: string;
  title: string;
  status: string;
  user_id: string;
  created_at: string;
  messages?: BackendMessage[];
}

interface BackendMessage {
  id: string;
  sender: string;
  content: string;
  created_at: string;
}

function mapSession(s: BackendSession, messages: BackendMessage[] = []): Conversation {
  const msgs = messages.map(mapMessage);
  const last = msgs[msgs.length - 1];
  return {
    id: s.id,
    user: s.title || "Visitor",
    sessionStatus: s.status as Conversation["sessionStatus"],
    lastMessage: last?.text || "No messages yet",
    unread: 0,
    messages: msgs,
  };
}

function mapMessage(m: BackendMessage): Message {
  return {
    id: m.id,
    sender: m.sender as Message["sender"],
    text: m.content,
    timestamp: m.created_at,
  };
}

// ─── Service ─────────────────────────────────────────────────────────────────

export const inboxService = {
  getConversations: async (): Promise<Conversation[]> => {
    const res = await api.get("/api/v1/sessions");
    const sessions: BackendSession[] = res.data;
    return sessions.map((s) => mapSession(s));
  },

  getMessages: async (sessionId: string): Promise<Message[]> => {
    const res = await api.get(`/api/v1/sessions/${sessionId}/messages`);
    const msgs: BackendMessage[] = res.data;
    return msgs.map(mapMessage);
  },

  createSession: async (title: string): Promise<Conversation> => {
    const res = await api.post("/api/v1/sessions", { title });
    return mapSession(res.data);
  },

  sendMessage: async (sessionId: string, content: string): Promise<Message> => {
    const res = await api.post(`/api/v1/sessions/${sessionId}/messages`, {
      content,
      sender_type: "human_agent",
    });
    return mapMessage(res.data);
  },

  takeoverSession: async (sessionId: string): Promise<void> => {
    await api.post(`/api/v1/sessions/${sessionId}/takeover`);
  },

  returnToBotMode: async (sessionId: string): Promise<void> => {
    await api.post(`/api/v1/sessions/${sessionId}/return-to-bot`);
  },
};
