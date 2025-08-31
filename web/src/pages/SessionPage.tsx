import { useParams } from "react-router-dom";
import { useEffect, useState } from "react";
import type { Question } from "../models/Question";
import { get, post } from "../functions/fetch";
import { apiUrl } from "../functions/api";
import type { Result } from "../models/Result";
import type {
  SubmitAnswerResponse,
  SubmitAnswerResponseCompleted,
} from "../responses/SubmitAnswerResponse";
import type { QuestionResponse } from "../responses/QuestionResponse";
import { useApiErrorHandler } from "../hooks/useApiErrorHandler";
import QuestionForm from "../components/question/QuestionForm";

export default function SessionPage() {
  const { sessionId } = useParams();
  const { handleError } = useApiErrorHandler();

  const [question, setQuestion] = useState<Question | null>(null);
  const [completed, setCompleted] = useState(false);
  const [result, setResult] = useState<Result | null>(null);

  useEffect(() => {
    if (!question) {
      async function fetchQuestion() {
        const qResponse = await get<QuestionResponse>(
          apiUrl(`/api/sessions/${sessionId}/question`),
          null,
          true
        );
        if (handleError(qResponse.error, qResponse.status)) return;

        if (qResponse.data!.completed === true) {
          setCompleted(qResponse.data!.completed ?? false);
          const result = await get<SubmitAnswerResponseCompleted>(
            apiUrl(`/api/sessions/${sessionId}/result`),
            null,
            true
          );
          if (handleError(result.error, result.status)) return;

          console.log(result.data!.result);
          setResult(result.data!.result);
        } else {
          setQuestion(qResponse.data!.question ?? null);
        }
      }
      fetchQuestion();
    }
  }, [question, sessionId]);

  const submitAnswer = async (answerId: string) => {
    const result = await post<SubmitAnswerResponse>(
      apiUrl(`/api/sessions/${sessionId}/answers`),
      null,
      {
        question_id: question?.ID,
        answer_ids: [answerId],
      },
      true
    );
    if (handleError(result.error, result.status)) return;

    if (result.data!.completed === true) {
      setCompleted(true);
      setResult(result.data!.result);
    } else {
      setQuestion(result.data!.nextQuestion);
    }
  };

  if (completed) {
    return (
      <div className="p-6">
        <h1 className="text-xl font-bold">Hotovo!</h1>
        <p>
          Skóre: {result?.score}/{result?.total} ({result?.percentage}%)
        </p>
      </div>
    );
  }

  if (!question) {
    return <p>Načítám otázku...</p>;
  }

  return (
    <div className="w-96 h-screen p-4">
      <QuestionForm
        text={question.Text}
        currentQIndex={2}
        totalQCount={40}
        answers={question.Answers}
      />
    </div>
    // <div className="p-6">
    //   <h2 className="text-lg font-semibold">{question.Text}</h2>
    //   <div className="mt-4 flex flex-col gap-2">
    //     {question.Answers.map((a) => (
    //       <button
    //         key={a.ID}
    //         onClick={() => submitAnswer(a.ID)}
    //         className="px-4 py-2 bg-gray-200 rounded-lg hover:bg-gray-300"
    //       >
    //         {a.Text}
    //       </button>
    //     ))}
    //   </div>
    // </div>
  );
}
