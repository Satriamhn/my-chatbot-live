import { useState, useEffect, useCallback, useRef } from "react";
import PageBreadcrumb from "../../components/common/PageBreadCrumb";
// import ComponentCard from "../../components/common/ComponentCard";
import Button from "../../components/ui/button/Button";
import { Send, Wifi, WifiOff, Plus, X, AlertCircle } from "lucide-react";
import MessageList from "../../components/inbox/MessageList";
import TakeoverButton from "../../components/inbox/TakeoverButton";
import StatusBadge from "../../components/inbox/StatusBadge";
import { inboxService, Conversation, Message } from "../../services/inbox";
import { useWebSocket } from "../../hooks/useWebSocket";
import { useAuth } from "../../context/AuthContext";

export default function Inbox() {
  const { token, user } = useAuth();
  const [conversations, setConversations] = useState<Conversation[]>([]);
  const [activeConvId, setActiveConvId] = useState<string | null>(null);
  const activeConvIdRef = useRef<string | null>(null);
  const [humanMode, setHumanMode] = useState(false);
  const [replyText, setReplyText] = useState("");
  const [streamingContent, setStreamingContent] = useState("");
  const streamingRef = useRef("");
  const [loadingConvs, setLoadingConvs] = useState(true);
  const [showNewChat, setShowNewChat] = useState(false);
  const [newChatTitle, setNewChatTitle] = useState("");
  const [creatingChat, setCreatingChat] = useState(false);
  const messagesEndRef = useRef<HTMLDivElement>(null);
  const [wsError, setWsError] = useState<string | null>(null);

  // Keep refs in sync
  useEffect(() => { activeConvIdRef.current = activeConvId; }, [activeConvId]);

  // Load conversations from real API
  useEffect(() => {
    inboxService
      .getConversations()
      .then((data) => {
        setConversations(data);
        if (data.length > 0) setActiveConvId(data[0].id);
      })
      .catch((err) => console.error("Failed to load conversations:", err))
      .finally(() => setLoadingConvs(false));
  }, []);

  const handleCreateChat = async () => {
    if (!newChatTitle.trim()) return;
    setCreatingChat(true);
    try {
      const newConv = await inboxService.createSession(newChatTitle.trim());
      setConversations((prev) => [newConv, ...prev]);
      setActiveConvId(newConv.id);
      setShowNewChat(false);
      setNewChatTitle("");
    } catch (err) {
      console.error("Failed to create session:", err);
    } finally {
      setCreatingChat(false);
    }
  };

  // Load messages when switching conversations
  useEffect(() => {
    if (!activeConvId) return;
    inboxService.getMessages(activeConvId).then((msgs) => {
      setConversations((prev) =>
        prev.map((c) =>
          c.id === activeConvId ? { ...c, messages: msgs } : c
        )
      );
    });
  }, [activeConvId]);

  // Auto-scroll to bottom
  useEffect(() => {
    messagesEndRef.current?.scrollIntoView({ behavior: "smooth" });
  }, [conversations, streamingContent]);

  // WebSocket for streaming AI responses
  // Uses refs to avoid recreating the callback on every token (which would disconnect WS)
  const handleToken = useCallback(
    (chunk: { content: string; is_last: boolean }) => {
      if (chunk.is_last) {
        const fullText = streamingRef.current;
        streamingRef.current = "";
        setStreamingContent("");
        if (!fullText.trim()) return;
        const finalMsg: Message = {
          id: Date.now().toString(),
          sender: "bot",
          text: fullText,
          timestamp: new Date().toISOString(),
        };
        setConversations((prev) =>
          prev.map((c) =>
            c.id === activeConvIdRef.current
              ? { ...c, messages: [...c.messages, finalMsg], lastMessage: finalMsg.text }
              : c
          )
        );
      } else {
        streamingRef.current += chunk.content;
        setStreamingContent((prev) => prev + chunk.content);
      }
    },
    [] // stable — uses refs only
  );

  const { isConnected, sendMessage: wsSendMessage } = useWebSocket({
    sessionId: activeConvId,
    orgId: user?.org_id || null,
    token,
    onToken: handleToken,
    onError: (msg) => {
      console.error("WS Error:", msg);
      setWsError(msg);
      setTimeout(() => setWsError(null), 6000);
    },
  });

  const activeConv = conversations.find((c) => c.id === activeConvId);

  const handleSendMessage = async () => {
    if (!replyText.trim() || !activeConvId) return;

    const msgText = replyText;
    setReplyText("");

    if (humanMode) {
      // Human agent sends via REST API
      try {
        const newMessage = await inboxService.sendMessage(activeConvId, msgText);
        setConversations((prev) =>
          prev.map((conv) =>
            conv.id === activeConvId
              ? { ...conv, messages: [...conv.messages, newMessage], lastMessage: newMessage.text }
              : conv
          )
        );
      } catch (error) {
        console.error("Failed to send message", error);
      }
    } else {
      // Bot mode: send via WebSocket for AI processing
      const userMsg: Message = {
        id: Date.now().toString(),
        sender: "user",
        text: msgText,
        timestamp: new Date().toISOString(),
      };
      setConversations((prev) =>
        prev.map((conv) =>
          conv.id === activeConvId
            ? { ...conv, messages: [...conv.messages, userMsg], lastMessage: msgText }
            : conv
        )
      );
      wsSendMessage(msgText);
    }
  };

  const handleTakeover = async () => {
    if (!activeConvId) return;
    try {
      if (!humanMode) {
        // Switch to human mode — update DB
        await inboxService.takeoverSession(activeConvId);
        setHumanMode(true);
      } else {
        // Return to bot mode — reset status in DB
        await inboxService.returnToBotMode(activeConvId);
        setHumanMode(false);
      }
    } catch (err) {
      console.error("Failed to toggle session mode:", err);
    }
  };

  return (
    <>
      <PageBreadcrumb pageTitle="Inbox" />

      <div className="flex flex-col lg:flex-row gap-6 h-[calc(100vh-160px)] min-h-[500px] overflow-hidden">

        {/* ── Sidebar ─────────────────────────────────────────────── */}
        <div className="lg:w-1/3 flex flex-col min-h-0 overflow-hidden rounded-2xl border border-gray-200 dark:border-gray-800 bg-white dark:bg-white/[0.03]">
          {/* Sidebar header */}
          <div className="px-5 py-4 border-b border-gray-100 dark:border-gray-800 flex-shrink-0">
            <div className="flex justify-between items-center">
              <span className="text-sm font-medium text-gray-800 dark:text-white/90">Conversations</span>
              <div className="flex items-center gap-2">
                <span className="text-xs text-gray-400">{conversations.length} percakapan</span>
                <Button size="sm" variant="primary" onClick={() => setShowNewChat(true)} className="px-3 py-1.5 text-xs flex items-center gap-1">
                  <Plus size={14} /> New Chat
                </Button>
              </div>
            </div>

            {/* New Chat input */}
            {showNewChat && (
              <div className="mt-3 p-3 rounded-lg border border-brand-500/30 bg-brand-500/5">
                <div className="flex items-center justify-between mb-2">
                  <span className="text-xs font-medium text-gray-700 dark:text-gray-300">Nama percakapan baru</span>
                  <button onClick={() => setShowNewChat(false)} className="text-gray-400 hover:text-gray-600"><X size={14} /></button>
                </div>
                <div className="flex gap-2">
                  <input
                    autoFocus type="text" value={newChatTitle}
                    onChange={(e) => setNewChatTitle(e.target.value)}
                    onKeyDown={(e) => e.key === "Enter" && handleCreateChat()}
                    placeholder="Contoh: Visitor dari Website"
                    className="flex-1 h-8 px-3 text-xs rounded-lg border border-gray-300 dark:border-gray-700 bg-transparent text-gray-800 dark:text-white focus:outline-none focus:ring-2 focus:ring-brand-500/20"
                  />
                  <Button size="sm" variant="primary" onClick={handleCreateChat} disabled={creatingChat || !newChatTitle.trim()} className="px-3 py-1 text-xs">
                    {creatingChat ? "..." : "Buat"}
                  </Button>
                </div>
              </div>
            )}
          </div>

          {/* Conversation list — scrollable */}
          <div className="flex-1 overflow-y-auto">
            {loadingConvs ? (
              <div className="flex items-center justify-center h-full text-sm text-gray-400">Memuat percakapan...</div>
            ) : conversations.length === 0 ? (
              <div className="flex items-center justify-center h-full text-sm text-gray-400">Belum ada percakapan.</div>
            ) : (
              conversations.map((conv) => {
                const isActive = activeConvId === conv.id;
                return (
                  <div
                    key={conv.id}
                    onClick={() => setActiveConvId(conv.id)}
                    className={`flex gap-3 p-4 cursor-pointer transition-colors border-b border-gray-100 dark:border-gray-800 ${
                      isActive ? "bg-gray-100 dark:bg-white/5" : "hover:bg-gray-50 dark:hover:bg-white/5"
                    }`}
                  >
                    <span className="relative block w-10 h-10 rounded-full flex-shrink-0">
                      <img width={40} height={40} src={conv.avatar || "/images/user/user-01.jpg"} alt={conv.user} className="w-full h-full object-cover rounded-full" />
                      <StatusBadge status={conv.status || "online"} />
                    </span>
                    <span className="block flex-1 min-w-0">
                      <span className="flex justify-between items-center mb-0.5">
                        <span className="font-medium text-sm text-gray-800 dark:text-white/90 truncate">{conv.user}</span>
                        {conv.unread > 0 && (
                          <span className="flex items-center justify-center bg-brand-500 text-white text-[10px] px-2 py-0.5 rounded-full flex-shrink-0">{conv.unread}</span>
                        )}
                      </span>
                      <span className="block text-xs text-gray-500 dark:text-gray-400 truncate">{conv.lastMessage}</span>
                    </span>
                  </div>
                );
              })
            )}
          </div>
        </div>

        {/* ── Chat Area ────────────────────────────────────────────── */}
        <div className="lg:w-2/3 flex flex-col min-h-0 overflow-hidden rounded-2xl border border-gray-200 dark:border-gray-800 bg-white dark:bg-white/[0.03]">
          {!activeConv ? (
            <div className="flex-1 flex items-center justify-center text-sm text-gray-400">Pilih percakapan</div>
          ) : (
            <>
              {/* Chat header */}
              <div className="px-5 py-4 border-b border-gray-100 dark:border-gray-800 flex-shrink-0 flex justify-between items-center">
                <div>
                  <h3 className="text-sm font-medium text-gray-800 dark:text-white/90">Chat with {activeConv.user}</h3>
                  <div className="flex items-center gap-1.5 mt-0.5 text-xs text-gray-500">
                    {isConnected ? (
                      <><Wifi size={12} className="text-green-500" /><span className="text-green-500">AI Terhubung</span></>
                    ) : (
                      <><WifiOff size={12} className="text-red-400" /><span className="text-red-400">Offline</span></>
                    )}
                  </div>
                </div>
                <TakeoverButton humanMode={humanMode} onClick={handleTakeover} />
              </div>

              {/* Messages — scrollable area */}
              <div className="flex-1 overflow-y-auto px-4 py-4">
                <MessageList messages={activeConv.messages} />
                {streamingContent && (
                  <div className="flex gap-2 px-2 py-1">
                    <span className="text-xs text-gray-400 italic">Bot: {streamingContent}▊</span>
                  </div>
                )}
                <div ref={messagesEndRef} />
              </div>

              {/* WS Error Toast */}
              {wsError && (
                <div className="mx-4 mb-2 flex items-start gap-2 rounded-lg bg-red-500/10 border border-red-500/20 px-3 py-2 text-xs text-red-400 flex-shrink-0">
                  <AlertCircle size={14} className="mt-0.5 flex-shrink-0" />
                  <span>{wsError}</span>
                </div>
              )}

              {/* Input */}
              <div className="px-4 py-4 border-t border-gray-100 dark:border-gray-800 flex gap-3 flex-shrink-0">
                <input
                  type="text" value={replyText}
                  onChange={(e) => setReplyText(e.target.value)}
                  placeholder={humanMode ? "Tulis pesan sebagai agen..." : "Tulis pesan ke AI..."}
                  data-testid="message-input"
                  className="flex-1 h-11 px-4 rounded-lg border border-gray-300 dark:border-gray-700 bg-transparent text-gray-800 dark:text-white focus:outline-none focus:ring-2 focus:ring-brand-500/20"
                  onKeyDown={(e) => { if (e.key === "Enter" && replyText.trim()) handleSendMessage(); }}
                />
                <Button variant="primary" disabled={!replyText.trim()} className="px-4" onClick={handleSendMessage} data-testid="send-button">
                  <Send size={18} />
                </Button>
              </div>
            </>
          )}
        </div>
      </div>
    </>
  );
}