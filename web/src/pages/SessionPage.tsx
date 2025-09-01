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
import ResultForm from "../components/ResultForm";


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

  const submitAnswer = async (answerIds: string[]) => {
    const result = await post<SubmitAnswerResponse>(
      apiUrl(`/api/sessions/${sessionId}/answers`),
      null,
      {
        question_id: question?.ID,
        answer_ids: answerIds,
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
    return <ResultForm sessionId={sessionId!} resultFrom={result}/>
  }

  if (!question) {
    return <p>Načítám otázku...</p>;
  }

  return (
    <div className="w-96 h-screen">
      <QuestionForm
        text={question.Text}
        currentQIndex={2}
        totalQCount={40}
        answers={question.Answers}
        submitAnswer={submitAnswer}
      />
    </div>
  );
}
