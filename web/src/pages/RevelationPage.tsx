import { useEffect, useState } from "react";
import type { Quiz } from "../models/Quiz";
import { apiUrl } from "../functions/api";
import { get } from "../functions/fetch";
import { useApiErrorHandler } from "../hooks/useApiErrorHandler";
import { useParams } from "react-router-dom";
import type { Question } from "../models/Question";
import type { RevealResponse } from "../responses/RevealResponse";
import QuestionReveal from "../components/question/QuestionReveal";
import QuestionStats from "../components/question/QuestionStats";
import QuestionPhoto from "../components/question/QuestionPhoto";
import FullscreenImageButton from "../components/FullScreenImageButton.tsx";

type QuestionPhotoImg = {
  src: string;
  name: string;
};

function pathToBucketQuery(path: string): QuestionPhotoImg | null {
  if (!path) return null;

  const parts = path.split("/");
  const name = parts.pop() || "";
  const bucket = parts.join("/");

  const src = apiUrl(
    `/bucket-data?bucket=${encodeURIComponent(
      bucket
    )}&name=${encodeURIComponent(name)}`
  );

  return { src, name };
}

function RevelationPage() {
  const { handleError } = useApiErrorHandler();
  const { quizId } = useParams();
  const [quiz, setQuiz] = useState<Quiz | null>(null);
  const [questionState, setQuestionState] = useState<{
    question: Question | null;
    nextIndex: number;
    totalQCount: number;
  }>({ question: null, totalQCount: 0, nextIndex: 0 });

  const [completed, setCompleted] = useState(false);

  const fetchResultOfQuestion = async function () {
    const result = await get<RevealResponse>(
      apiUrl(`/api/quiz/${quizId}/reveal`),
      { index: questionState.nextIndex },
      true
    );
    if (handleError(result.error, result.status)) return;

    if (!result.data?.goNext) {
      setCompleted(true);
      return;
    }

    setQuestionState({
      question: result.data!.question,
      nextIndex: result.data!.nextIndex,
      totalQCount: result.data.totalQCount, // TODO: tohle tady nemá být
    });
  };

  useEffect(() => {
    if (!quiz) {
      async function fetchQuiz() {
        const result = await get<Quiz>(apiUrl(`/api/quiz/${quizId}`), {}, true);
        if (handleError(result.error, result.status)) return;

        setQuiz(result.data);
      }
      fetchQuiz();
    }
    if (!questionState.question) {
      fetchResultOfQuestion();
    }
  }, [quizId]);

  if (completed) {
    return <div></div>;
  }

  if (!questionState.question) {
    return <p>Načítám otázku...</p>;
  }

  const questionPhotoImg = pathToBucketQuery(questionState.question.PhotoPath);

  return (
    <div className="flex flex-row h-screen w-screen items-center place-content-center">
      <div className="w-96 h-2/3">
        <QuestionReveal
          text={questionState.question.Text}
          type={questionState.question.Type}
          currentQIndex={questionState.nextIndex}
          totalQCount={questionState.totalQCount}
          answers={questionState.question.Answers}
          nextQuestion={fetchResultOfQuestion}
        ></QuestionReveal>
      </div>
      <div className="w-[768px] grid grid-rows-11 items-center place-items-center h-2/3 gap-2">
        <div className="h-full w-full flex flex-col items-center row-span-5">
          {questionPhotoImg && (
            <QuestionPhoto
              src={questionPhotoImg.src}
              alt={questionPhotoImg.name}
            />
          )}
        </div>
        <div className="w-1/6 row-span-1">
          {questionPhotoImg && (
            <FullscreenImageButton src={questionPhotoImg.src} label="Zvětšit" />
          )}
        </div>
        <div className="row-span-5 w-1/2 h-3/4">
          <QuestionStats id={questionState.question.ID} />
        </div>
      </div>
    </div>
  );
}

export default RevelationPage;
