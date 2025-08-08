import { useState } from "react";
import InputText from "../components/input_text/InputText";
import Button from "../components/Button";

function LoginPage() {
  // Input Text state
  const [nameValue, setNameValue] = useState("");

  const [passwordValue, setPasswordValue] = useState("");

  async function handleJoin() {}

  return (
    <div>
      <h2>Login</h2>
      <InputText
        placeholder="Enter your name"
        value={nameValue}
        onChange={(e) => setNameValue(e.target.value)}
      />
      <InputText
        placeholder="Enter your password"
        value={passwordValue}
        onChange={(e) => setPasswordValue(e.target.value)}
      />
      <Button label="Login" onClick={handleJoin}></Button>
    </div>
  );
}

export default LoginPage;
