import { useState, useEffect } from "react";
import PageBreadcrumb from "../../components/common/PageBreadCrumb";
import ComponentCard from "../../components/common/ComponentCard";
import Button from "../../components/ui/button/Button";
import { Table, TableBody, TableCell, TableHeader, TableRow } from "../../components/ui/table";
import Badge from "../../components/ui/badge/Badge";
import { Upload, Link2, Edit3, Trash2, BookOpen } from "lucide-react";

interface KnowledgeItem {
  id: string;
  name: string;
  sourceType: string;
  status: string;
  dateAdded: string;
}

export default function KnowledgeBase() {
  const [knowledge, setKnowledge] = useState<KnowledgeItem[]>([]);

  useEffect(() => {
    // Load from mock
    fetch("/src/mocks/knowledge.json")
      .then(res => res.json())
      .then(data => setKnowledge(data))
      .catch(err => console.error("Failed to load knowledge mock", err));
  }, []);

  return (
    <>
      <PageBreadcrumb pageTitle="Knowledge Base" />

      <div className="flex flex-col gap-6">
        <ComponentCard title="Manage Sources" desc="Upload files, sync URLs, or enter manual Q&A pairs for your bot.">
          
          <div className="flex gap-4 mb-6">
            <Button variant="primary" startIcon={<Upload size={16} />}>Upload File</Button>
            <Button variant="outline" startIcon={<Link2 size={16} />}>URL Sync</Button>
            <Button variant="outline" startIcon={<BookOpen size={16} />}>Manual Q&A</Button>
          </div>

          <div className="overflow-hidden rounded-xl border border-gray-200 bg-white dark:border-white/[0.05] dark:bg-white/[0.03]">
            <div className="max-w-full overflow-x-auto">
              <Table>
                <TableHeader className="border-b border-gray-200 dark:border-white/[0.05]">
                  <TableRow>
                    <TableCell isHeader className="px-5 py-3 font-medium text-gray-500 text-start text-theme-xs dark:text-gray-400">
                      Document Name
                    </TableCell>
                    <TableCell isHeader className="px-5 py-3 font-medium text-gray-500 text-start text-theme-xs dark:text-gray-400">
                      Source Type
                    </TableCell>
                    <TableCell isHeader className="px-5 py-3 font-medium text-gray-500 text-start text-theme-xs dark:text-gray-400">
                      Status
                    </TableCell>
                    <TableCell isHeader className="px-5 py-3 font-medium text-gray-500 text-start text-theme-xs dark:text-gray-400">
                      Date Added
                    </TableCell>
                    <TableCell isHeader className="px-5 py-3 font-medium text-gray-500 text-start text-theme-xs dark:text-gray-400">
                      Actions
                    </TableCell>
                  </TableRow>
                </TableHeader>

                <TableBody className="divide-y divide-gray-200 dark:divide-white/[0.05]">
                  {knowledge.map((item) => (
                    <TableRow key={item.id}>
                      <TableCell className="px-5 py-4 sm:px-6 text-start">
                        <span className="block font-medium text-gray-800 dark:text-white/90">
                          {item.name}
                        </span>
                      </TableCell>
                      <TableCell className="px-5 py-4 sm:px-6 text-start text-gray-500 dark:text-gray-400">
                        {item.sourceType}
                      </TableCell>
                      <TableCell className="px-5 py-4 sm:px-6 text-start">
                        <Badge 
                          size="sm" 
                          color={item.status === 'Active' ? 'success' : item.status === 'Processing' ? 'warning' : 'error'}
                        >
                          {item.status}
                        </Badge>
                      </TableCell>
                      <TableCell className="px-5 py-4 sm:px-6 text-start text-gray-500 dark:text-gray-400">
                        {item.dateAdded}
                      </TableCell>
                      <TableCell className="px-5 py-4 sm:px-6 text-start">
                        <div className="flex items-center gap-3">
                          <button className="text-gray-500 hover:text-brand-500 dark:text-gray-400 dark:hover:text-brand-400">
                            <Edit3 size={16} />
                          </button>
                          <button className="text-gray-500 hover:text-error-500 dark:text-gray-400 dark:hover:text-error-400">
                            <Trash2 size={16} />
                          </button>
                        </div>
                      </TableCell>
                    </TableRow>
                  ))}
                </TableBody>
              </Table>
            </div>
          </div>
        </ComponentCard>
      </div>
    </>
  );
}
