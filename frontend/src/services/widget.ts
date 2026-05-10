import axios from "axios";

const widgetApi = axios.create({
	baseURL: import.meta.env.VITE_API_URL || "http://localhost:8080",
	headers: {
		"Content-Type": "application/json",
	},
});

export interface WidgetSettings {
	bot_name: string;
	welcome_message: string;
}

export const widgetService = {
	getSettings: async (orgId: string): Promise<WidgetSettings> => {
		const response = await widgetApi.get<WidgetSettings>("/api/v1/widget/settings", {
			params: { org_id: orgId },
		});
		return response.data;
	},
};
