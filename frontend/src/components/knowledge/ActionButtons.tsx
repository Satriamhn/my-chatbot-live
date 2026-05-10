import { useState } from "react";
import { Upload, Link2, MessageSquarePlus } from "lucide-react";
import Button from "../ui/button/Button";
import { knowledgeService } from "../../services/knowledge";

interface ActionButtonsProps {
  onAdded?: () => void; // callback to refresh table
}

export default function ActionButtons({ onAdded }: ActionButtonsProps) {
  const [urlInput, setUrlInput] = useState("");
  const [qaContent, setQaContent] = useState("");
  const [showUrlModal, setShowUrlModal] = useState(false);
  const [showQaModal, setShowQaModal] = useState(false);
  const [loading, setLoading] = useState(false);

  // Needs a default KB ID — in production this would come from selected KB
  const DEFAULT_KB_ID = "00000000-0000-0000-0000-000000000001";

  const handleUpload = async (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0];
    if (!file) return;
    setLoading(true);
    try {
      await knowledgeService.createItem({
        knowledge_base_id: DEFAULT_KB_ID,
        type: "file",
        content: file.name,
        metadata: JSON.stringify({ filename: file.name, size: file.size }),
      });
      onAdded?.();
    } catch (err) {
      console.error("Upload failed:", err);
    } finally {
      setLoading(false);
      e.target.value = "";
    }
  };

  const handleAddUrl = async () => {
    if (!urlInput.trim()) return;
    setLoading(true);
    try {
      await knowledgeService.createItem({
        knowledge_base_id: DEFAULT_KB_ID,
        type: "url",
        content: urlInput.trim(),
      });
      setUrlInput("");
      setShowUrlModal(false);
      onAdded?.();
    } catch (err) {
      console.error("Add URL failed:", err);
    } finally {
      setLoading(false);
    }
  };

  const handleAddQA = async () => {
    if (!qaContent.trim()) return;
    setLoading(true);
    try {
      await knowledgeService.createItem({
        knowledge_base_id: DEFAULT_KB_ID,
        type: "manual",
        content: qaContent.trim(),
      });
      setQaContent("");
      setShowQaModal(false);
      onAdded?.();
    } catch (err) {
      console.error("Add Q&A failed:", err);
    } finally {
      setLoading(false);
    }
  };

  return (
    <>
      <div className="flex flex-wrap items-center gap-3">
        {/* Upload Dokumen */}
        <label className="inline-flex items-center gap-2 rounded-lg border border-gray-300 bg-white px-4 py-2 text-theme-sm font-medium text-gray-700 shadow-theme-xs hover:bg-gray-50 hover:text-gray-800 dark:border-gray-700 dark:bg-gray-800 dark:text-gray-400 dark:hover:bg-white/[0.03] dark:hover:text-gray-200 cursor-pointer">
          <Upload size={15} />
          {loading ? "Mengupload..." : "Upload Dokumen"}
          <input type="file" className="hidden" onChange={handleUpload} disabled={loading} />
        </label>

        {/* URL Sync */}
        <button
          onClick={() => setShowUrlModal(true)}
          className="inline-flex items-center gap-2 rounded-lg border border-gray-300 bg-white px-4 py-2 text-theme-sm font-medium text-gray-700 shadow-theme-xs hover:bg-gray-50 hover:text-gray-800 dark:border-gray-700 dark:bg-gray-800 dark:text-gray-400 dark:hover:bg-white/[0.03] dark:hover:text-gray-200"
        >
          <Link2 size={15} />
          URL Sync
        </button>

        {/* Manual Q&A */}
        <button
          onClick={() => setShowQaModal(true)}
          className="inline-flex items-center gap-2 rounded-lg border border-gray-300 bg-white px-4 py-2 text-theme-sm font-medium text-gray-700 shadow-theme-xs hover:bg-gray-50 hover:text-gray-800 dark:border-gray-700 dark:bg-gray-800 dark:text-gray-400 dark:hover:bg-white/[0.03] dark:hover:text-gray-200"
        >
          <MessageSquarePlus size={15} />
          Manual Q&A
        </button>
      </div>

      {/* URL Modal */}
      {showUrlModal && (
        <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/40 backdrop-blur-sm">
          <div className="w-full max-w-md rounded-2xl border border-gray-200 bg-white p-6 shadow-xl dark:border-gray-700 dark:bg-gray-900">
            <h3 className="mb-4 text-lg font-semibold text-gray-800 dark:text-white">Tambah URL</h3>
            <input
              type="url"
              value={urlInput}
              onChange={(e) => setUrlInput(e.target.value)}
              placeholder="https://contoh.com/artikel"
              className="w-full rounded-lg border border-gray-300 px-4 py-2.5 text-sm text-gray-800 dark:border-gray-700 dark:bg-gray-800 dark:text-white focus:outline-none focus:ring-2 focus:ring-brand-500/20 mb-4"
            />
            <div className="flex gap-3 justify-end">
              <Button variant="outline" onClick={() => setShowUrlModal(false)}>Batal</Button>
              <Button variant="primary" onClick={handleAddUrl} disabled={loading || !urlInput.trim()}>
                {loading ? "Menambahkan..." : "Tambah"}
              </Button>
            </div>
          </div>
        </div>
      )}

      {/* Q&A Modal */}
      {showQaModal && (
        <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/40 backdrop-blur-sm">
          <div className="w-full max-w-md rounded-2xl border border-gray-200 bg-white p-6 shadow-xl dark:border-gray-700 dark:bg-gray-900">
            <h3 className="mb-4 text-lg font-semibold text-gray-800 dark:text-white">Tambah Q&A Manual</h3>
            <textarea
              value={qaContent}
              onChange={(e) => setQaContent(e.target.value)}
              placeholder="Q: Apa jam operasional toko?&#10;A: Kami buka Senin-Sabtu pukul 09.00-21.00."
              rows={5}
              className="w-full rounded-lg border border-gray-300 px-4 py-2.5 text-sm text-gray-800 dark:border-gray-700 dark:bg-gray-800 dark:text-white focus:outline-none focus:ring-2 focus:ring-brand-500/20 mb-4 resize-none"
            />
            <div className="flex gap-3 justify-end">
              <Button variant="outline" onClick={() => setShowQaModal(false)}>Batal</Button>
              <Button variant="primary" onClick={handleAddQA} disabled={loading || !qaContent.trim()}>
                {loading ? "Menyimpan..." : "Simpan"}
              </Button>
            </div>
          </div>
        </div>
      )}
    </>
  );
}
