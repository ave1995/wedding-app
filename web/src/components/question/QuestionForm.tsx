import { useState } from "react";
import type { Answer } from "../../models/Answer";
import Button from "../Button";

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
  submitAnswer
}: QuestionForm) {
  const [selectedAnswers, setSelectedAnswers] = useState<string[]>([]);

  const toggleAnswer = (id: string) => {
    setSelectedAnswers((prev) =>
      prev.includes(id) ? prev.filter((a) => a !== id) : [...prev, id]
    );
  };

  return (
    <div className="flex flex-col w-full h-full">
      <div className="flex-grow place-items-start p-6">
        <p className="text-xs italic">
          Otázka <span className="text-pink-500">{currentQIndex}</span> ze  <span className="text-pink-500">{totalQCount}</span>
        </p>
        <h2 className="text-lg font-semibold">{text}</h2>
      </div>
      <div className="flex-grow p-6">
        <div className="w-full flex flex-col gap-2 items-center justify-center overflow-x-auto">
          {answers.map((a) => {
            const isSelected = selectedAnswers.includes(a.ID);
            return (
              <button
                key={a.ID}
                onClick={() => toggleAnswer(a.ID)}
                className={`
    w-full px-5 py-3 border-2 border-b-4 rounded-lg font-medium cursor-pointer
    transition-all duration-150 ease-in-out
    ${
      isSelected
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
          })}
        </div>
      </div>
      <div className="border-t-2 border-gray-300 p-6">
        <Button label="Odpovědět" onClick={() => submitAnswer(selectedAnswers)} />
      </div>
    </div>
  );
}
