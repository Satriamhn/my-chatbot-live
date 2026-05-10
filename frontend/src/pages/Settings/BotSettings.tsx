import { useState, useEffect } from "react";
import PageBreadcrumb from "../../components/common/PageBreadCrumb";
import ComponentCard from "../../components/common/ComponentCard";
import Input from "../../components/form/input/InputField";
import TextArea from "../../components/form/input/TextArea";
import Label from "../../components/form/Label";
import Select from "../../components/form/Select";
import Button from "../../components/ui/button/Button";
import { Send, Bot, Loader2, Sparkles, Key, Zap, ShieldCheck } from "lucide-react";
import { settingsService, AIProvider, PROVIDER_MODELS } from "../../services/settings";

export default function BotSettings() {
  // Bot config
  const [botName, setBotName] = useState("Chatbot");
  const [themeColor, setThemeColor] = useState("#3b82f6");
  const [welcomeMessage, setWelcomeMessage] = useState("Hi! How can I help you today?");
  const [systemPrompt, setSystemPrompt] = useState("You are a helpful assistant.");

  // AI Provider (hybrid)
  const [aiProvider, setAiProvider] = useState<AIProvider>("gemini");
  const [modelName, setModelName] = useState("gemini-2.0-flash");
  const [apiKey, setApiKey] = useState(""); // BYOK input
  const [showApiKey, setShowApiKey] = useState(false);
  const [hasByokKey, setHasByokKey] = useState(false);
  const [dailyCount, setDailyCount] = useState(0);
  const [dailyLimit, setDailyLimit] = useState(100);

  // UI state
  const [loading, setLoading] = useState(true);
  const [saving, setSaving] = useState(false);
  const [saved, setSaved] = useState(false);

  useEffect(() => {
    let mounted = true;
    settingsService.getSettings()
      .then((data) => {
        if (!mounted) return;
        if (data.bot_name) setBotName(data.bot_name);
        if (data.welcome_message) setWelcomeMessage(data.welcome_message);
        if (data.system_prompt) setSystemPrompt(data.system_prompt);
        setAiProvider("gemini");
        if (data.ai_provider === "gemini" && data.model_name) {
          setModelName(data.model_name);
        } else {
          setModelName(PROVIDER_MODELS.gemini.models[0].value);
        }
        setHasByokKey(data.has_byok_key ?? false);
        setDailyCount(data.daily_message_count ?? 0);
        setDailyLimit(data.daily_message_limit ?? 100);
      })
      .catch(console.error)
      .finally(() => { if (mounted) setLoading(false); });
    return () => { mounted = false; };
  }, []);

  // Reset model when provider changes
  const handleProviderChange = (p: AIProvider) => {
    setAiProvider(p);
    setModelName(PROVIDER_MODELS[p].models[0].value);
  };

  const handleSave = async () => {
    setSaving(true);
    try {
      const payload: Parameters<typeof settingsService.updateSettings>[0] = {
        bot_name: botName,
        welcome_message: welcomeMessage,
        system_prompt: systemPrompt,
        ai_provider: aiProvider,
        model_name: modelName,
      };
      if (apiKey.trim()) {
        payload.api_key = apiKey.trim();
      }
      const updated = await settingsService.updateSettings(payload);
      setHasByokKey(updated.has_byok_key);
      setApiKey(""); // clear after save
      setSaved(true);
      setTimeout(() => setSaved(false), 2500);
    } catch (err) {
      console.error("Failed to save settings:", err);
    } finally {
      setSaving(false);
    }
  };

  if (loading) {
    return (
      <div className="flex items-center justify-center h-[calc(100vh-200px)]">
        <Loader2 className="w-8 h-8 animate-spin text-brand-500" />
      </div>
    );
  }

  const usagePercent = dailyLimit > 0 ? Math.min((dailyCount / dailyLimit) * 100, 100) : 0;

  return (
    <>
      <PageBreadcrumb pageTitle="Bot Configuration" />

      <div className="grid grid-cols-1 gap-8 lg:grid-cols-2 lg:gap-10">
        {/* ─── Left: Settings Form ─── */}
        <div className="flex flex-col gap-6">

          {/* General Settings */}
          <ComponentCard title="General Settings" desc="Konfigurasi identitas dan perilaku chatbot.">
            <div className="flex flex-col gap-5">
              <div>
                <Label htmlFor="botName">Bot Name</Label>
                <Input id="botName" type="text" value={botName} onChange={(e) => setBotName(e.target.value)} placeholder="Contoh: Asisten Toko" />
              </div>

              <div>
                <Label htmlFor="themeColor">Theme Color</Label>
                <div className="flex items-center gap-3">
                  <div className="relative flex items-center justify-center h-11 w-16 rounded-lg overflow-hidden border border-gray-200 shadow-theme-xs dark:border-gray-800 focus-within:ring-2 focus-within:ring-brand-500/20">
                    <input type="color" id="themeColor" value={themeColor} onChange={(e) => setThemeColor(e.target.value)} className="absolute inset-0 w-full h-full cursor-pointer opacity-0" />
                    <div className="w-full h-full pointer-events-none" style={{ backgroundColor: themeColor }} />
                  </div>
                  <Input type="text" value={themeColor} onChange={(e) => setThemeColor(e.target.value)} className="flex-1 font-mono uppercase tracking-wider text-sm" />
                </div>
              </div>

              <div>
                <Label htmlFor="welcomeMessage">Welcome Message</Label>
                <TextArea value={welcomeMessage} onChange={setWelcomeMessage} placeholder="Pesan pertama yang dilihat user" rows={2} hint="Pesan yang muncul saat chat widget dibuka." />
              </div>

              <div>
                <Label htmlFor="systemPrompt">
                  <div className="flex items-center gap-1.5">
                    <Sparkles size={16} className="text-brand-500" />
                    System Prompt
                  </div>
                </Label>
                <TextArea value={systemPrompt} onChange={setSystemPrompt} placeholder="You are a helpful customer support assistant..." rows={4} hint="Instruksi rahasia untuk bot. Tidak ditampilkan ke user." />
              </div>
            </div>
          </ComponentCard>

          {/* AI Provider Settings */}
          <ComponentCard title="AI Provider" desc="Pilih provider AI dan masukkan API key kamu sendiri (opsional).">
            <div className="flex flex-col gap-5">

              {/* Provider selector */}
              <div>
                <Label>Provider</Label>
                <div className="grid grid-cols-2 gap-3 mt-1">
                  {(Object.entries(PROVIDER_MODELS) as [AIProvider, typeof PROVIDER_MODELS[AIProvider]][])
                    .filter(([key]) => key === "gemini")
                    .map(([key, val]) => (
                    <button
                      key={key}
                      onClick={() => handleProviderChange(key)}
                      className={`flex items-center gap-2 p-3 rounded-xl border-2 text-sm font-medium transition-all ${
                        aiProvider === key
                          ? "border-brand-500 bg-brand-50 text-brand-700 dark:bg-brand-500/10 dark:text-brand-400"
                          : "border-gray-200 dark:border-gray-700 text-gray-600 dark:text-gray-400 hover:border-gray-300"
                      }`}
                      >
                        <Zap size={16} />
                        {val.label}
                      </button>
                  ))}
                </div>
              </div>

              {/* Model selector */}
              <div>
                <Label>Model</Label>
                <div className="mt-1">
                  <Select
                    options={PROVIDER_MODELS[aiProvider].models}
                    value={modelName}
                    onChange={(val) => setModelName(val)}
                  />
                </div>
              </div>

              {/* BYOK API Key */}
              <div>
                <Label>
                  <div className="flex items-center gap-1.5">
                    <Key size={16} className="text-brand-500" />
                    API Key Sendiri (BYOK — Opsional)
                  </div>
                </Label>
                <div className="relative">
                  <Input
                    type={showApiKey ? "text" : "password"}
                    placeholder={hasByokKey ? "••••••••••••••••••• (sudah tersimpan)" : "Masukkan API key untuk unlimited usage"}
                    value={apiKey}
                    onChange={(e) => setApiKey(e.target.value)}
                  />
                  <button
                    type="button"
                    onClick={() => setShowApiKey(!showApiKey)}
                    className="absolute right-3 top-1/2 -translate-y-1/2 text-xs text-gray-400 hover:text-gray-600"
                  >
                    {showApiKey ? "Sembunyikan" : "Tampilkan"}
                  </button>
                </div>
                {hasByokKey ? (
                  <p className="mt-1 text-xs text-green-600 dark:text-green-400 flex items-center gap-1">
                    <ShieldCheck size={12} /> API key tersimpan — unlimited usage aktif
                  </p>
                ) : (
                  <p className="mt-1 text-xs text-gray-400">
                    Tanpa API key sendiri, kamu pakai platform key ({dailyCount}/{dailyLimit} pesan hari ini).
                  </p>
                )}
              </div>

              {/* Usage bar (only for platform key users) */}
              {!hasByokKey && (
                <div>
                  <div className="flex justify-between text-xs text-gray-500 mb-1">
                    <span>Penggunaan hari ini</span>
                    <span>{dailyCount} / {dailyLimit} pesan</span>
                  </div>
                  <div className="w-full h-2 bg-gray-100 dark:bg-gray-800 rounded-full overflow-hidden">
                    <div
                      className={`h-full rounded-full transition-all ${usagePercent >= 90 ? "bg-red-500" : usagePercent >= 70 ? "bg-yellow-500" : "bg-brand-500"}`}
                      style={{ width: `${usagePercent}%` }}
                    />
                  </div>
                </div>
              )}
            </div>
          </ComponentCard>

          <div>
            <Button variant="primary" onClick={handleSave} disabled={saving} className="w-full sm:w-auto min-w-[160px] flex justify-center">
              {saving ? <Loader2 size={18} className="animate-spin" /> : saved ? "✓ Tersimpan!" : "Simpan Perubahan"}
            </Button>
          </div>
        </div>

        {/* ─── Right: Live Preview ─── */}
        <div className="flex flex-col gap-6 lg:sticky lg:top-24 self-start">
          <ComponentCard title="Live Preview" desc="Visualisasi real-time chat widget kamu.">
            <div className="flex justify-center items-center p-8 bg-gradient-to-br from-gray-50/50 to-gray-100/80 dark:from-gray-800/20 dark:to-gray-900/40 rounded-xl min-h-[500px] border border-gray-100/50 dark:border-gray-800/50 relative overflow-hidden shadow-inner">
              <div className="absolute top-0 right-0 w-64 h-64 bg-current opacity-[0.03] blur-3xl rounded-full" style={{ color: themeColor }} />
              <div className="absolute bottom-0 left-0 w-48 h-48 bg-current opacity-[0.03] blur-2xl rounded-full" style={{ color: themeColor }} />

              <div className="w-[340px] bg-white dark:bg-gray-900 rounded-[24px] shadow-2xl overflow-hidden border border-gray-200/60 dark:border-gray-800 flex flex-col relative z-10">
                {/* Header */}
                <div className="px-5 py-4 flex items-center gap-3 text-white relative overflow-hidden transition-colors duration-300" style={{ backgroundColor: themeColor }}>
                  <div className="absolute inset-0 bg-black/5" />
                  <div className="w-10 h-10 rounded-full bg-white/20 backdrop-blur-sm flex items-center justify-center shadow-inner relative z-10">
                    <Bot size={20} className="text-white drop-shadow-sm" />
                  </div>
                  <div className="relative z-10">
                    <h4 className="font-semibold text-[15px] leading-tight tracking-tight drop-shadow-sm">{botName || "Unnamed Bot"}</h4>
                    <div className="flex items-center gap-1.5 mt-0.5">
                      <span className="w-1.5 h-1.5 rounded-full bg-green-400 shadow-[0_0_8px_rgba(74,222,128,0.8)]" />
                      <p className="text-[11px] font-medium text-white/90 uppercase tracking-wider">Online</p>
                    </div>
                  </div>
                </div>

                {/* Provider badge */}
                <div className="px-4 py-2 bg-gray-50 dark:bg-gray-800/50 border-b border-gray-100 dark:border-gray-800 flex items-center gap-2">
                  <Zap size={12} className="text-brand-500" />
                  <span className="text-xs text-gray-500 dark:text-gray-400">
                    {PROVIDER_MODELS[aiProvider].label} — {modelName}
                  </span>
                </div>

                {/* Chat Body */}
                <div className="p-5 flex-1 h-[240px] bg-[#f8fafc] dark:bg-[#0f172a] flex flex-col gap-4 overflow-y-auto">
                  <div className="flex gap-2.5 max-w-[90%]">
                    <div className="w-7 h-7 rounded-full flex items-center justify-center flex-shrink-0 text-white shadow-sm mt-1 transition-colors duration-300" style={{ backgroundColor: themeColor }}>
                      <Bot size={14} />
                    </div>
                    <div className="bg-white dark:bg-gray-800 border border-gray-100 dark:border-gray-700/60 p-3.5 rounded-2xl rounded-tl-sm text-[14px] leading-relaxed text-gray-700 dark:text-gray-200 shadow-sm">
                      {welcomeMessage || "..."}
                    </div>
                  </div>
                </div>

                {/* Input Area */}
                <div className="p-3.5 bg-white dark:bg-gray-900 border-t border-gray-100 dark:border-gray-800 flex items-center gap-3">
                  <div className="flex-1 bg-gray-50 dark:bg-gray-800/50 rounded-full px-4 py-2.5 flex items-center border border-gray-100 dark:border-gray-800">
                    <input type="text" placeholder="Type a message..." className="w-full text-[14px] bg-transparent border-none focus:outline-none dark:text-white placeholder:text-gray-400 dark:placeholder:text-gray-500 cursor-not-allowed" disabled />
                  </div>
                  <button className="w-10 h-10 rounded-full text-white flex items-center justify-center transition-all duration-300 shadow-md flex-shrink-0" style={{ backgroundColor: themeColor }} disabled>
                    <Send size={16} className="ml-0.5" />
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
