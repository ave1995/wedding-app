const _getJsonHeaders = () => {
  const headers = new Headers();
  headers.append("Content-Type", "application/json");
  headers.append("Accept", "application/json");
  return headers;
};

const _getTextHeaders = () => {
  const headers = new Headers();
  headers.append("Content-Type", "text/plain; charset=utf-8");
  headers.append("Accept", "text/plain; charset=utf-8");
  return headers;
};

type FetchResult<T> = {
  data: T | null;
  error: string | null;
};

const _fetch = async <T>(req: {
  url: string;
  method: string;
  func: "json" | "text";
  content?: any;
}): Promise<FetchResult<T>> => {
  let init: RequestInit;

  if (req.content) {
    init = {
      headers: req.func == "json" ? _getJsonHeaders() : _getTextHeaders(),
      method: req.method,
      body: JSON.stringify(req.content),
    };
  } else {
    init = { method: req.method };
  }

  try {
    const response = await fetch(req.url, init);

    if (!response.ok) {
      // Try to parse error message from JSON or fallback to status text
      let errorText = response.statusText;
      try {
        const errorJson = await response.json();
        errorText = errorJson.message || JSON.stringify(errorJson);
      } catch {
        // Ignore JSON parse error here
      }
      return { data: null, error: errorText };
    }

    const data = await response[req.func]();
    return { data, error: null };
  } catch (networkError: any) {
    return { data: null, error: networkError.message || "Network error" };
  }
};

type TContent = Record<any, any> | null | any;

const _parseQuery = (query: Record<any, any>) =>
  Object.keys(query)
    .map((key) => {
      const value = query[key];
      if (Array.isArray(value)) {
        return value.map((s) => `${key}=${encodeURIComponent(s)}`).join("&");
      }
      return `${key}=${encodeURIComponent(query[key])}`;
    })
    .filter((p) => p)
    .join("&");

const _parseUrl = (url: string, query: TContent = null) =>
  query ? `${url}?${_parseQuery(query)}` : url;

export const get = async <T>(url: string, query: TContent = null) =>
  _fetch<T>({
    url: _parseUrl(url, query),
    method: "GET",
    func: "json",
  }) as Promise<FetchResult<T>>;

export const getText = async <T>(url: string, query: TContent = null) =>
  _fetch<T>({
    url: _parseUrl(url, query),
    method: "GET",
    func: "text",
  }) as Promise<FetchResult<T>>;
