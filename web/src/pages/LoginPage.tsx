import { useState } from "react";
import InputText from "../components/input_text/InputText";
import Button from "../components/Button";
import { post } from "../functions/fetch";
import { apiUrl } from "../functions/api";
import type { SimpleResponse } from "../responses/SimpleResponse";
import { useApiErrorHandler } from "../hooks/useApiErrorHandler";

function LoginPage() {
  const { handleError } = useApiErrorHandler();
  // Input Text state
  const [emailValue, setEmailValue] = useState("");

  const [passwordValue, setPasswordValue] = useState("");

  async function handleJoin() {
    const result = await post<SimpleResponse>(
      apiUrl("/auth/login"),
      null,
      {
        email: emailValue,
        password: passwordValue,
      },
      true
    );
    if (handleError(result.error, result.status)) return;

    console.log(result.data);
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
