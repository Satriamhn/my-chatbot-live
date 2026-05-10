import api from "./api";

// ─── Types ──────────────────────────────────────────────────────────────────

export interface KnowledgeBase {
  id: string;
  name: string;
  description: string;
  created_at: string;
}

export interface KnowledgeItem {
  id: string;
  knowledge_base_id: string;
  type: "file" | "url" | "manual";
  status: "queued" | "processing" | "ready" | "failed";
  content: string;
  created_at: string;
}

// ─── Service ─────────────────────────────────────────────────────────────────

export const knowledgeService = {
  getItems: async (knowledgeBaseId?: string): Promise<KnowledgeItem[]> => {
    const params = knowledgeBaseId ? { knowledge_base_id: knowledgeBaseId } : {};
    const res = await api.get("/api/v1/knowledge-base/items", { params });
    return res.data;
  },

  createItem: async (payload: {
    knowledge_base_id: string;
    type: "file" | "url" | "manual";
    content: string;
    metadata?: string;
  }): Promise<KnowledgeItem> => {
    const res = await api.post("/api/v1/knowledge-base/items", payload);
    return res.data;
  },

  updateStatus: async (
    id: string,
    status: "queued" | "processing" | "ready" | "failed"
  ): Promise<KnowledgeItem> => {
    const res = await api.patch(`/api/v1/knowledge-base/items/${id}/status`, { status });
    return res.data;
  },

  deleteItem: async (id: string): Promise<void> => {
    await api.delete(`/api/v1/knowledge-base/items/${id}`);
  },
};
