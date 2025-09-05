import { useState } from "react";
import type { Answer } from "../../models/Answer";
import Button from "../Button";
import QuestionSkeleton from "./QuestionSkeleton";

interface QuestionForm {
  text: string;
  currentQIndex: number;
  totalQCount: number;
  answers: Answer[];
  submitAnswer: (answerId: string[]) => Promise<void>;
}

export default function QuestionForm({
  text,
  currentQIndex,
  totalQCount,
  answers,
  submitAnswer,
}: QuestionForm) {
  const [selectedAnswers, setSelectedAnswers] = useState<string[]>([]);

  const toggleAnswer = (id: string) => {
    setSelectedAnswers((prev) =>
      prev.includes(id) ? prev.filter((a) => a !== id) : [...prev, id]
    );
  };

  return (
    <QuestionSkeleton
      text={text}
      currentQIndex={currentQIndex}
      totalQCount={totalQCount}
      answers={answers}
      button={
        <Button
          label="Odpovědět"
          onClick={() => {
            submitAnswer(selectedAnswers);
            setSelectedAnswers([]);
          }}
        />
      }
      renderAnswer={(a) => {
        const selected = selectedAnswers.includes(a.ID);
        return (
          <button
            key={a.ID}
            onClick={() => toggleAnswer?.(a.ID)}
            className={`
            w-full px-5 py-3 border-2 border-b-4 rounded-lg font-medium cursor-pointer
            transition-all duration-150 ease-in-out
            ${
              selected
                ? " text-pink-500 border-pink-500 shadow-md"
                : " text-gray-800 border-gray-300 hover:bg-gray-100 hover:border-gray-400"
            }
            active:scale-95 active:shadow-sm
            focus:outline-none font-semibold
          `}
          >
            {a.Text}
          </button>
        );
      }}
    ></QuestionSkeleton>
  );
}
