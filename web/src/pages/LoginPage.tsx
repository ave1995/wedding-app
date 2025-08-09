import { useState } from "react";
import InputText from "../components/input_text/InputText";
import Button from "../components/Button";
import { post } from "../functions/fetch";
import type { AccessToken } from "../models/AccessToken";
import { apiUrl } from "../functions/api";

function LoginPage() {
  // Input Text state
  const [emailValue, setEmailValue] = useState("");

  const [passwordValue, setPasswordValue] = useState("");

  async function handleJoin() {
    const result = await post<AccessToken>(
      apiUrl("/auth/login"),
      {},
      {
        email: emailValue,
        password: passwordValue,
      },
      true
    );
  }

  return (
    <div>
      <h2>Login</h2>
      <InputText
        placeholder="Enter your email"
        value={emailValue}
        onChange={(e) => setEmailValue(e.target.value)}
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
