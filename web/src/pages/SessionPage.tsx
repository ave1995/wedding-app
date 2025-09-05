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

  const [completed, setCompleted] = useState(false);
  const [result, setResult] = useState<Result | null>(null);
  const [questionState, setQuestionState] = useState<{
    question: Question | null;
    currentQIndex: number;
    totalQCount: number;
  }>({ question: null, currentQIndex: 0, totalQCount: 0 });

  useEffect(() => {
    if (!questionState.question) {
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
          setQuestionState({
            question: qResponse.data!.question ?? null,
            currentQIndex: qResponse.data!.currentQIndex ?? 0,
            totalQCount: qResponse.data!.totalQCount ?? 0,
          });
        }
      }
      fetchQuestion();
    }
  }, [sessionId]);

  const submitAnswer = async (answerIds: string[]) => {
    const result = await post<SubmitAnswerResponse>(
      apiUrl(`/api/sessions/${sessionId}/answers`),
      null,
      {
        question_id: questionState.question?.ID,
        answer_ids: answerIds,
      },
      true
    );
    if (handleError(result.error, result.status)) return;

    if (result.data!.completed === true) {
      setCompleted(true);
      setResult(result.data!.result);
    } else {
      setQuestionState({
        question: result.data!.question ?? null,
        currentQIndex: result.data!.currentQIndex ?? 0,
        totalQCount: result.data!.totalQCount ?? 0,
      });
    }
  };

  if (completed) {
    return <ResultForm sessionId={sessionId!} resultFrom={result} />;
  }

  if (!questionState.question) {
    return <p>Načítám otázku...</p>;
  }

  return (
    <div className="w-96 h-full">
      <QuestionForm
        text={questionState.question.Text}
        type={questionState.question.Type}
        currentQIndex={questionState.currentQIndex ?? 0}
        totalQCount={questionState.totalQCount ?? 0}
        answers={questionState.question!.Answers}
        submitAnswer={submitAnswer}
      />
    </div>
  );
}
