import { useState } from "react";
import InputText from "../components/input_text/InputText";
import Button from "../components/Button";
import { post } from "../functions/fetch";
import { apiUrl } from "../functions/api";
import type { SimpleResponse } from "../models/SimpleResponse";

function LoginPage() {
  // Input Text state
  const [emailValue, setEmailValue] = useState("");

  const [passwordValue, setPasswordValue] = useState("");

  async function handleJoin() {
    const result = await post<SimpleResponse>(
      apiUrl("/auth/login"),
      {},
      {
        email: emailValue,
        password: passwordValue,
      },
      true
    );

    if (result.error) {
      console.error(result.error);
    } else {
      console.log(result.data);
    }
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
