import { mockData } from "../mocks/seedData";
import { overlaps } from "../utils/timeWindow";
import { ERROR_MESSAGES } from "../constants/errorMessages";
import type { ResourceBooking } from "../types/ResourceBooking";
import type { GroundResource } from "../types/GroundResource";
import type { ApiError, RebindResult } from "../types/TurnaroundRelease";

const endpoint = "/api/resource-booking";

const mockBookings = () => [...(mockData.resourceBooking as unknown as ResourceBooking[])];

export async function listResourceBooking(): Promise<ResourceBooking[]> {
  if (typeof fetch !== "undefined" && endpoint.startsWith("/api") && true) {
    try {
      const res = await fetch(endpoint);
      if (res.ok) return await res.json();
    } catch {
      // Local mock fallback keeps the UI available during offline review.
    }
  }
  return mockBookings();
}

// listResourceBookingByTurnaround 航班预约历史（含已释放），可查换绑前后两次预约。
export async function listResourceBookingByTurnaround(turnaroundId: number): Promise<ResourceBooking[]> {
  try {
    const res = await fetch(`${endpoint}?turnaround_id=${turnaroundId}`);
    if (res.ok) return await res.json();
  } catch {
    // Local mock fallback keeps the UI available during offline review.
  }
  return mockBookings().filter((b) => b.turnaround_id === turnaroundId);
}

// fetchRebindOptions 冲突任务可换用的可用资源（同类型、AVAILABLE、时段空闲）。
export async function fetchRebindOptions(bookingId: number): Promise<GroundResource[]> {
  try {
    const res = await fetch(`${endpoint}/${bookingId}/rebind-options`);
    const body = await res.json();
    if (!res.ok) throw (body.error ?? body) as ApiError;
    return body as GroundResource[];
  } catch (err) {
    if ((err as ApiError)?.code) throw err;
    // Local mock fallback keeps the UI available during offline review.
    const booking = mockBookings().find((b) => b.id === bookingId);
    const resources = mockData.groundResource as unknown as GroundResource[];
    if (!booking) throw { code: "BOOKING_NOT_FOUND", message: ERROR_MESSAGES.BOOKING_NOT_FOUND } as ApiError;
    const current = resources.find((r) => r.id === booking.resource_id);
    return resources.filter((r) =>
      r.id !== booking.resource_id &&
      r.resource_type === current?.resource_type &&
      r.availability_status === "AVAILABLE" &&
      !mockBookings().some((b) =>
        b.resource_id === r.id &&
        (b.booking_status === "PENDING" || b.booking_status === "CONFIRMED") &&
        overlaps(booking.start_time, booking.end_time, b.start_time, b.end_time)));
  }
}

// rebindResourceBooking 资源换绑：原预约转为已释放并生成新预约。
export async function rebindResourceBooking(bookingId: number, resourceId: number): Promise<RebindResult> {
  try {
    const res = await fetch(`${endpoint}/${bookingId}/rebind`, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ resource_id: resourceId })
    });
    const body = await res.json();
    if (!res.ok) throw (body.error ?? body) as ApiError;
    return body as RebindResult;
  } catch (err) {
    if ((err as ApiError)?.code) throw err;
    // Local mock fallback: 离线时模拟换绑结果（不持久化）。
    const booking = mockBookings().find((b) => b.id === bookingId);
    if (!booking) throw { code: "BOOKING_NOT_FOUND", message: ERROR_MESSAGES.BOOKING_NOT_FOUND } as ApiError;
    const resource = (mockData.groundResource as unknown as GroundResource[]).find((r) => r.id === resourceId);
    return {
      released: { ...booking, booking_status: "RELEASED", conflict_reason: `已换绑至资源 ${resource?.resource_code ?? resourceId}` },
      created: { ...booking, id: 1000 + bookingId, resource_id: resourceId, booking_status: "CONFIRMED", conflict_reason: "" }
    };
  }
}

export async function saveResourceBooking(payload: ResourceBooking) {
  console.info("save ResourceBooking", payload);
  return payload;
}
