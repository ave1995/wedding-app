import { useState, useCallback, useMemo } from "react";

export type QuestionStatus = "Přemýšlí" | "Odpověděl";

export const questionStatusColors: Record<QuestionStatus, string> = {
  Přemýšlí: "bg-sky-200",
  Odpověděl: "bg-green-500",
};

export type QuestionInfo = {
  QuestionID: string;
  UserID: string;
  QuestionText: string;
  UserIconUrl: string;
  Username: string;
  Status: QuestionStatus;
};

export function useQuestions() {
  const [questionsMap, setQuestionsMap] = useState<Map<string, QuestionInfo>>(
    new Map()
  );

  const upsertQuestion = useCallback((newQuestion: QuestionInfo) => {
    setQuestionsMap((prev) => {
      const newMap = new Map(prev);
      // remove existing so reinsertion updates order
      newMap.delete(newQuestion.UserID);
      newMap.set(newQuestion.UserID, newQuestion);
      return newMap;
    });
  }, []);

  // Convert Map -> array (latest first)
  const questions = useMemo(
    () => Array.from(questionsMap.values()).reverse(),
    [questionsMap]
  );

  return { questions, upsertQuestion };
}
