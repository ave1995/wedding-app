import Dashboard from "../components/dashboard/Dashboard";
import Thanks from "../assets/Thanks.png";

function DashboardPage() {
  return (
    <div className="grid grid-cols-2 w-screen gap-20 p-5 h-full min-h-screen">
      <Dashboard />
      <div className="flex items-center place-content-center p-8">
        <img
          src={Thanks}
          alt="Thanks"
          className="max-w-full max-h-full object-contain"
        />
      </div>
    </div>
  );
}

export default DashboardPage;
