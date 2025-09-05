import { useEffect, useState } from "react";
import { apiUrl } from "../../functions/api";
import DashboardHead from "./DashboardHead";
import { useQuestions } from "../../hooks/useQuestionEventHandler";
import DashboardQuestions from "./DashboardQuestions";
import DashboardActive from "./DashboardActive";
import { useAuthCheck } from "../../hooks/useAuthCheck";
import { get } from "../../functions/fetch";
import { useApiErrorHandler } from "../../hooks/useApiErrorHandler";

export type DashboardStatus = "Online" | "Offline" | "Error";

type ActiveSessions = {
  count: number;
};

interface Dashboard {
  quizID: string;
}

export default function Dashboard({ quizID }: Dashboard) {
  const { handleError } = useApiErrorHandler();
  const [status, setStatus] = useState<DashboardStatus>("Offline");
  const [activeSessions, setActiveSessions] = useState(0);
  const { questions, upsertQuestion } = useQuestions();

  const fetchAuthCheck = useAuthCheck();

  useEffect(() => {
    const active = async () => {
      const result = await get<ActiveSessions>(
        apiUrl(`/api/quiz/${quizID}/active`),
        null,
        true
      );
      if (handleError(result.error, result.status)) return;

      setActiveSessions(result.data?.count!);
    };
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
          } else if (eventData.topic === "heartbeat") {
            console.log("Received:", eventData.data);

            const pong = {
              topic: "heartbeat",
              data: "pong!",
            };
            // send pong back to server
            const pongMessage = JSON.stringify(pong);

            if (socket.readyState === WebSocket.OPEN) {
              socket.send(pongMessage);
              console.log("Sent:", pong.data);
            }
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

    active();
    init();
  }, [upsertQuestion]);

  return (
    <div className="flex flex-col p-6 w-1/2 h-screen">
      <div className="flex-none pb-8">
        <DashboardHead status={status} />
      </div>
      <div className="flex-none pb-8">
        <DashboardActive activeSessions={activeSessions}></DashboardActive>
      </div>
      <h2 className="text-lg font-semibold text-left mx-6 mb-2 relative inline-block">
        <span className="relative z-10 pr-4 bg-[#f6f4f3]">Aktivita</span>
        <span className="absolute left-0 top-1/2 w-full h-px bg-gray-300 -z-10"></span>
      </h2>
      <div className="flex-grow overflow-auto">
        <DashboardQuestions questions={questions}></DashboardQuestions>
      </div>
    </div>
  );
}
