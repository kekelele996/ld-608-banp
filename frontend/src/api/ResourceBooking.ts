import { request } from "./client";
import type { RebookCandidate, RebookResult, BookingHistoryItem } from "../types/ResourceBooking";

// rebookBooking swaps a conflicting booking onto another usable resource.
// The backend flips the old booking to RELEASED and creates a linked new one
// inside one transaction.
export function rebookBooking(
  bookingId: number,
  payload: { new_resource_id: number; new_start_time?: string; new_end_time?: string }
): Promise<RebookResult> {
  return request<RebookResult>(`/api/bookings/${bookingId}/rebook`, {
    method: "POST",
    body: JSON.stringify({ actor: "dispatcher", ...payload })
  });
}

export function listRebookCandidates(bookingId: number): Promise<RebookCandidate[]> {
  return request<RebookCandidate[]>(`/api/bookings/${bookingId}/candidates`);
}

export function listBookingHistory(taskId: number): Promise<BookingHistoryItem[]> {
  return request<BookingHistoryItem[]>(`/api/bookings/task/${taskId}/history`);
}
