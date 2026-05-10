interface StatusBadgeProps {
  status?: 'online' | 'offline' | 'busy' | 'away';
  className?: string;
}

export default function StatusBadge({ status = 'online', className = '' }: StatusBadgeProps) {
  let bgColor = 'bg-success-500';
  if (status === 'offline') bgColor = 'bg-gray-500';
  if (status === 'busy') bgColor = 'bg-error-500';
  if (status === 'away') bgColor = 'bg-warning-500';

  return (
    <span className={`absolute bottom-0 right-0 z-10 h-2.5 w-full max-w-2.5 rounded-full border-[1.5px] border-white ${bgColor} dark:border-gray-900 ${className}`}></span>
  );
}