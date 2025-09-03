import { useEffect, useState } from "react";
import { apiUrl } from "../../functions/api";
import DashboardHead from "./DashboardHead";
import { useQuestions } from "../../hooks/useQuestionEventHandler";
import DashboardQuestions from "./DashboardQuestions";
import DashboardActive from "./DashboardActive";
import { useAuthCheck } from "../../hooks/useAuthCheck";

export type DashboardStatus = "Online" | "Offline" | "Error";

export default function Dashboard() {
  const [status, setStatus] = useState<DashboardStatus>("Offline");
  const [activeSessions, setActiveSessions] = useState(0);
  const { questions, upsertQuestion } = useQuestions();

  const fetchAuthCheck = useAuthCheck();

  useEffect(() => {
    const init = async () => {
      const authorized = await fetchAuthCheck();
      if (!authorized) return;

      const httpUrl = apiUrl(
        "/api/ws?topics=answer_submit,session_start,session_end,question_open"
      );
      const wsUrl = httpUrl.replace(/^http/, "ws");

      const socket = new WebSocket(wsUrl);

      socket.onopen = () => setStatus("Online");

      socket.onmessage = (event: MessageEvent) => {
        try {
          const eventData = JSON.parse(event.data);
          if (eventData.topic === "session_start") {
            setActiveSessions((prev) => prev + 1);
          } else if (eventData.topic === "session_end") {
            setActiveSessions((prev) => Math.max(prev - 1, 0));
          } else if (eventData.topic === "answer_submit") {
            upsertQuestion({
              UserID: eventData.data.UserID,
              QuestionID: eventData.data.QuestionID,
              QuestionText: eventData.data.QuestionText,
              UserIconUrl: eventData.data.UserIconUrl,
              Username: eventData.data.Username,
              Status: "Odpověděl",
            });
          } else if (eventData.topic === "question_open") {
            upsertQuestion({
              UserID: eventData.data.UserID,
              QuestionID: eventData.data.QuestionID,
              QuestionText: eventData.data.QuestionText,
              UserIconUrl: eventData.data.UserIconUrl,
              Username: eventData.data.Username,
              Status: "Přemýšlí",
            });
          }
        } catch (e) {
          console.error("Chyba při parsování WebSocket zprávy:", e);
        }
      };

      socket.onclose = (event) => {
        console.log("Closed:", event.code, event.reason);
        setStatus("Offline");
      };

      socket.onerror = (err) => {
        console.error("WebSocket error", err);
        setStatus("Error");
      };

      // Cleanup
      return () => {
        socket.close();
      };
    };

    init();
  }, [upsertQuestion]);

  return (
    <div className="flex flex-col gap-8 p-6 w-full h-full">
      <DashboardHead status={status} />
      <DashboardActive activeSessions={activeSessions}></DashboardActive>
      <DashboardQuestions questions={questions}></DashboardQuestions>
    </div>
  );
}
