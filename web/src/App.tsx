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
// import EventPage from "./pages/EventPage";
import CenteredLayout from "./CenteredLayout";
import RevelationPage from "./pages/RevelationPage";

function App() {
  return (
    <BrowserRouter>
      <Routes>
        <Route
          path="/"
          element={
            <CenteredLayout>
              <HomePage />
            </CenteredLayout>
          }
        />
        <Route
          path="/invite"
          element={
            <CenteredLayout>
              <InvitePage />
            </CenteredLayout>
          }
        />
        <Route
          path="/login"
          element={
            <CenteredLayout>
              <LoginPage />
            </CenteredLayout>
          }
        />
        <Route
          path="/quiz/:quizId"
          element={
            <CenteredLayout>
              <QuizPage />
            </CenteredLayout>
          }
        />

        <Route path="/quiz/:quizId/reveal" element={<RevelationPage />} />

        <Route
          path="/session/:sessionId"
          element={
            <CenteredLayout>
              <SessionPage />
            </CenteredLayout>
          }
        />

        <Route path="/dashboard/:quizId" element={<DashboardPage />} />
        {/* <Route path="/events" element={<EventPage />} /> */}

        <Route
          path="/unauthorized"
          element={
            <CenteredLayout>
              <UnauthorizedPage />
            </CenteredLayout>
          }
        />
        <Route
          path="*"
          element={
            <CenteredLayout>
              <NotFoundPage />
            </CenteredLayout>
          }
        />
      </Routes>
    </BrowserRouter>
  );
}

export default App;
