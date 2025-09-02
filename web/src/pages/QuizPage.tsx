import { useEffect, useState } from "react";
import { useNavigate, useParams } from "react-router-dom";
import type { Quiz } from "../models/Quiz";
import { get, post } from "../functions/fetch";
import { apiUrl } from "../functions/api";
import type { StartSessionResponse } from "../responses/StartSessionResponse";
import Button, { ButtonTypeEnum } from "../components/Button";
import { useApiErrorHandler } from "../hooks/useApiErrorHandler";

function QuizPage() {
  const { handleError } = useApiErrorHandler();
  const navigate = useNavigate();
  const { quizId } = useParams();
  const [quiz, setQuiz] = useState<Quiz | null>(null);

  useEffect(() => {
    if (!quiz) {
      async function fetchQuiz() {
        const result = await get<Quiz>(apiUrl(`/api/quiz/${quizId}`), {}, true);
        if (handleError(result.error, result.status)) return;


        setQuiz(result.data);
      }
      fetchQuiz();
    }
  }, [quiz, quizId]);

  const startSession = async () => {
    const result = await post<StartSessionResponse>(
      apiUrl(`/api/quizzes/${quizId}/sessions`),
      null,
      null,
      true
    );
    if (handleError(result.error, result.status)) return;

    navigate(`/session/${result.data?.session_id}`);
  };

  // TODO: create loader
  if (!quiz) return <p>Loading quiz...</p>;

  return (
    <div className="flex flex-col w-96 h-screen items-center justify-center p-6 gap-4">
      <h1 className="text-xl font-bold">{quiz.Name}</h1>
      <p className="text-left p-3 font-medium">
        Milí <span className="font-semibold text-pink-500">hoste</span>, vítej na naší <span className="font-semibold text-pink-500">svatbě</span>! Jsme moc rádi, že jsi dnes tady s
        námi a sdílíš s námi tenhle výjimečný den. Tvoje přítomnost je pro nás
        <span className="font-semibold text-pink-500"> největší dar</span> – děkujeme, že slavíš, směješ se a vytváříš s námi
        nezapomenutelné chvíle. A protože se chceme pobavit společně, připravili
        jsme si pro tebe <span className="font-semibold text-pink-500">malý svatební kvíz</span>. Přejeme ti hodně štěstí a hlavně
        spoustu zábavy!
      </p>
      <div className="w-full">
        <Button
          onClick={startSession}
          label="Spustit kvíz"
          type={ButtonTypeEnum.Basic}
        ></Button>
      </div>
    </div>
  );
}

export default QuizPage;
