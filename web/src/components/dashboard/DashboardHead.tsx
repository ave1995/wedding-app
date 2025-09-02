import type { Status } from "./Dashboard";

interface DashboardHead {
  status: Status;
  activeSessions: number;
}

export default function DashboardHead({
  status,
  activeSessions,
}: DashboardHead) {
  const statusColors: Record<Status, string> = {
    Online: "bg-green-500",
    Offline: "bg-gray-400",
    Error: "bg-red-500",
  };

  return (
    <div className="flex items-center place-content-between gap-2 border rounded-xl border-b-4 p-4 bg-white/60 border-gray-400 text-gray-800">
      <div className="flex gap-4">
        <h1 className="text-2xl font-bold">Nástěnka Kvízu</h1>
        <div className="flex items-center gap-2">
          <span
            className={`h-3 w-3 rounded-full ${statusColors[status]}`}
          ></span>
          <span className="text-sm font-medium">{status}</span>
        </div>
      </div>
      <div className="flex gap-12 text-xl font-semibold">
        <p>
          Aktivní hosti:{" "}
          <span className="font-semibold text-pink-500">{activeSessions}</span>
        </p>
      </div>
    </div>
  );
}
