import type { Question } from "../models/Question";

export type RevealResponse = {
  question: Question;
  goNext: boolean;
  nextIndex: number;
  totalQCount: number;
};
