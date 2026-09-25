import type { ConflictBooking } from "../../types/FlightTurnaround";
import { formatTaskType, formatTime } from "../../utils/formatters";
import { StatusBadge } from "../common/StatusBadge";
import { ConflictBadge } from "../common/ConflictBadge";
import "./ConflictPanel.css";

interface Props {
  conflicts: ConflictBooking[];
  onRebook: (conflict: ConflictBooking) => void;
  reboundBookingIds: number[];
}

// ConflictPanel lists every active booking whose time window overlaps
// another turnaround on the same resource, with a one-click rebook action.
export function ConflictPanel({ conflicts, onRebook, reboundBookingIds }: Props) {
  if (conflicts.length === 0) {
    return <p className="panel-ok">✓ 无预约时间冲突</p>;
  }
  return (
    <div className="conflict-list">
      {conflicts.map((c) => {
        const cleared = reboundBookingIds.includes(c.booking.id);
        return (
          <article key={c.booking.id} className={"conflict-card" + (cleared ? " cleared" : "")}>
            <div className="conflict-main">
              <div className="conflict-title">
                <ConflictBadge value={cleared ? "CLEARED" : "CONFLICT"} />
                <strong>预约 #{c.booking.id}</strong>
                <span className="dim">{formatTaskType(c.task.task_type)} · 任务 #{c.task.id}</span>
              </div>
              <p className="conflict-detail">{c.reason}</p>
              <div className="conflict-meta">
                <span>资源 <strong>{c.resource.resource_code}</strong>（{c.resource.location}）</span>
                <span>占用 {formatTime(c.booking.start_time)} ~ {formatTime(c.booking.end_time)}</span>
                <span>撞车预约 #{c.other_booking.id}（航班 #{c.other_booking.turnaround_id}，{formatTime(c.other_booking.start_time)}~{formatTime(c.other_booking.end_time)}）</span>
              </div>
              {cleared && <p className="cleared-note">该预约已换绑，冲突已解除（详情见换绑历史）</p>}
            </div>
            <div className="conflict-side">
              <div className="choices-preview">
                <StatusBadge value={c.available_choices.length > 0 ? "AVAILABLE" : "MAINTENANCE"} />
                <span>{c.available_choices.length} 个可换用资源</span>
              </div>
              <button className="btn warning" disabled={cleared} onClick={() => onRebook(c)}>
                换用可用资源
              </button>
            </div>
          </article>
        );
      })}
    </div>
  );
}
