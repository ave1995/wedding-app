import { useParams } from "react-router-dom";
import Dashboard from "../components/dashboard/Dashboard";

function DashboardPage() {
  const { quizId } = useParams();

  if (!quizId) {
    return null;
  }

  return (
    <div className="flex flex-col w-screen h-screen items-center place-content-start">
      <Dashboard quizID={quizId} />
    </div>
  );
}

export default DashboardPage;
