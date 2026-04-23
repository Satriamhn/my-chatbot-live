import { useState } from "react";
import PageBreadcrumb from "../../components/common/PageBreadCrumb";
import ComponentCard from "../../components/common/ComponentCard";
import Input from "../../components/form/input/InputField";
import Label from "../../components/form/Label";
import Button from "../../components/ui/button/Button";
import { Send, Bot } from "lucide-react";

export default function BotSettings() {
  const [botName, setBotName] = useState("Acme Helper");
  const [themeColor, setThemeColor] = useState("#3b82f6");
  const [welcomeMessage, setWelcomeMessage] = useState("Hi! How can I help you today?");

  return (
    <>
      <PageBreadcrumb pageTitle="Bot Configuration" />

      <div className="grid grid-cols-1 gap-6 lg:grid-cols-2">
        {/* Settings Form */}
        <div className="flex flex-col gap-6">
          <ComponentCard title="General Settings">
            <div className="flex flex-col gap-4">
              <div>
                <Label htmlFor="botName">Bot Name</Label>
                <Input
                  id="botName"
                  type="text"
                  value={botName}
                  onChange={(e) => setBotName(e.target.value)}
                  placeholder="Enter bot name"
                />
              </div>

              <div>
                <Label htmlFor="welcomeMessage">Welcome Message</Label>
                <Input
                  id="welcomeMessage"
                  type="text"
                  value={welcomeMessage}
                  onChange={(e) => setWelcomeMessage(e.target.value)}
                  placeholder="Enter welcome message"
                />
              </div>

              <div>
                <Label htmlFor="themeColor">Theme Color</Label>
                <div className="flex items-center gap-3">
                  <input
                    type="color"
                    id="themeColor"
                    value={themeColor}
                    onChange={(e) => setThemeColor(e.target.value)}
                    className="h-11 w-14 rounded-lg cursor-pointer border border-gray-200 dark:border-gray-800"
                  />
                  <Input
                    type="text"
                    value={themeColor}
                    onChange={(e) => setThemeColor(e.target.value)}
                    className="flex-1"
                  />
                </div>
              </div>

              <div className="mt-2">
                <Button variant="primary">Save Changes</Button>
              </div>
            </div>
          </ComponentCard>
        </div>

        {/* Live Preview */}
        <div className="flex flex-col gap-6">
          <ComponentCard title="Live Preview" desc="See how your bot appears to users.">
            <div className="flex justify-center items-center p-6 bg-gray-50 dark:bg-gray-800/50 rounded-xl min-h-[400px]">
              
              {/* Mock Chat Widget */}
              <div className="w-[320px] bg-white dark:bg-gray-900 rounded-2xl shadow-xl overflow-hidden border border-gray-200 dark:border-gray-800 flex flex-col">
                {/* Header */}
                <div 
                  className="px-4 py-3 flex items-center gap-3 text-white"
                  style={{ backgroundColor: themeColor }}
                >
                  <div className="w-8 h-8 rounded-full bg-white/20 flex items-center justify-center">
                    <Bot size={18} />
                  </div>
                  <div>
                    <h4 className="font-medium text-sm">{botName}</h4>
                    <p className="text-xs text-white/80">Online</p>
                  </div>
                </div>

                {/* Chat Body */}
                <div className="p-4 flex-1 h-[240px] bg-gray-50 dark:bg-gray-800 flex flex-col gap-3">
                  {/* Bot Message */}
                  <div className="flex gap-2">
                    <div className="w-6 h-6 rounded-full flex items-center justify-center flex-shrink-0 text-white" style={{ backgroundColor: themeColor }}>
                      <Bot size={14} />
                    </div>
                    <div className="bg-white dark:bg-gray-900 border border-gray-100 dark:border-gray-700 p-2.5 rounded-lg rounded-tl-none text-sm text-gray-700 dark:text-gray-300 shadow-sm">
                      {welcomeMessage}
                    </div>
                  </div>
                </div>

                {/* Input Area */}
                <div className="p-3 bg-white dark:bg-gray-900 border-t border-gray-100 dark:border-gray-800 flex items-center gap-2">
                  <input 
                    type="text" 
                    placeholder="Type a message..." 
                    className="flex-1 text-sm bg-transparent border-none focus:outline-none dark:text-white"
                    disabled
                  />
                  <button 
                    className="p-1.5 rounded-md text-white flex items-center justify-center transition-opacity hover:opacity-90"
                    style={{ backgroundColor: themeColor }}
                    disabled
                  >
                    <Send size={16} />
                  </button>
                </div>
              </div>

            </div>
          </ComponentCard>
        </div>
      </div>
    </>
  );
}
