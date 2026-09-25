import type { FlightTurnaround } from "../types/FlightTurnaround";
import type { GroundTask } from "../types/GroundTask";
import type { GroundResource } from "../types/GroundResource";
import type { ResourceBooking } from "../types/ResourceBooking";
import type { DelayEvent } from "../types/DelayEvent";
import type { BookingConflict, ReleaseBlocker, ReleaseSummary } from "../types/TurnaroundRelease";
import { overlaps } from "../utils/timeWindow";

export const createReleaseBlocker = (kind: ReleaseBlocker["kind"], refId: number, message: string): ReleaseBlocker => ({
  kind,
  ref_id: refId,
  message
});

export const createEmptyReleaseSummary = (turnaround: FlightTurnaround): ReleaseSummary => ({
  turnaround,
  unfinished_tasks: [],
  open_delays: [],
  booking_conflicts: [],
  conclusion: "RELEASABLE",
  blockers: []
});

const isActiveBooking = (b: ResourceBooking) => b.booking_status === "PENDING" || b.booking_status === "CONFIRMED";

// buildReleaseSummaryResponse 由原始实体构造放行汇总响应对象，
// 与后端判定规则保持一致，供离线回退使用，页面/store 不得自行拼装。
export function buildReleaseSummaryResponse(
  turnaround: FlightTurnaround,
  tasks: GroundTask[],
  delays: DelayEvent[],
  bookings: ResourceBooking[],
  resources: GroundResource[]
): ReleaseSummary {
  const summary = createEmptyReleaseSummary(turnaround);
  summary.unfinished_tasks = tasks.filter((t) => t.turnaround_id === turnaround.id && t.status !== "COMPLETED");
  summary.unfinished_tasks.forEach((t) => summary.blockers.push(
    createReleaseBlocker("TASK", t.id, `航班 ${turnaround.flight_no} 任务 #${t.id}（${t.task_type}）未完成，当前状态 ${t.status}`)
  ));
  summary.open_delays = delays.filter((d) => d.turnaround_id === turnaround.id && !d.resolved_at);
  summary.open_delays.forEach((d) => summary.blockers.push(
    createReleaseBlocker("DELAY", d.id, `航班 ${turnaround.flight_no} 延误事件 #${d.id}（${d.delay_type}，${d.minutes} 分钟）未关闭`)
  ));
  const mine = bookings.filter((b) => b.turnaround_id === turnaround.id && isActiveBooking(b));
  mine.forEach((b) => {
    bookings.filter((o) => o.id !== b.id && o.resource_id === b.resource_id && isActiveBooking(o))
      .filter((o) => overlaps(b.start_time, b.end_time, o.start_time, o.end_time))
      .forEach((o) => {
        const resource = resources.find((r) => r.id === b.resource_id) ?? resources[0];
        const conflict: BookingConflict = { booking: b, conflicts_with: o, resource, reason: `${b.start_time}~${b.end_time}` };
        summary.booking_conflicts.push(conflict);
        summary.blockers.push(createReleaseBlocker("BOOKING", b.id,
          `资源 ${resource?.resource_code ?? b.resource_id} 预约 #${b.id} 与预约 #${o.id} 时段重叠（${conflict.reason}）`));
      });
  });
  summary.conclusion = summary.blockers.length > 0 ? "BLOCKED" : "RELEASABLE";
  return summary;
}
