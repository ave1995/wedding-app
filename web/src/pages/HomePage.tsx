import { useState } from "react";
import { getText } from "../functions/fetch";
import Button, { ButtonTypeEnum } from "../components/Button";
import Toast from "../components/toast/Toast";
import { apiUrl } from "../functions/api";

function HomePage() {
  const [message, setMessage] = useState<string | null>(null);

  const sleep = (ms: number) =>
    new Promise((resolve) => setTimeout(resolve, ms));

  const handleClick = async () => {
    await sleep(1000);
    const result = await getText<string>(apiUrl("/ping"));
    if (result.error) {
      console.error(result.error);
    } else {
      setMessage(result.data);
    }
  };

  return (
    <div>
      {message && (
        <Toast message={message} onClose={() => setMessage(null)}></Toast>
      )}
      <Button onClick={handleClick} label="Ping" type={ButtonTypeEnum.Basic} />
    </div>
  );
}

export default HomePage;
