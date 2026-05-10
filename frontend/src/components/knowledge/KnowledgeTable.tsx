import { useEffect, useState } from "react";
import {
  Table,
  TableBody,
  TableCell,
  TableHeader,
  TableRow,
} from "../ui/table";
import StatusBadge from "./StatusBadge";
import { knowledgeService, KnowledgeItem } from "../../services/knowledge";
import { Trash2 } from "lucide-react";
import Button from "../ui/button/Button";

interface KnowledgeTableProps {
  refreshSignal?: number; // increment this to trigger a refresh
}

export default function KnowledgeTable({ refreshSignal }: KnowledgeTableProps) {
  const [data, setData] = useState<KnowledgeItem[]>([]);
  const [isLoading, setIsLoading] = useState(true);

  const fetchData = async () => {
    setIsLoading(true);
    try {
      const result = await knowledgeService.getItems();
      setData(result);
    } catch (error) {
      console.error("Failed to fetch knowledge data", error);
    } finally {
      setIsLoading(false);
    }
  };

  useEffect(() => {
    fetchData();
  }, [refreshSignal]);

  const handleDelete = async (id: string) => {
    if (!confirm("Yakin hapus item ini?")) return;
    try {
      await knowledgeService.deleteItem(id);
      setData((prev) => prev.filter((d) => d.id !== id));
    } catch (err) {
      console.error("Failed to delete item:", err);
    }
  };

  const typeLabel = (type: string) => {
    const map: Record<string, string> = {
      file: "Dokumen",
      url: "URL",
      manual: "Manual Q&A",
    };
    return map[type] || type;
  };

  return (
    <div className="overflow-hidden rounded-2xl border border-gray-200 bg-white px-4 pb-3 pt-4 dark:border-gray-800 dark:bg-white/[0.03] sm:px-6">
      <div className="flex flex-col gap-2 mb-4 sm:flex-row sm:items-center sm:justify-between">
        <h3 className="text-lg font-semibold text-gray-800 dark:text-white/90">
          Trained Data
        </h3>
        <span className="text-sm text-gray-400">{data.length} item</span>
      </div>
      <div className="max-w-full overflow-x-auto">
        <Table>
          <TableHeader className="border-gray-100 dark:border-gray-800 border-y">
            <TableRow>
              <TableCell isHeader className="py-3 font-medium text-gray-500 text-start text-theme-xs dark:text-gray-400">
                Konten
              </TableCell>
              <TableCell isHeader className="py-3 font-medium text-gray-500 text-start text-theme-xs dark:text-gray-400">
                Tipe
              </TableCell>
              <TableCell isHeader className="py-3 font-medium text-gray-500 text-start text-theme-xs dark:text-gray-400">
                Status
              </TableCell>
              <TableCell isHeader className="py-3 font-medium text-gray-500 text-start text-theme-xs dark:text-gray-400">
                Tanggal
              </TableCell>
              <TableCell isHeader className="py-3 font-medium text-gray-500 text-start text-theme-xs dark:text-gray-400">
                Aksi
              </TableCell>
            </TableRow>
          </TableHeader>

          <TableBody className="divide-y divide-gray-100 dark:divide-gray-800">
            {isLoading ? (
              <TableRow>
                <TableCell colSpan={5} className="py-6 text-center text-gray-500 text-theme-sm">
                  Memuat data...
                </TableCell>
              </TableRow>
            ) : data.length === 0 ? (
              <TableRow>
                <TableCell colSpan={5} className="py-6 text-center text-gray-500 text-theme-sm">
                  Belum ada data. Upload dokumen atau tambahkan URL untuk mulai melatih bot.
                </TableCell>
              </TableRow>
            ) : (
              data.map((item) => (
                <TableRow key={item.id}>
                  <TableCell className="py-3 text-gray-800 text-theme-sm dark:text-white/90 max-w-[200px] truncate">
                    {item.content || "—"}
                  </TableCell>
                  <TableCell className="py-3 text-gray-500 text-theme-sm dark:text-gray-400">
                    {typeLabel(item.type)}
                  </TableCell>
                  <TableCell className="py-3 text-gray-500 text-theme-sm dark:text-gray-400">
                    <StatusBadge status={item.status} />
                  </TableCell>
                  <TableCell className="py-3 text-gray-500 text-theme-sm dark:text-gray-400">
                    {item.created_at ? new Date(item.created_at).toLocaleDateString("id-ID") : "—"}
                  </TableCell>
                  <TableCell className="py-3">
                    <Button
                      size="sm"
                      variant="outline"
                      onClick={() => handleDelete(item.id)}
                      className="text-red-500 hover:text-red-600 border-red-200 hover:border-red-300 px-2 py-1"
                    >
                      <Trash2 size={14} />
                    </Button>
                  </TableCell>
                </TableRow>
              ))
            )}
          </TableBody>
        </Table>
      </div>
    </div>
  );
}
