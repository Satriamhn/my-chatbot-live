import PageBreadcrumb from "../../components/common/PageBreadCrumb";
import ComponentCard from "../../components/common/ComponentCard";

export default function Contact() {
  return (
    <>
      <PageBreadcrumb pageTitle="Hubungi Saya (Owner)" />

      <div className="flex flex-col gap-6">
        <ComponentCard title="Contact Details">
          <div className="flex flex-col gap-4">
            <p className="text-gray-600 dark:text-gray-400">
              Silakan hubungi saya untuk keperluan bisnis, feedback, atau laporan masalah.
            </p>
            <div className="bg-gray-50 dark:bg-gray-800 p-4 rounded-lg">
              <p className="font-medium">Email: owner@saas-chatbot.com</p>
              <p className="font-medium">Phone: +62 812 3456 7890</p>
            </div>
          </div>
        </ComponentCard>
      </div>
    </>
  );
}
