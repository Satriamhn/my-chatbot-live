import { useState, useRef, useEffect } from "react";
import { useSearchParams } from "react-router-dom";
import { Bot, Send, Loader2, User } from "lucide-react";
import { useWidgetWebSocket, StreamingToken } from "../../hooks/useWidgetWebSocket";
import axios from "axios";
import ReactMarkdown from "react-markdown";
import { renderToStaticMarkup } from "react-dom/server";
import { widgetService, WidgetSettings } from "../../services/widget";

const API_BASE = import.meta.env.VITE_API_URL || "http://localhost:8080/api/v1";

const FALLBACK_WIDGET_SETTINGS: WidgetSettings = {
	bot_name: "Customer Support",
	welcome_message: "Hello! How can we help you today?",
};

const INVISIBLE_CONTENT_RE = /[\s\u200B-\u200D\u2060\uFEFF]/g;

const hasVisibleBotContent = (content: string) => {
  if (content.replace(INVISIBLE_CONTENT_RE, "").length === 0) {
    return false;
  }

  const renderedMarkup = renderToStaticMarkup(<ReactMarkdown>{content}</ReactMarkdown>);

  if (!renderedMarkup.trim()) {
    return false;
  }

  const probe = document.createElement("div");
  probe.innerHTML = renderedMarkup;

  const renderedText = probe.textContent?.replace(INVISIBLE_CONTENT_RE, "").trim() ?? "";

  if (renderedText.length > 0) {
    return true;
  }

  return Boolean(probe.querySelector("hr, img"));
};

const shouldRenderMessage = (message: Message) => message.sender !== "bot" || hasVisibleBotContent(message.content);

interface Message {
  id: string;
  sender: "user" | "bot";
  content: string;
}

interface MessageBubbleProps {
  message: Message;
}

export function MessageBubble({ message }: MessageBubbleProps) {
  if (!shouldRenderMessage(message)) {
    return null;
  }

  return (
    <div
      className={`flex min-w-0 gap-2.5 max-w-[90%] ${message.sender === "user" ? "ml-auto flex-row-reverse" : "mr-auto"}`}
    >
      <div
        className={`w-7 h-7 rounded-full flex items-center justify-center shrink-0 mt-1 text-white ${message.sender === "user" ? "bg-gray-800" : "bg-brand-500"}`}
      >
        {message.sender === "user" ? <User size={14} /> : <Bot size={14} />}
      </div>
      <div
        className={`px-3.5 py-2.5 rounded-2xl overflow-hidden ${message.sender === "user" ? "bg-gray-800 text-white rounded-tr-sm" : "bg-white border border-gray-100 text-gray-700 rounded-tl-sm shadow-sm"}`}
      >
        {message.sender === "bot" ? (
          <div className="prose prose-sm max-w-full prose-p:leading-relaxed prose-pre:bg-gray-100 prose-pre:text-gray-800 break-words overflow-x-auto">
            <ReactMarkdown>{message.content || ""}</ReactMarkdown>
          </div>
        ) : (
          <span className="whitespace-pre-wrap break-words">{message.content || ""}</span>
        )}
      </div>
    </div>
  );
}

