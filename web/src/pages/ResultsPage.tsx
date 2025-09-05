import { useEffect, useMemo, useState } from "react";
import { useParams } from "react-router-dom";
import { get } from "../functions/fetch";
import { useApiErrorHandler } from "../hooks/useApiErrorHandler";
import { apiUrl } from "../functions/api";
import type { UserResult } from "../responses/UserResultsResponse";

type RankedResult = UserResult & { rank: number };

function ResultsPage() {
  const { handleError } = useApiErrorHandler();
  const { quizId } = useParams();
  const [userResults, setResults] = useState<UserResult[]>([]);

  useEffect(() => {
    async function fetchResults() {
      const result = await get<UserResult[]>(
        apiUrl(`/api/quiz/${quizId}/results`),
        null,
        true
      );
      if (handleError(result.error, result.status)) return;

      setResults(result.data!);
    }
    fetchResults();
  }, [quizId]);

  const rankedResults: RankedResult[] = useMemo(() => {
    const sorted = [...userResults].sort(
      (a, b) => b.result.Score - a.result.Score
    );

    const ranked: RankedResult[] = [];
    let currentRank = 1;

    for (let i = 0; i < sorted.length; i++) {
      if (i > 0 && sorted[i].result.Score === sorted[i - 1].result.Score) {
        // same score → same rank
        ranked.push({ ...sorted[i], rank: ranked[i - 1].rank });
      } else {
        // new score → assign current rank
        ranked.push({ ...sorted[i], rank: currentRank });
      }
      // increment rank counter
      currentRank = ranked[i].rank + 1;
    }

    return ranked;
  }, [userResults]);

  return (
    <div className="flex flex-col p-6 w-screen items-center place-content-center">
      <div className="w-1/3">
        <h1 className="text-2xl font-bold p-6">Výsledky</h1>
        <ul className="space-y-2 relative">
          {rankedResults.map((result, idx) => (
            <li
              key={idx}
              className={`relative flex items-center border rounded-xl border-b-4 px-3 py-1 bg-white/60 text-gray-800
        ${
          result.rank === 1
            ? "border-pink-500"
            : result.rank === 2
            ? "border-[#38C172]"
            : result.rank === 3
            ? "border-[#3D52D5]"
            : idx % 2 === 0
            ? "border-gray-300"
            : "border-gray-400"
        }`}
            >
              {/* User */}
              <div className="flex items-center gap-2">
                <img
                  className="w-11 h-11"
                  src={result.user.IconUrl}
                  alt={result.user.Username}
                />
                <span className="font-medium text-gray-800">
                  {result.user.Username}
                </span>
              </div>

              {/* Medal (absolute center) */}
              <div className="absolute left-1/2 transform -translate-x-1/2">
                {result.rank === 1 && (
                  <span className="text-yellow-500 text-3xl">🥇</span>
                )}
                {result.rank === 2 && (
                  <span className="text-gray-400 text-3xl">🥈</span>
                )}
                {result.rank === 3 && (
                  <span className="text-amber-600 text-3xl">🥉</span>
                )}
              </div>

              {/* Score */}
              <p className="ml-auto font-semibold">
                Skóre: {result.result.Score}/{result.result.Total} (
                {result.result.Percentage}%)
              </p>
            </li>
          ))}
        </ul>
      </div>
    </div>
  );
}

export default ResultsPage;
