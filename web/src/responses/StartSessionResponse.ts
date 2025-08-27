import type { Question } from "../models/Question";

export type StartSessionResponse = {
    session_id: string;
    question: Question;
}