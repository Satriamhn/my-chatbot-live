import api from "./api";

export interface DashboardStats {
  total_messages_today: number;
  active_sessions: number;
  total_users: number;
  total_knowledge: number;
}

export const statsService = {
  getStats: async (): Promise<DashboardStats> => {
    const response = await api.get<DashboardStats>("/api/v1/stats");
    return response.data;
  },
};
