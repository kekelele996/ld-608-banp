import { mockData } from "../mocks/seedData";
import { buildReleaseSummaryResponse } from "../constructors/TurnaroundReleaseConstructor";
import { ERROR_MESSAGES } from "../constants/errorMessages";
import type { ApiError, ReleaseResult, ReleaseSummary } from "../types/TurnaroundRelease";
import type { FlightTurnaround } from "../types/FlightTurnaround";
import type { GroundTask } from "../types/GroundTask";
import type { GroundResource } from "../types/GroundResource";
import type { ResourceBooking } from "../types/ResourceBooking";
import type { DelayEvent } from "../types/DelayEvent";

const endpoint = "/api/flight-turnaround";

const mockSummary = (turnaroundId: number): ReleaseSummary => {
  const turnaround = (mockData.flightTurnaround as unknown as FlightTurnaround[]).find((t) => t.id === turnaroundId);
  if (!turnaround) {
    throw { code: "TURNAROUND_NOT_FOUND", message: ERROR_MESSAGES.TURNAROUND_NOT_FOUND } as ApiError;
  }
  return buildReleaseSummaryResponse(
    turnaround,
    mockData.groundTask as unknown as GroundTask[],
    mockData.delayEvent as unknown as DelayEvent[],
    mockData.resourceBooking as unknown as ResourceBooking[],
    mockData.groundResource as unknown as GroundResource[]
  );
};

// fetchReleaseSummary 航班详情汇总：未完成任务、未关闭延误、预约冲突与放行结论。
export async function fetchReleaseSummary(turnaroundId: number): Promise<ReleaseSummary> {
  try {
    const res = await fetch(`${endpoint}/${turnaroundId}/release-summary`);
    const body = await res.json();
    if (!res.ok) throw (body.error ?? body) as ApiError;
    return body as ReleaseSummary;
  } catch (err) {
    if ((err as ApiError)?.code) throw err;
    // Local mock fallback keeps the UI available during offline review.
    return mockSummary(turnaroundId);
  }
}

// submitRelease 提交放行：后端重新核对三类条件，不满足则抛出带 blockers 的 ApiError。
export async function submitRelease(turnaroundId: number): Promise<ReleaseResult> {
  try {
    const res = await fetch(`${endpoint}/${turnaroundId}/release`, { method: "POST" });
    const body = await res.json();
    if (!res.ok) throw (body.error ?? body) as ApiError;
    return body as ReleaseResult;
  } catch (err) {
    if ((err as ApiError)?.code) throw err;
    // Local mock fallback: 离线时按同一套规则给出结论。
    const summary = mockSummary(turnaroundId);
    if (summary.conclusion === "BLOCKED") {
      throw { code: "RELEASE_CHECK_FAILED", message: ERROR_MESSAGES.RELEASE_CHECK_FAILED, blockers: summary.blockers } as ApiError;
    }
    return { ok: true, turnaround: { ...summary.turnaround, turnaround_status: "READY" } };
  }
}
