import type { Question } from "../models/Question";

export type QuestionResponse = {
    completed: boolean;
    question: Question;
}