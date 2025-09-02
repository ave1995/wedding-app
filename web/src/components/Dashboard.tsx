import { useEffect, useState } from "react";
import { apiUrl } from "../functions/api";

export default function Dashboard() {
  const [messages, setMessages] = useState<string[]>([]);
  const [status, setStatus] = useState<"Connected" | "Disconnected" | "Error">(
    "Disconnected"
  );

  useEffect(() => {
    // Connect to WebSocket with desired topics
    const httpUrl = apiUrl("/ws?topics=answer_submit, session_start, session_end");
    const wsUrl = httpUrl.replace(/^http/, "ws");

     const socket = new WebSocket(wsUrl);

    socket.onopen = () => {
      setStatus("Connected");
    };

    socket.onmessage = (event: MessageEvent) => {
      setMessages((prev) => [...prev, event.data]);
    };

    socket.onclose = () => {
      setStatus("Disconnected");
    };

    socket.onerror = () => {
      setStatus("Error");
    };

    return () => {
      socket.close();
    };
  }, []);

  return (
    <div className="p-6 min-h-screen bg-gray-100">
      <h1 className="text-2xl font-bold mb-4">Dashboard</h1>
      <p className="mb-2">
        Status: <span className="font-semibold">{status}</span>
      </p>

      <div className="bg-white shadow-md rounded-xl p-4 max-h-[400px] overflow-y-auto">
        <h2 className="text-lg font-semibold mb-2">Messages</h2>
        {messages.length === 0 ? (
          <p className="text-gray-500">No messages yet…</p>
        ) : (
          <ul className="space-y-2">
            {messages.map((msg, idx) => (
              <li
                key={idx}
                className="p-2 bg-gray-50 rounded border border-gray-200"
              >
                {msg}
              </li>
            ))}
          </ul>
        )}
      </div>
    </div>
  );
}
