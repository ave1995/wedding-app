interface DashboardActive {
  activeSessions: number;
}

export default function DashboardActive({ activeSessions }: DashboardActive) {
  return (
    <div className="flex flex-col">
      <p className="text-7xl font-bold text-pink-500">{activeSessions}</p>
      <p className="font-medium text-gray-800">Aktivních hostů</p>
    </div>
  );
}
