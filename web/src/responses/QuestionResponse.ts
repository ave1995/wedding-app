import type { Question } from "../models/Question";

// TODO: duplicity with SubmitAnswerResponsePending
export type QuestionResponse = {
  session_id: string;
  completed: boolean;
  question: Question;
  currentQIndex: number;
  totalQCount: number;
};
