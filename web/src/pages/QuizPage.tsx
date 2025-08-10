import { useEffect, useState } from "react";
import { useLocation, useParams } from "react-router-dom";
import type { Quiz } from "../models/Quiz";
import { get } from "../functions/fetch";
import { apiUrl } from "../functions/api";

function QuizPage() {
  const { state } = useLocation();
  const { quizId } = useParams();
  const [quiz, setQuiz] = useState<Quiz | null>(state?.quiz || null);

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

  // TODO: create loader
  if (!quiz) return <p>Loading quiz...</p>;

  return <div className="border-2">{quiz?.Name}</div>;
}

export default QuizPage;
