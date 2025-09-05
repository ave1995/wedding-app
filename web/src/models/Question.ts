import type { Answer } from "./Answer";

export type QuestionType = "single_choice" | "multiple_choice";

export type Question = {
  ID: string;
  Text: string;
  Type: QuestionType;
  Answers: Answer[];
};
