import { useEffect } from "react";
import { useSearchParams } from "react-router-dom";
import { get } from "../functions/fetch";
import InputText from "../components/input_text/InputText";

function InvitePage() {
  const API_BASE_URL = import.meta.env.VITE_API_URL;

  const [searchParams] = useSearchParams();
  const code = searchParams.get("code");

  useEffect(() => {
    if (code) {
      async function fetchQuiz() {
        const result = await get<any>(`${API_BASE_URL}/api/join-quiz`, {
          invite: code,
        });
        console.log(result);
        // You can set state here if needed
      }
      fetchQuiz();
    }
  }, [code]);

  return (
    <div>
      <h2>New Guest</h2>
      <InputText></InputText>
    </div>
  );
}
export default InvitePage;
