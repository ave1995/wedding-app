import { BrowserRouter, Route, Routes } from "react-router-dom";
import "./App.css";
import HomePage from "./pages/HomePage";
import InvitePage from "./pages/InvitePage";
import NotFoundPage from "./pages/NotFound";
import LoginPage from "./pages/LoginPage";
import QuizPage from "./pages/QuizPage";
import SessionPage from "./pages/SessionPage";
import UnauthorizedPage from "./pages/Unauthorized";
import DashboardPage from "./pages/DashboardPage";

function App() {
  return (
    <BrowserRouter>
      <Routes>
        <Route path="/" element={<HomePage />} />
        <Route path="/invite" element={<InvitePage />} />
        <Route path="/login" element={<LoginPage />} />
        <Route path="/quiz/:quizId" element={<QuizPage />} />
        <Route path="/session/:sessionId" element={<SessionPage />} />

        <Route path="/dashboard" element={<DashboardPage />} />

        <Route path="/unauthorized" element={<UnauthorizedPage />} />
        <Route path="*" element={<NotFoundPage />} />
      </Routes>
    </BrowserRouter>
  );
}

export default App;
