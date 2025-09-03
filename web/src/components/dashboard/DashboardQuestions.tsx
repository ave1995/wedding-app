import {
  questionStatusColors,
  type QuestionInfo,
} from "../../hooks/useQuestionEventHandler";

interface DashboardQuestions {
  questions: QuestionInfo[];
}

export default function DashboardQuestions({ questions }: DashboardQuestions) {
  return (
    <div className="flex flex-col gap-2 p-4 text-gray-800 text-left h-[600px] overflow-auto">
      <h2 className="text-lg font-semibold">Aktivita</h2>
      <ul className="space-y-2">
        {questions.map((msg, idx) => (
          <li
            key={idx}
            className={`flex gap-4 items-center place-content-between border rounded-xl border-b-4 px-3 py-1 bg-white/60  text-gray-800
                ${idx % 2 === 0 ? "border-gray-300" : "border-gray-400"}`}
          >
            {/* Uživatel */}
            <div className="flex gap-2 items-center">
              <img
                className="w-11 h-11"
                src={msg.UserIconUrl}
                alt={msg.Username}
              />
              <span className="font-medium text-gray-800">{msg.Username}</span>
            </div>
            {/* Otázka */}
            <p className="flex-grow">
              <span className="font-semibold">Otázka:</span> {msg.QuestionText}
            </p>
            {/* Status */}
            <div className="flex">
              <span
                className={`flex items-center place-content-center min-w-[102px] py-1 ${
                  questionStatusColors[msg.Status]
                } rounded-xl`}
              >
                {msg.Status}
              </span>
            </div>
          </li>
        ))}
      </ul>
    </div>
  );
}
