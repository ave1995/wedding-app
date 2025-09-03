import { apiUrl } from "../functions/api";
import { get } from "../functions/fetch";
import { useApiErrorHandler } from "./useApiErrorHandler";

type AuthCheck = {
  status: string;
};

export function useAuthCheck() {
  const { handleError } = useApiErrorHandler();

  return async function fetchAuthCheck(): Promise<boolean> {
    const result = await get<AuthCheck>(apiUrl("/api/auth-check"), null, true);
    if (handleError(result.error, result.status)) return false;
    return true;
  };
}
