import type { ReactNode } from "react";
import type { Answer } from "../../models/Answer";
import type { QuestionType } from "../../models/Question";

interface QuestionSkeleton {
  text: string;
  type: QuestionType;
  currentQIndex: number;
  totalQCount: number;
  answers: Answer[];
  button: ReactNode;
  renderAnswer: (a: Answer) => ReactNode;
}

export default function QuestionSkeleton({
  text,
  type,
  currentQIndex,
  totalQCount,
  answers,
  button,
  renderAnswer,
}: QuestionSkeleton) {
  return (
    <div className="flex flex-col w-full h-full">
      <div className="flex-grow p-6">
        <div className="flex place-content-between items-center">
          <p className="text-xs italic">
            Otázka <span className="text-pink-500">{currentQIndex}</span> ze{" "}
            <span className="text-pink-500">{totalQCount}</span>
          </p>
          <p className="w-7 h-7 bg-pink-500 text-white text-sm rounded-full flex items-center justify-center">
            {type === "multiple_choice" ? "2+" : "1"}
          </p>
        </div>
        <h2 className="text-lg font-semibold text-left">{text}</h2>
      </div>
      <div className="flex-grow p-6">
        <div className="w-full flex flex-col gap-2 items-center justify-center overflow-x-auto">
          {answers.map((a) => {
            return (
              <div key={a.ID} className="w-full">
                {renderAnswer(a)}
              </div>
            );
          })}
        </div>
      </div>
      <div className="border-t-2 border-gray-300 p-6">{button}</div>
    </div>
  );
}
