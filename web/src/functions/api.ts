const API_BASE_URL = import.meta.env.VITE_API_URL as string;

export const apiUrl = (path: string) => `${API_BASE_URL}${path}`;