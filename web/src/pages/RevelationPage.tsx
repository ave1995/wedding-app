import { useEffect, useState } from "react";
import type { Quiz } from "../models/Quiz";
import { apiUrl } from "../functions/api";
import { get } from "../functions/fetch";
import { useApiErrorHandler } from "../hooks/useApiErrorHandler";
import { useParams } from "react-router-dom";
import type { Question } from "../models/Question";
import type { RevealResponse } from "../responses/RevealResponse";
import QuestionReveal from "../components/question/QuestionReveal";
import QuestionStats from "../components/question/QuestionStats";
import QuestionPhoto from "../components/question/QuestionPhoto";

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

  if (!questionState.question) {
    return <p>Načítám otázku...</p>;
  }

  return (
    <div className="flex flex-row h-screen w-screen items-center place-content-center">
      <div className="w-96 h-full">
        <QuestionReveal
          text={questionState.question.Text}
          type={questionState.question.Type}
          currentQIndex={questionState.nextIndex}
          totalQCount={questionState.totalQCount}
          answers={questionState.question.Answers}
          nextQuestion={fetchResultOfQuestion}
        ></QuestionReveal>
      </div>
      <div className="p-6 w-1/3 h-full">
        <QuestionPhoto path={questionState.question.PhotoPath} />
      </div>
      <div className="p-6">
        <QuestionStats id={questionState.question.ID} />
      </div>
    </div>
  );
}

export default RevelationPage;
