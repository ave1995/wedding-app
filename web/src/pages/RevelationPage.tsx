import { useEffect, useState } from "react";
import type { Quiz } from "../models/Quiz";
import { apiUrl } from "../functions/api";
import { get } from "../functions/fetch";
import { useApiErrorHandler } from "../hooks/useApiErrorHandler";
import { useParams } from "react-router-dom";
import type { Question } from "../models/Question";
import type { RevealResponse } from "../responses/RevealResponse";
import QuestionReveal from "../components/question/QuestionReveal";

function RevelationPage() {
  const { handleError } = useApiErrorHandler();
  const { quizId } = useParams();
  const [quiz, setQuiz] = useState<Quiz | null>(null);
  const [questionState, setQuestionState] = useState<{
    question: Question | null;
    nextIndex: number;
    totalQCount: number;
  }>({ question: null, totalQCount: 0, nextIndex: 0 });

  const [completed, setCompleted] = useState(false);

  const fetchResultOfQuestion = async function () {
    const result = await get<RevealResponse>(
      apiUrl(`/api/quiz/${quizId}/reveal`),
      { index: questionState.nextIndex },
      true
    );
    if (handleError(result.error, result.status)) return;

    if (!result.data?.goNext) {
      setCompleted(true);
      return;
    }

    setQuestionState({
      question: result.data!.question,
      nextIndex: result.data!.nextIndex,
      totalQCount: result.data.totalQCount, // TODO: tohle tady nemá být
    });
  };

  useEffect(() => {
    if (!quiz) {
      async function fetchQuiz() {
        const result = await get<Quiz>(apiUrl(`/api/quiz/${quizId}`), {}, true);
        if (handleError(result.error, result.status)) return;

        setQuiz(result.data);
      }
      fetchQuiz();
    }
    if (!questionState.question) {
      fetchResultOfQuestion();
    }
  }, [quizId]);

  if (completed) {
    return <div></div>;
  }

  return (
    <div className="w-96 h-screen">
      <QuestionReveal
        text={questionState.question?.Text ?? ""}
        currentQIndex={questionState.nextIndex}
        totalQCount={questionState.totalQCount}
        answers={questionState.question?.Answers ?? []}
        nextQuestion={fetchResultOfQuestion}
      ></QuestionReveal>
    </div>
  );
}

export default RevelationPage;
