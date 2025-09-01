import type { Result } from "../models/Result";
import Thanks from "../assets/Thanks.png";
import { useEffect, useState } from "react";
import confetti from "canvas-confetti";
import type { SubmitAnswerResponseCompleted } from "../responses/SubmitAnswerResponse";
import { get } from "../functions/fetch";
import { apiUrl } from "../functions/api";
import { useApiErrorHandler } from "../hooks/useApiErrorHandler";

interface ResultForm {
  sessionId: string;
  resultFrom: Result | null;
}

export default function ResultForm({ sessionId, resultFrom }: ResultForm) {
  const [result, setResult] = useState<Result | null>(resultFrom);
  const { handleError } = useApiErrorHandler();

  useEffect(() => {
    if (!result) {
      async function fetchResult() {
        const result = await get<SubmitAnswerResponseCompleted>(
          apiUrl(`/api/sessions/${sessionId}/result`),
          null,
          true
        );
        if (handleError(result.error, result.status)) return;

        console.log(result.data!.result);
        setResult(result.data!.result);
      }
      fetchResult();
    }
  }, [result]);

  useEffect(() => {
    const interval = setInterval(() => {
      confetti({
        particleCount: 8,
        angle: 60,
        spread: 55,
        origin: { x: 0 },
      });
      confetti({
        particleCount: 8,
        angle: 120,
        spread: 55,
        origin: { x: 1 },
      });
    }, 500);

    return () => clearInterval(interval); // cleanup
  }, []);

  return (
    <div className="flex flex-col w-96 h-screen items-center justify-center p-6 gap-4">
      <h1 className="text-xl font-bold">Hotovo!</h1>
      <p className="text-left p-3 font-medium">
        <span className="font-semibold text-pink-500">Děkujeme</span>, že jste se s námi zapojili do kvízu! Doufáme, že vás pobavil a
        třeba jste se o nás i něco nového dozvěděli. Teď už si pojďme společně
        dál užívat <span className="font-semibold text-pink-500">naši svatbu</span> – ať je plná smíchu, tance a krásných vzpomínek. S láskou <span className="font-semibold text-pink-500">Bednářovi</span>!
      </p>
      <p className="font-semibold">
        Skóre: {result?.score}/{result?.total} ({result?.percentage}%)
      </p>
      <img src={Thanks} alt="Thanks" className="w-36 h-28"></img>
    </div>
  );
}
