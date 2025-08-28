import type { Answer } from "./Answer";

export type Question = {
    ID: string;
    Text: string;
    Type: string;
    Answers: Answer[]
};