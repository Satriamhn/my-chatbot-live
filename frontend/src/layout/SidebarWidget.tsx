import { Link } from "react-router";

export default function SidebarWidget() {
  return (
    <div className="mx-auto mt-auto mb-6 w-full max-w-[250px] px-4">
      <div className="relative overflow-hidden rounded-2xl bg-gradient-to-br from-gray-50 to-gray-100 p-4 shadow-sm ring-1 ring-gray-200/50 dark:from-white/[0.05] dark:to-white/[0.02] dark:ring-white/10">
        {/* Decorative background element */}
        <div className="absolute -right-4 -top-4 size-16 rounded-full bg-brand-500/10 blur-xl dark:bg-brand-400/20" />
        
        <div className="relative flex items-center gap-3">
          <div className="relative flex size-10 shrink-0 items-center justify-center rounded-full bg-brand-100 text-brand-600 dark:bg-brand-500/20 dark:text-brand-300">
            <span className="text-sm font-bold">OP</span>
            <span className="absolute bottom-0 right-0 size-2.5 rounded-full border-2 border-white bg-green-500 dark:border-gray-900" />
          </div>
          
          <div className="flex flex-col">
            <h3 className="text-sm font-semibold text-gray-900 dark:text-white">
              System Owner
            </h3>
            <p className="text-xs font-medium text-gray-500 dark:text-gray-400">
              Admin & Support
            </p>
          </div>
        </div>

        <div className="relative mt-4">
          <Link
            to="/contact"
            className="flex w-full items-center justify-center rounded-lg bg-white px-3 py-2 text-xs font-semibold text-gray-700 shadow-sm ring-1 ring-inset ring-gray-300 transition-all hover:bg-gray-50 dark:bg-gray-800 dark:text-gray-200 dark:ring-gray-700 dark:hover:bg-gray-700/50 dark:hover:ring-gray-600"
          >
            Contact Owner
          </Link>
        </div>
      </div>
    </div>
  );
}
