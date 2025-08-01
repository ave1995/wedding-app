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

const _fetch = async <T>(req: {
  url: string;
  method: string;
  func: "json" | "text";
  content?: any;
  raw: boolean;
}) => {
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

  const response = await fetch(req.url, init);
  if (req.raw) {
    return response;
  }

  return (await response[req.func]()) as T;
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
    raw: false,
  }) as Promise<T>;

export const getText = async <T>(url: string, query: TContent = null) =>
  _fetch<T>({
    url: _parseUrl(url, query),
    method: "GET",
    func: "text",
    raw: false,
  }) as Promise<T>;
