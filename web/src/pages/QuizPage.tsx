import { useEffect, useState } from "react";
import { useNavigate, useParams } from "react-router-dom";
import type { Quiz } from "../models/Quiz";
import { get, post } from "../functions/fetch";
import { apiUrl } from "../functions/api";
import type { StartSessionResponse } from "../responses/StartSessionResponse";
import Button, { ButtonTypeEnum } from "../components/Button";

function QuizPage() {
  const navigate = useNavigate();
  const { quizId } = useParams();
  const [quiz, setQuiz] = useState<Quiz | null>(null);

  useEffect(() => {
    if (!quiz) {
      async function fetchQuiz() {
        const result = await get<Quiz>(apiUrl(`/api/quiz/${quizId}`), {}, true);
        if (result.error) {
          console.error(result.error);
        } else {
          setQuiz(result.data);
        }
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
    if (result.error) {
      console.error(result.error);
    } else {
      navigate(`/session/${result.data?.session_id}`);
    }
  };

  // TODO: create loader
  if (!quiz) return <p>Loading quiz...</p>;

  return (
    <div className="p-6">
      <h1 className="text-xl font-bold">Quiz {quiz.Name}</h1>
      <Button
        onClick={startSession}
        label="Start Quiz"
        type={ButtonTypeEnum.Basic}
      ></Button>
    </div>
  );
}

export default QuizPage;
