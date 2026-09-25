import type { ApiErrorBody } from "../types/FlightTurnaround";

// request wraps the uniform { ok, data } / { ok:false, code, message }
// envelope. Non-2xx responses are thrown as ApiError carrying the server
// violation list so pages can pinpoint the blocking flight/task/resource.
export class ApiError extends Error {
  status: number;
  code: string;
  violations: ApiErrorBody["violations"];

  constructor(status: number, body: ApiErrorBody | null) {
    super(body?.message ?? `请求失败（HTTP ${status}）`);
    this.status = status;
    this.code = body?.code ?? "NETWORK_ERROR";
    this.violations = body?.violations ?? (body?.violation ? [body.violation] : undefined);
  }
}

export async function request<T>(input: string, init?: RequestInit): Promise<T> {
  const res = await fetch(input, {
    headers: { "Content-Type": "application/json", "X-Actor": "dispatcher", ...(init?.headers ?? {}) },
    ...init
  });
  let body: unknown = null;
  try {
    body = await res.json();
  } catch {
    // keep null for non-JSON responses
  }
  if (!res.ok) {
    throw new ApiError(res.status, body as ApiErrorBody | null);
  }
  const envelope = body as { ok: boolean; data: T };
  return envelope.data;
}
