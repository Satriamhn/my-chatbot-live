import Badge from "../ui/badge/Badge";

interface StatusBadgeProps {
  status: "queued" | "processing" | "ready" | "failed";
}

export default function StatusBadge({ status }: StatusBadgeProps) {
  let color: "success" | "warning" | "error" | "info" = "info";
  let label: string = status;

  switch (status) {
    case "ready":
      color = "success";
      label = "Ready";
      break;
    case "processing":
      color = "warning";
      label = "Processing";
      break;
    case "failed":
      color = "error";
      label = "Failed";
      break;
    case "queued":
      color = "info"; // Will use info/gray style depending on Badge impl
      label = "Queued";
      break;
  }

  return (
    <Badge size="sm" color={color}>
      {label}
    </Badge>
  );
}
