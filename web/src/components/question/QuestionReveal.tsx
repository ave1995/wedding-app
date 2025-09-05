import type { Answer } from "../../models/Answer";
import type { QuestionType } from "../../models/Question";
import Button from "../Button";
import QuestionSkeleton from "./QuestionSkeleton";

interface QuestionReveal {
  text: string;
  type: QuestionType;
  currentQIndex: number;
  totalQCount: number;
  answers: Answer[];
  nextQuestion: () => Promise<void>;
}

export default function QuestionReveal({
  text,
  type,
  currentQIndex,
  totalQCount,
  answers,
  nextQuestion,
}: QuestionReveal) {
  return (
    <QuestionSkeleton
      text={text}
      type={type}
      currentQIndex={currentQIndex}
      totalQCount={totalQCount}
      answers={answers}
      button={
        currentQIndex >= totalQCount ? (
          <span className="invisible">
            <Button label="Pokračovat" onClick={() => {}}></Button>
          </span>
        ) : (
          <Button label="Pokračovat" onClick={nextQuestion} />
        )
      }
      renderAnswer={(a) => {
        return (
          <button
            key={a.ID}
            className={`
                w-full px-5 py-3 border-2 border-b-4 rounded-lg font-medium
                transition-all duration-150 ease-in-out
                ${
                  a.IsCorrect
                    ? " text-white border-[#5DD48F] bg-[#38C172] shadow-md"
                    : " text-gray-800 border-gray-300"
                }
               
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
