import { useEffect, useState } from "react";
import { apiUrl } from "../../functions/api";
import DashboardHead from "./DashboardHead";
import { useQuestions } from "../../hooks/useQuestionEventHandler";
import DashboardQuestions from "./DashboardQuestions";
import DashboardActive from "./DashboardActive";

export type DashboardStatus = "Online" | "Offline" | "Error";

export default function Dashboard() {
  const [status, setStatus] = useState<DashboardStatus>("Offline");
  const [activeSessions, setActiveSessions] = useState(0);
  const { questions, upsertQuestion } = useQuestions();

  useEffect(() => {
    // Connect to WebSocket with desired topics
    const httpUrl = apiUrl(
      "/ws?topics=answer_submit,session_start,session_end,question_open"
    );
    const wsUrl = httpUrl.replace(/^http/, "ws");

    const socket = new WebSocket(wsUrl);

    socket.onopen = () => {
      setStatus("Online");
    };

    socket.onmessage = (event: MessageEvent) => {
      console.log(event.data);
      try {
        const eventData = JSON.parse(event.data);

        if (eventData.topic === "session_start") {
          setActiveSessions((prev) => prev + 1);
        } else if (eventData.topic === "session_end") {
          setActiveSessions((prev) => Math.max(prev - 1, 0)); // nikdy méně než 0
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

    socket.onclose = () => {
      setStatus("Offline");
    };

    socket.onerror = () => {
      setStatus("Error");
    };

    return () => {
      socket.close();
    };
  }, [upsertQuestion]);

  return (
    <div className="flex flex-col gap-8 p-6 w-full h-full">
      <DashboardHead status={status} />
      <DashboardActive activeSessions={activeSessions}></DashboardActive>
      <DashboardQuestions questions={questions}></DashboardQuestions>
    </div>
  );
}
