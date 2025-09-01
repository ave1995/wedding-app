import { useEffect, useState } from "react";
import { useSearchParams, useNavigate } from "react-router-dom";
import { get, post } from "../functions/fetch";
import IconSelector, {
  type SvgItem,
} from "../components/icon_selector/IconSelector";
import Button from "../components/Button";
import type { Quiz } from "../models/Quiz";
import { apiUrl } from "../functions/api";
import type { SimpleResponse } from "../responses/SimpleResponse";
import { useApiErrorHandler } from "../hooks/useApiErrorHandler";
import { FaQuestion, FaExclamation } from "react-icons/fa";

type InviteQuizResult = {
  quiz: Quiz;
  authenticated: boolean;
};

function InvitePage() {
  const { handleError } = useApiErrorHandler();
  const [searchParams] = useSearchParams();
  const code = searchParams.get("code");
  const navigate = useNavigate();

  const [quiz, setQuiz] = useState<Quiz | null>(null);

  //OnMount
  useEffect(() => {
    if (code) {
      async function joinQuiz() {
        // Trying get quiz and also check if i'm authenticated
        const result = await get<InviteQuizResult>(
          apiUrl("/auth/join-quiz"),
          {
            invite: code,
          },
          true
        );
        // Check If I didn't find quiz
        if (handleError(result.error, result.status)) return;

        // I know I find it so lets save it
        setQuiz(result.data!.quiz);
        // Check if I'm authenticated and if so, let's go straight to the quiz
        if (result.data!.authenticated) {
          navigate(`/quiz/${result.data!.quiz!.ID}`);
        }
      }
      joinQuiz();
    } else {
      // TODO: měl bych udělat možnost vyplnit code
      console.error("No code from you!");
      navigate("/not-found");
    }
  }, [code, navigate]);

  // Icon state
  const [selectedIcon, setSelectedIcon] = useState<SvgItem | null>(null);
  const [showSelector, setShowSelector] = useState(false);

  // Input Text state
  const [inputValue, setInputValue] = useState("");
  const [error, setError] = useState("");
  // TODO: add validation
  async function handleJoin() {
    if (!inputValue.trim()) {
      setError("Musíš zadat přezdívku.");
      return;
    }
    if (selectedIcon === null) {
      setShowSelector(true);
      return;
    }

    setError("");
    const result = await post<SimpleResponse>(
      apiUrl("/auth/create-guest"),
      null,
      {
        username: inputValue,
        iconurl: selectedIcon!.URL,
        quizID: quiz!.ID,
      },
      true
    );
    if (handleError(result.error, result.status)) return;

    navigate(`/quiz/${quiz!.ID}`);
  }

  return (
    <div className="flex flex-col w-96 h-screen items-center justify-center p-6 gap-6">
      <h1 className="text-xl font-bold">Nový host</h1>
      <div
        className={`flex border-2 border-b-4 rounded-xl w-72 focus-within:border-[#3D52D5] bg-white/60 border-gray-400
          ${error ? "border-red-500" : ""}
          ${showSelector ? "border-[#3D52D5]" : ""}`}
      >
        <div className="relative">
          <div
            className="w-11 h-11 cursor-pointer"
            onClick={() => setShowSelector((t) => !t)}
          >
            {selectedIcon ? (
              <img src={selectedIcon.URL} alt={selectedIcon.Name} />
            ) : (
              <FaQuestion className="w-full h-full p-3"/>
            )}
          </div>

          {showSelector && (
            <div
              style={{ position: "absolute", top: "100%", left: 0, zIndex: 10 }}
              className="mt-1 w-72 -translate-x-0.5 "
            >
              <IconSelector
                onSelect={(icon) => setSelectedIcon(icon)}
                onClose={() => setShowSelector(false)}
              />
            </div>
          )}
        </div>
        <input
          type="text"
          placeholder="Zadej svojí přezdívku"
          value={inputValue}
          onChange={(e) => {
            setInputValue(e.target.value);
            if (error) setError("");
          }}
          className={`w-full px-3 py-2 text-black border-l-2 border-gray-400 outline-0 placeholder:text-black text-base
            focus-within:border-[#3D52D5] focus-within:placeholder:text-gray-400 caret-[#3D52D5]
            ${error ? "border-red-500" : ""}
            ${showSelector ? "border-[#3D52D5]" : ""}`}
        />
      </div>
      {error && <div className="flex items-center place-content-center gap-1"><FaExclamation className="text-red-500"/><p className="text-red-500 text-sm mt-1">{error}</p></div>}
      <div className="w-72">
        <Button label="Připoj se" onClick={handleJoin}></Button>
      </div>
    </div>
  );
}
export default InvitePage;
