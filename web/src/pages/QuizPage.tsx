import { useEffect, useState } from "react";
import { useLocation, useParams } from "react-router-dom";
import type { Quiz } from "../models/Quiz";

function QuizPage() {
  const { state } = useLocation();
  const { quizId } = useParams();
  const [quiz, setQuiz] = useState<Quiz | null>(state?.quiz || null);

//   useEffect(() => {});
  return <div>{quiz?.ID}</div>;
}

export default QuizPage;
