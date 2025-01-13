/**
 * Subset of FetchRequestConfig
 */
export interface RequestConfig<TData = unknown> {
  baseURL?: string;
  url?: string;
  method: 'GET' | 'PUT' | 'PATCH' | 'POST' | 'DELETE' | 'OPTIONS';
  params?: unknown;
  data?: TData | FormData;
  responseType?: 'arraybuffer' | 'blob' | 'document' | 'json' | 'text' | 'stream';
  signal?: AbortSignal;
  headers?: [string, string][] | Record<string, string>;
}
/**
 * Subset of FetchResponse
 */
export interface ResponseConfig<TData = unknown> {
  data: TData;
  status: number;
  statusText: string;
  headers?: [string, string][] | Record<string, string>;
}

export const fetchClient = async <TData, TError = unknown, TVariables = unknown>(config: RequestConfig<TVariables>): Promise<ResponseConfig<TData>> => {
  const searchParams = new URLSearchParams(config.params as string).toString()
  const response = await fetch([config.baseURL, config.url].filter(Boolean).join('') + (searchParams ? `?${searchParams}` : ''), {
    method: config.method.toUpperCase(),
    body: JSON.stringify(config.data),
    signal: config.signal,
    headers: config.headers,
    credentials: 'include'
  })

  const data = (await response.json()) as TData

  return {
    data,
    status: response.status,
    statusText: response.statusText
  }
}

fetchClient.getConfig = () => {
  throw new Error('Not supported')
}
fetchClient.setConfig = () => {
  throw new Error('Not supported')
}

export default fetchClient
