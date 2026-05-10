import { useState } from "react";
import KnowledgeTable from "../../components/knowledge/KnowledgeTable";
import ActionButtons from "../../components/knowledge/ActionButtons";
import PageMeta from "../../components/common/PageMeta";

export default function KnowledgePage() {
  const [refreshSignal, setRefreshSignal] = useState(0);

  return (
    <>
      <PageMeta
        title="Knowledge Base | Chatbot Dashboard"
        description="Manage your chatbot's knowledge base and trained data."
      />
      <div className="space-y-6">
        <div className="flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between">
          <div>
            <h2 className="text-2xl font-bold text-gray-800 dark:text-white/90">
              Knowledge Base
            </h2>
            <p className="mt-1 text-sm text-gray-500 dark:text-gray-400">
              Upload dokumen, sync URL, atau tambah Q&A manual untuk melatih chatbot.
            </p>
          </div>
          <ActionButtons onAdded={() => setRefreshSignal((s) => s + 1)} />
        </div>

        <KnowledgeTable refreshSignal={refreshSignal} />
      </div>
    </>
  );
}
