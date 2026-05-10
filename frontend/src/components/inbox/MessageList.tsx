import { Bot, User, UserCog } from 'lucide-react';
import ReactMarkdown from 'react-markdown';

export interface Message {
  id: string;
  sender: "user" | "bot" | "human_agent" | "system";
  text: string;
  timestamp?: string;
  meta?: { source: string; confidence: string };
}

interface MessageListProps {
  messages: Message[];
}

export default function MessageList({ messages }: MessageListProps) {
  return (
    <div className="flex-1 overflow-y-auto space-y-6 mb-4 pr-2 custom-scrollbar">
      {messages.map((msg: Message) => {
        const isBot = msg.sender === 'bot';
        const isHumanModeMsg = msg.sender === 'human_agent';

        return (
          <div key={msg.id} data-testid={`message-${msg.sender}`} className={`flex gap-3 ${isBot || isHumanModeMsg ? 'justify-start' : 'justify-end'}`}>
            {(isBot || isHumanModeMsg) && (
              <div className={`w-8 h-8 rounded-full flex items-center justify-center flex-shrink-0 text-white ${isHumanModeMsg ? 'bg-orange-500' : 'bg-brand-500'}`}>
                {isHumanModeMsg ? <UserCog size={16} /> : <Bot size={16} />}
              </div>
            )}
            
            <div className={`flex flex-col min-w-0 flex-1 max-w-[85%] lg:max-w-[75%] ${isBot || isHumanModeMsg ? 'items-start' : 'items-end'}`}>
              <div className={`p-3 rounded-2xl shadow-sm text-sm overflow-hidden ${
                isBot || isHumanModeMsg
                  ? 'bg-gray-100 dark:bg-gray-800 text-gray-800 dark:text-white/90 rounded-tl-none' 
                  : 'bg-brand-500 text-white rounded-tr-none'
              }`}>
                {(isBot || isHumanModeMsg) ? (
                  <div className="prose prose-sm dark:prose-invert max-w-full prose-p:leading-relaxed prose-pre:bg-gray-900 prose-pre:text-white break-words [&_*]:break-words">
                    <ReactMarkdown>
                      {msg.text || ""}
                    </ReactMarkdown>
                  </div>
                ) : (
                  <span className="whitespace-pre-wrap break-words">{msg.text || ""}</span>
                )}
              </div>
              
              {/* Metadata for bot */}
              {isBot && msg.meta && (
                <div data-testid="bot-metadata" className="mt-1 text-xs text-gray-400 flex gap-2 px-1">
                  <span>Source: {msg.meta.source}</span>
                  <span>•</span>
                  <span>Confidence: {msg.meta.confidence}</span>
                </div>
              )}
              {!isBot && msg.timestamp && (
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
  );
}