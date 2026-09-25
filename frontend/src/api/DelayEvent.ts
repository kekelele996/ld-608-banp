import { request } from "./client";
import type { DelayEvent } from "../types/DelayEvent";

export function listDelayEvents(): Promise<DelayEvent[]> {
  return request<DelayEvent[]>("/api/delay-events");
}

export function resolveDelay(id: number): Promise<DelayEvent> {
  return request<DelayEvent>(`/api/delay-events/${id}/resolve`, { method: "POST", body: "{}" });
}
