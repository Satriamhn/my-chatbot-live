import api from "./api";

// ─── Types ──────────────────────────────────────────────────────────────────

export type AIProvider = "gemini" | "openai";

export const PROVIDER_MODELS: Record<AIProvider, { label: string; models: { value: string; label: string }[] }> = {
  gemini: {
    label: "Google Gemini",
    models: [
      { value: "gemini-2.0-flash", label: "Gemini 2.0 Flash (Cepat & Murah)" },
      { value: "gemini-2.0-pro", label: "Gemini 2.0 Pro (Lebih Pintar)" },
      { value: "gemini-1.5-flash", label: "Gemini 1.5 Flash" },
    ],
  },
  openai: {
    label: "OpenAI",
    models: [
      { value: "gpt-4o-mini", label: "GPT-4o Mini (Cepat & Murah)" },
      { value: "gpt-4o", label: "GPT-4o (Terbaik)" },
      { value: "gpt-3.5-turbo", label: "GPT-3.5 Turbo" },
    ],
  },
};

export interface BotSettings {
  bot_name: string;
  welcome_message: string;
  system_prompt: string;
  ai_provider: AIProvider;
  model_name: string;
  has_byok_key: boolean;
  daily_message_limit: number;
  daily_message_count: number;
}

export interface UpdateBotSettingsPayload {
  bot_name: string;
  welcome_message: string;
  system_prompt: string;
  ai_provider?: AIProvider;
  model_name?: string;
  api_key?: string; // BYOK — optional
}

// ─── Service ─────────────────────────────────────────────────────────────────

export const settingsService = {
  getSettings: async (): Promise<BotSettings> => {
    const res = await api.get("/api/v1/settings/bot");
    return res.data;
  },

  updateSettings: async (payload: UpdateBotSettingsPayload): Promise<BotSettings> => {
    const res = await api.put("/api/v1/settings/bot", payload);
    return res.data;
  },
};
