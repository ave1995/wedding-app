import { useEffect, useState } from "react";
import { getText } from "../functions/fetch";
import Button, { ButtonTypeEnum } from "../components/Button";
import Toast from "../components/toast/Toast";

function HomePage() {
  const [message, setMessage] = useState<string | null>(null);

  const API_BASE_URL = import.meta.env.VITE_API_URL;

  useEffect(() => {
    if (!API_BASE_URL) {
      alert(`NO API BASE URL! ${API_BASE_URL}`);
      return;
    }
  }, [API_BASE_URL]);

  const sleep = (ms: number) =>
    new Promise((resolve) => setTimeout(resolve, ms));

  const handleClick = async () => {
    await sleep(1000);
    const result = await getText<string>(`${API_BASE_URL}/api/ping`);
    setMessage(result);
  };

  return (
    <div>
      {message && (
        <Toast message={message} onClose={() => setMessage(null)}></Toast>
      )}
      <Button
        onClick={handleClick}
        label="Ping"
        type={ButtonTypeEnum.Basic}
      ></Button>
    </div>
  );
}

export default HomePage;