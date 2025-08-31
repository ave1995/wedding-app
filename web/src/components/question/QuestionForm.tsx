import { useState } from "react";
import type { Answer } from "../../models/Answer";

interface QuestionForm {
  text: string;
  currentQIndex: number;
  totalQCount: number;
  answers: Answer[];
}

export default function QuestionForm({
  text,
  currentQIndex,
  totalQCount,
  answers,
}: QuestionForm) {
  const [selectedAnswers, setSelectedAnswers] = useState<string[]>([]);

  const toggleAnswer = (id: string) => {
    setSelectedAnswers((prev) =>
      prev.includes(id) ? prev.filter((a) => a !== id) : [...prev, id]
    );
  };

  return (
    <div className="flex flex-col w-full h-full">
      <div className="flex-grow place-items-start">
        <p className="text-sm">
          Otázka {currentQIndex} ze {totalQCount}
        </p>
        <h2 className="text-lg font-semibold">{text}</h2>
      </div>
      <div className="flex-grow">
        <div className="w-full flex flex-col gap-2 items-center justify-center">
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
    focus:outline-none
  `}
              >
                {a.Text}
              </button>
            );
          })}
        </div>
      </div>
      <button
        className="w-full px-5 py-3 text-center border-2 border-b-4 rounded-lg text-white border-pink-300 bg-pink-500
       active:scale-95 active:shadow-sm
       transition-all duration-150 ease-in-out
    focus:outline-none font-medium cursor-pointer"
      >
        Odpovědět
      </button>
    </div>
  );
}
