import { useNavigate } from "react-router-dom";

export function useApiErrorHandler() {
  const navigate = useNavigate();

  function handleError(error: string | null, status: number): boolean {
    if (error === null) return false;

    console.error(`I handle an error: ${error}`);
    
    if (status === 401 || status === 403) {
      navigate("/unauthorized");
    } else if (status === 404) {
      navigate("/not-found");
    } else if (status === -1) {
      navigate("/network-error");
    }

    return true;
  }

  return { handleError };
}
