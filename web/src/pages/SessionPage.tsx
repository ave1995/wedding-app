import { useParams, useLocation } from "react-router-dom";
import { useEffect, useState } from "react";
import type { Question } from "../models/Question";
import { get, post } from "../functions/fetch";
import { apiUrl } from "../functions/api";
import type { Result } from "../models/Result";
import type { SubmitAnswerResponse } from "../responses/SubmitAnswerResponse";

export default function SessionPage() {
  const { state } = useLocation();
  const { sessionId } = useParams();

  const [question, setQuestion] = useState<Question | null>(state?.question);
  const [completed, setCompleted] = useState(false);
  const [result, setResult] = useState<Result | null>(null);

  useEffect(() => {
    if (!question) {
      async function fetchQuestion() {
        const result = await get<Question>(
          apiUrl(`/api/sessions/${sessionId}/question`),
          null,
          true
        );
        if (result.error) {
          console.error(result.error);
        } else {
          setQuestion(result.data);
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
        answer_id: answerId,
      },
      true
    );
    if (result.error) {
      console.error(result.error);
      return;
    }
    if (result.data?.completed === true) {
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
          Skóre: {result?.Score}/{result?.Total} ({result?.Percentage}%)
        </p>
      </div>
    );
  }

  if (!question) {
    return <p>Načítám otázku...</p>;
  }

  return (
    <div className="p-6">
      <h2 className="text-lg font-semibold">{question.Text}</h2>
      <div className="mt-4 flex flex-col gap-2">
        {question.Answers.map((a) => (
          <button
            key={a.ID}
            onClick={() => submitAnswer(a.ID)}
            className="px-4 py-2 bg-gray-200 rounded-lg hover:bg-gray-300"
          >
            {a.Text}
          </button>
        ))}
      </div>
    </div>
  );
}
