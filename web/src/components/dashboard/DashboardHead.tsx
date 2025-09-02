import type { DashboardStatus } from "./Dashboard";

interface DashboardHead {
  status: DashboardStatus;
}

export default function DashboardHead({ status }: DashboardHead) {
  const statusColors: Record<DashboardStatus, string> = {
    Online: "bg-green-500",
    Offline: "bg-gray-400",
    Error: "bg-red-500",
  };

  return (
    <div className="flex items-center place-content-between p-4 text-gray-800">
      <h1 className="text-2xl font-bold">Nástěnka Kvízu</h1>
      <div className="flex items-center gap-2">
        <span className={`h-3 w-3 rounded-full animate-pulse ${statusColors[status]}`}></span>
        <span className="text-sm font-medium">{status}</span>
      </div>
    </div>
  );
}
