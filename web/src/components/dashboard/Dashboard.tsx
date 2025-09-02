import { useEffect, useState } from "react";
import { apiUrl } from "../../functions/api";
import DashboardHead from "./DashboardHead";

export type Status = "Online" | "Offline" | "Error";

export type AnswerInfo = {
  QuestionText: string;
  UserIconUrl: string;
  Username: string;
};

export default function Dashboard() {
  const [status, setStatus] = useState<Status>("Offline");
  const [activeSessions, setActiveSessions] = useState(0);
  const [answers, setAnswer] = useState<AnswerInfo[]>([]);

  useEffect(() => {
    // Connect to WebSocket with desired topics
    const httpUrl = apiUrl(
      "/ws?topics=answer_submit,session_start,session_end"
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
          const answerInfo: AnswerInfo = {
            QuestionText: eventData.data.QuestionText,
            UserIconUrl: eventData.data.UserIconUrl,
            Username: eventData.data.Username,
          };
          setAnswer((prev) => [...prev, answerInfo]);
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
  }, []);

  return (
    <div className="flex flex-col gap-8 p-6 w-full h-full">
      <DashboardHead activeSessions={activeSessions} status={status} />
      <div className="flex flex-col gap-2 border rounded-xl border-b-4 p-4 bg-white/60 border-gray-400 text-gray-800">
        <h2 className="text-xl font-semibold">Odpovědi</h2>
        <div className="grid grid-cols-[30%_70%] gap-4 font-semibold border-b border-pink-500 pb-2 mb-2">
          <div>Uživatel</div>
          <div>Otázka</div>
        </div>
        <ul className="space-y-2">
          {answers.map((msg, idx) => (
            <li
              key={idx}
              className="grid grid-cols-[40%_60%] gap-4 items-center place-items-start border rounded-xl border-b-4 p-1 bg-white/60 border-[#3D52D5] text-gray-800"
            >
              <div className="flex gap-2 items-center place-content-center">
                <img
                  className="w-11 h-11"
                  src={msg.UserIconUrl}
                  alt={msg.Username}
                />
                {msg.Username}
              </div>
              <div>{msg.QuestionText}</div>
            </li>
          ))}
        </ul>
      </div>
    </div>
  );
}
