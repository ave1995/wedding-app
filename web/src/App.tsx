import { useState } from "react";
import "./App.css";
import Button, { ButtonTypeEnum } from "./components/Button";
import { getText } from "./functions/fetch";
import Toast from "./components/toast/Toast";

function App() {
  const [message, setMessage] = useState<string | null>(null);

  const handleClick = async () => {
    const result = await getText<string>("/api/ping");
    setMessage(result);
  };

  return (
    <div>
      {message && (
        <Toast message={message} onClose={() => setMessage(null)}></Toast>
      )}
      <Button
        onClickAsync={handleClick}
        label="Ping"
        type={ButtonTypeEnum.Basic}
      ></Button>
    </div>
  );
}

export default App;
