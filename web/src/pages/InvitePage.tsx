import { useEffect, useState } from "react";
import { useSearchParams, useNavigate } from "react-router-dom";
import { get } from "../functions/fetch";
import InputText from "../components/input_text/InputText";
import IconSelector, {
  type SvgItem,
} from "../components/icon_selector/IconSelector";
import Button from "../components/Button";
import type { Quiz } from "../models/Quiz";
import { apiUrl } from "../functions/api";

type InviteQuizResult = {
  quiz: Quiz;
  authenticated: boolean;
};

function InvitePage() {
  const [searchParams] = useSearchParams();
  const code = searchParams.get("code");
  const navigate = useNavigate();

  const [quiz, setQuiz] = useState<Quiz | null>(null);

  //OnMount
  useEffect(() => {
    if (code) {
      async function fetchQuiz() {
        // Trying get quiz and also check if i am authenticated
        const result = await get<InviteQuizResult>(
          apiUrl("/auth/join-quiz"),
          {
            invite: code,
          },
          true
        );
        // Check If I didn't find quiz
        if (result.error) {
          console.error(result.error);
          navigate("/not-found");
        }
        // I know I find it so lets save it
        setQuiz(result.data!.quiz);
        // Check if I am authenticated and if so, let's go straight to the quiz
        if (result.data!.authenticated) {
          navigate(`/quiz/${quiz!.ID}`, {
            state: { quiz: quiz! },
          });
        }
      }
      fetchQuiz();
    } else {
      console.error("No code from you!");
      navigate("/not-found");
    }
  }, [code, navigate]);

  // Icon state
  const [selectedIcon, setSelectedIcon] = useState<SvgItem | null>(null);
  const [showSelector, setShowSelector] = useState(false);

  // Input Text state
  const [inputValue, setInputValue] = useState("");

  async function handleJoin() {}

  return (
    <div>
      <h2>New Guest</h2>
      <div className="flex">
        <div className="relative">
          <div
            className="w-12 h-12 cursor-pointer border-2 rounded-lg"
            onClick={() => setShowSelector(true)}
          >
            {selectedIcon && (
              <img src={selectedIcon.URL} alt={selectedIcon.Name} />
            )}
          </div>

          {showSelector && (
            <div
              style={{ position: "absolute", top: "100%", left: 0, zIndex: 10 }}
              className="mt-1 w-max"
            >
              <IconSelector
                onSelect={(icon) => setSelectedIcon(icon)}
                onClose={() => setShowSelector(false)}
              />
            </div>
          )}
        </div>
        <InputText
          placeholder="Enter your name"
          value={inputValue}
          onChange={(e) => setInputValue(e.target.value)}
        />
      </div>
      <Button label="Join" onClick={handleJoin}></Button>
    </div>
  );
}
export default InvitePage;
