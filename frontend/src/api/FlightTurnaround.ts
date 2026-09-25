import { request } from "./client";
import type { FlightTurnaround, ReleaseCheckResult, ReleaseResponse } from "../types/FlightTurnaround";

const endpoint = "/api/turnarounds";

export function listTurnarounds(): Promise<FlightTurnaround[]> {
  return request<FlightTurnaround[]>(endpoint);
}

export function getTurnaround(id: number) {
  return request<{
    turnaround: FlightTurnaround;
    tasks: import("../types/GroundTask").GroundTask[];
    delays: import("../types/DelayEvent").DelayEvent[];
    bookings: import("../types/ResourceBooking").BookingHistoryItem[];
  }>(`${endpoint}/${id}`);
}

// getReleaseCheck aggregates unfinished tasks, open delays and booking
// conflicts for the flight detail release panel (read-only).
export function getReleaseCheck(id: number): Promise<ReleaseCheckResult> {
  return request<ReleaseCheckResult>(`${endpoint}/${id}/release-check`);
}

// submitRelease re-verifies all three conditions server-side. On failure the
// rejected ApiError carries violations naming the exact tasks/delays/bookings.
export function submitRelease(id: number, actor = "dispatcher"): Promise<ReleaseResponse> {
  return request<ReleaseResponse>(`${endpoint}/${id}/release`, {
    method: "POST",
    body: JSON.stringify({ actor })
  });
}
