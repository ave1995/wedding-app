import type { Question } from "../models/Question";
import type { Result } from "../models/Result";

export type SubmitAnswerResponseCompleted = {
  completed: true;
  result: Result;
};

export type SubmitAnswerResponsePending = {
  completed: false;
  session_id: string;
  question: Question;
  currentQIndex: number;
  totalQCount: number;
};

export type SubmitAnswerResponse =
  | SubmitAnswerResponseCompleted
  | SubmitAnswerResponsePending;
