import type { User } from "../models/User";

export type UserResult = {
  result: DeleteThisStupidResult;
  user: User;
};

export type DeleteThisStupidResult = {
  Score: number;
  Total: number;
  Percentage: number;
};