export default function Widget() {
  const [searchParams] = useSearchParams();
  const orgId: string | null = searchParams.get("org_id");

  const [widgetSettings, setWidgetSettings] = useState<WidgetSettings>(FALLBACK_WIDGET_SETTINGS);
  const [sessionId, setSessionId] = useState<string | null>(null);
  const [messages, setMessages] = useState<Message[]>([]);
  const [input, setInput] = useState("");
  const [isInitializing, setIsInitializing] = useState(false);
  const [error, setError] = useState<string | null>(null);

  // Streaming state
  const streamingContentRef = useRef("");
  const [streamingRender, setStreamingRender] = useState("");
  const messagesEndRef = useRef<HTMLDivElement>(null);

  const { isConnected, sendMessage: sendWSMessage } = useWidgetWebSocket({
    sessionId,
    orgId,
    onToken: (chunk: StreamingToken) => {
      streamingContentRef.current += chunk.content;
      setStreamingRender(streamingContentRef.current);

      if (chunk.is_last) {
        if (hasVisibleBotContent(streamingContentRef.current)) {
          setMessages((prev) => [
            ...prev,
            { id: Date.now().toString(), sender: "bot", content: streamingContentRef.current },
          ]);
        }
        streamingContentRef.current = "";
        setStreamingRender("");
      }
    },
    onError: (msg) => {
      setError(msg);
      // Remove temporary bot "thinking" message if present and error occurred
      streamingContentRef.current = "";
      setStreamingRender("");
    },
  });

  const scrollToBottom = () => {
    messagesEndRef.current?.scrollIntoView({ behavior: "smooth" });
  };

  useEffect(() => {
    scrollToBottom();
  }, [messages, streamingRender]);

  useEffect(() => {
    document.title = widgetSettings.bot_name;
  }, [widgetSettings.bot_name]);

  useEffect(() => {
    let mounted = true;
    const loadWidgetSettings = async () => {
      if (!orgId) return;

      try {
        const settings = await widgetService.getSettings(orgId);
        if (mounted) {
          setWidgetSettings({
            bot_name: settings.bot_name?.trim() || FALLBACK_WIDGET_SETTINGS.bot_name,
            welcome_message: settings.welcome_message?.trim() || FALLBACK_WIDGET_SETTINGS.welcome_message,
          });
        }
      } catch {
        if (mounted) {
          setWidgetSettings(FALLBACK_WIDGET_SETTINGS);
        }
      }
    };

    loadWidgetSettings();

    return () => {
      mounted = false;
    };
  }, [orgId]);

  useEffect(() => {
    let mounted = true;
    const initSessionOnLoad = async () => {
      if (!orgId || sessionId) return;
      setIsInitializing(true);
      try {
        const res = await axios.post(`${API_BASE}/widget/sessions`, {}, {
          headers: { "X-Org-ID": orgId }
        });
        if (mounted) {
          setSessionId(res.data.id);
        }
      } catch (err: unknown) {
        if (mounted) {
          const message = axios.isAxiosError(err) ? err.response?.data?.error : null;
          setError(message || "Failed to initialize chat session.");
        }
      } finally {
        if (mounted) setIsInitializing(false);
      }
    };
    initSessionOnLoad();
    return () => { mounted = false; };
  }, [orgId, sessionId]);

  const handleSend = async (e?: React.FormEvent) => {
    e?.preventDefault();
    if (!input.trim() || !orgId || !sessionId) return;

    const msgContent = input;
    setInput("");
    setError(null);

    // Optimistically add user message
    setMessages((prev) => [...prev, { id: Date.now().toString(), sender: "user", content: msgContent }]);

    // Wait a brief moment for WS to connect if it was just initialized
    let attempts = 0;
    const checkAndSend = () => {
      if (sendWSMessage(msgContent)) return;
      if (attempts < 10) {
        attempts++;
        setTimeout(checkAndSend, 200);
      } else {
        setError("Failed to connect to chat server.");
      }
    };
    checkAndSend();
  };

  if (!orgId) {
    return (
      <div className="flex h-screen items-center justify-center bg-gray-50 p-4 text-center">
        <div role="alert" className="max-w-sm rounded-2xl border border-red-100 bg-white px-6 py-5 text-sm text-gray-700 shadow-sm">
          <div className="font-semibold text-red-700">Missing org_id</div>
          <div className="mt-2 leading-6">This widget needs a tenant id in the URL to load.</div>
        </div>
      </div>
    );
  }

  return (
    <div className="flex flex-col h-screen bg-white font-sans text-sm overflow-hidden">
      {/* Header */}
      <div className="flex items-center gap-3 px-4 py-3 bg-brand-500 text-white shadow-sm shrink-0">
        <div className="w-8 h-8 bg-white/20 rounded-full flex items-center justify-center">
          <Bot size={18} />
        </div>
        <div>
          <h3 className="font-semibold leading-tight">{widgetSettings.bot_name}</h3>
          <p className="text-[11px] opacity-90 flex items-center gap-1.5">
            <span className={`w-1.5 h-1.5 rounded-full ${isConnected ? "bg-green-400" : "bg-white/40"}`} />
            {isConnected ? "Online" : "Connecting..."}
          </p>
        </div>
      </div>

      {/* Messages */}
      <div className="flex-1 overflow-y-auto p-4 flex flex-col gap-4 bg-gray-50/50">
        {messages.length === 0 && (
          <div className="mx-auto mt-10 max-w-sm rounded-3xl border border-brand-100 bg-white px-5 py-6 text-center shadow-sm">
            <div className="mx-auto mb-3 flex h-10 w-10 items-center justify-center rounded-full bg-brand-50 text-brand-600">
              <Bot size={18} />
            </div>
            <h4 className="text-sm font-semibold text-gray-900">{widgetSettings.bot_name}</h4>
            <p className="mt-2 text-sm leading-6 text-gray-600">{widgetSettings.welcome_message}</p>
            <p className="mt-3 text-xs text-gray-400">Send a message to start chatting.</p>
          </div>
        )}

        {messages.map((msg) =>
          shouldRenderMessage(msg) ? <MessageBubble key={msg.id} message={msg} /> : null,
        )}

        {hasVisibleBotContent(streamingRender) && (
          <div className="flex min-w-0 gap-2.5 max-w-[90%] mr-auto">
             <div className="w-7 h-7 rounded-full flex items-center justify-center shrink-0 mt-1 text-white bg-brand-500">
               <Bot size={14} />
             </div>
             <div className="px-3.5 py-2.5 rounded-2xl overflow-hidden bg-white border border-gray-100 text-gray-700 rounded-tl-sm shadow-sm min-w-0">
               <div className="prose prose-sm max-w-full prose-p:leading-relaxed prose-pre:bg-gray-100 prose-pre:text-gray-800 break-words overflow-x-auto">
                 <ReactMarkdown>
                   {streamingRender + " ▊"}
                 </ReactMarkdown>
               </div>
             </div>
          </div>
        )}
        <div ref={messagesEndRef} />
      </div>

      {/* Error Bar */}
      {error && (
        <div role="alert" className="px-4 py-2 bg-red-50 text-red-600 text-xs text-center border-t border-red-100 shrink-0">
          {error}
        </div>
      )}

      {/* Input */}
      <div className="p-3 bg-white border-t border-gray-100 shrink-0">
        <form onSubmit={handleSend} className="flex items-center gap-2 relative">
          <input
            type="text"
            value={input}
            onChange={(e) => setInput(e.target.value)}
            disabled={isInitializing}
            placeholder="Type your message..."
            className="flex-1 bg-gray-50 border border-gray-200 rounded-full px-4 py-2.5 text-sm focus:outline-none focus:border-brand-300 focus:ring-2 focus:ring-brand-100 disabled:opacity-50"
          />
          <button
            type="submit"
            disabled={!input.trim() || isInitializing}
            className="w-10 h-10 rounded-full bg-brand-500 text-white flex items-center justify-center disabled:opacity-50 disabled:cursor-not-allowed hover:bg-brand-600 transition-colors shrink-0"
          >
            {isInitializing ? <Loader2 size={16} className="animate-spin" /> : <Send size={16} className="ml-0.5" />}
          </button>
        </form>
      </div>
    </div>
  );
}
