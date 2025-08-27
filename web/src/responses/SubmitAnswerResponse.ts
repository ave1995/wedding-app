import type { Question } from "../models/Question";
import type { Result } from "../models/Result";

export type SubmitAnswerResponseCompleted = {
    completed: true;
    result: Result
}

export type SubmitAnswerResponsePending = {
    completed: false;
    nextQuestion: Question;
}

export type SubmitAnswerResponse = SubmitAnswerResponseCompleted | SubmitAnswerResponsePending;