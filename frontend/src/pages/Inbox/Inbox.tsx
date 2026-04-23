import { useState, useEffect } from "react";
import PageBreadcrumb from "../../components/common/PageBreadCrumb";
import ComponentCard from "../../components/common/ComponentCard";
import Button from "../../components/ui/button/Button";
import { User, Bot, UserCog, Send, Power } from "lucide-react";

interface Message {
  id: string;
  sender: "user" | "bot" | "human";
  text: string;
  timestamp?: string;
  meta?: { source: string; confidence: string };
}

interface Conversation {
  id: string;
  user: string;
  lastMessage: string;
  unread: number;
  messages: Message[];
}

export default function Inbox() {
  const [conversations, setConversations] = useState<Conversation[]>([]);
  const [activeConvId, setActiveConvId] = useState<string | null>(null);
  const [humanMode, setHumanMode] = useState(false);
  const [replyText, setReplyText] = useState("");

  useEffect(() => {
    fetch("/src/mocks/inbox.json")
      .then(res => res.json())
      .then(data => {
        setConversations(data.conversations);
        if (data.conversations.length > 0) {
          setActiveConvId(data.conversations[0].id);
        }
      })
      .catch(err => console.error(err));
  }, []);

  const activeConv = conversations.find(c => c.id === activeConvId);

  return (
    <>
      <PageBreadcrumb pageTitle="Inbox" />

      <div className="flex flex-col lg:flex-row gap-6 h-[calc(100vh-160px)] min-h-[600px]">
        
        {/* Sidebar */}
        <div className="lg:w-1/3 flex flex-col gap-4">
          <ComponentCard title="Conversations" className="h-full flex flex-col">
            <div className="flex-1 overflow-y-auto space-y-2 -mx-4 px-4">
              {conversations.map((conv) => (
                <div 
                  key={conv.id}
                  onClick={() => setActiveConvId(conv.id)}
                  className={`p-4 rounded-xl cursor-pointer transition-colors border ${
                    activeConvId === conv.id 
                      ? 'border-brand-500 bg-brand-50 dark:bg-brand-500/10' 
                      : 'border-transparent hover:bg-gray-50 dark:hover:bg-gray-800'
                  }`}
                >
                  <div className="flex justify-between items-center mb-1">
                    <h4 className="font-semibold text-gray-800 dark:text-white/90">{conv.user}</h4>
                    {conv.unread > 0 && (
                      <span className="bg-brand-500 text-white text-xs px-2 py-0.5 rounded-full">
                        {conv.unread}
                      </span>
                    )}
                  </div>
                  <p className="text-sm text-gray-500 dark:text-gray-400 truncate">
                    {conv.lastMessage}
                  </p>
                </div>
              ))}
            </div>
          </ComponentCard>
        </div>

        {/* Chat Area */}
        <div className="lg:w-2/3 flex flex-col">
          <ComponentCard title={activeConv ? `Chat with ${activeConv.user}` : "Select a conversation"} className="h-full flex flex-col flex-1">
            {activeConv && (
              <div className="flex flex-col h-full">
                
                {/* Chat Header Actions */}
                <div className="flex justify-end mb-4 border-b border-gray-100 dark:border-gray-800 pb-4">
                  <Button 
                    variant={humanMode ? "primary" : "outline"}
                    onClick={() => setHumanMode(!humanMode)}
                    startIcon={<Power size={16} />}
                  >
                    {humanMode ? "Human Mode Active" : "Take Over (Human Mode)"}
                  </Button>
                </div>

                {/* Messages */}
                <div className="flex-1 overflow-y-auto space-y-6 mb-4 pr-2">
                  {activeConv.messages.map((msg: Message) => {
                    const isBot = msg.sender === 'bot';
                    const isHumanModeMsg = msg.sender === 'human';

                    return (
                      <div key={msg.id} className={`flex gap-3 ${isBot || isHumanModeMsg ? 'justify-start' : 'justify-end'}`}>
                        {(isBot || isHumanModeMsg) && (
                          <div className={`w-8 h-8 rounded-full flex items-center justify-center flex-shrink-0 text-white ${isHumanModeMsg ? 'bg-orange-500' : 'bg-brand-500'}`}>
                            {isHumanModeMsg ? <UserCog size={16} /> : <Bot size={16} />}
                          </div>
                        )}
                        
                        <div className={`flex flex-col max-w-[75%] ${isBot || isHumanModeMsg ? 'items-start' : 'items-end'}`}>
                          <div className={`p-3 rounded-2xl shadow-sm text-sm ${
                            isBot || isHumanModeMsg
                              ? 'bg-gray-100 dark:bg-gray-800 text-gray-800 dark:text-white/90 rounded-tl-none' 
                              : 'bg-brand-500 text-white rounded-tr-none'
                          }`}>
                            {msg.text}
                          </div>
                          
                          {/* Metadata for bot */}
                          {isBot && msg.meta && (
                            <div className="mt-1 text-xs text-gray-400 flex gap-2 px-1">
                              <span>Source: {msg.meta.source}</span>
                              <span>•</span>
                              <span>Confidence: {msg.meta.confidence}</span>
                            </div>
                          )}
                          {!isBot && (
                            <div className="mt-1 text-xs text-gray-400 px-1">
                              {msg.timestamp}
                            </div>
                          )}
                        </div>

                        {!isBot && !isHumanModeMsg && (
                          <div className="w-8 h-8 rounded-full flex items-center justify-center flex-shrink-0 bg-gray-200 dark:bg-gray-700 text-gray-600 dark:text-gray-300">
                            <User size={16} />
                          </div>
                        )}
                      </div>
                    );
                  })}
                </div>

                {/* Input */}
                <div className="pt-4 border-t border-gray-100 dark:border-gray-800 flex gap-3">
                  <input
                    type="text"
                    value={replyText}
                    onChange={(e) => setReplyText(e.target.value)}
                    placeholder={humanMode ? "Type your message..." : "Bot is handling this. Take over to reply."}
                    disabled={!humanMode}
                    className="flex-1 h-11 px-4 rounded-lg border border-gray-300 dark:border-gray-700 bg-transparent text-gray-800 dark:text-white disabled:opacity-50 focus:outline-none focus:ring-2 focus:ring-brand-500/20"
                    onKeyDown={(e) => {
                      if (e.key === 'Enter' && replyText.trim() && humanMode) {
                        setReplyText("");
                        // Ideally we would add to state here
                      }
                    }}
                  />
                  <Button 
                    variant="primary" 
                    disabled={!humanMode || !replyText.trim()}
                    className="px-4"
                  >
                    <Send size={18} />
                  </Button>
                </div>
              </div>
            )}
          </ComponentCard>
        </div>

      </div>
    </>
  );
}
