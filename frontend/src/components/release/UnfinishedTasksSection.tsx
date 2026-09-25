import type { GroundTask } from "../../types/GroundTask";
import { formatTaskType, formatTaskStatus, formatTime } from "../../utils/formatters";
import { StatusBadge } from "../common/StatusBadge";
import "./BlockerSections.css";

interface Props {
  tasks: GroundTask[];
  completingId: number | null;
  onComplete: (id: number) => void;
}

// UnfinishedTasksSection lists tasks not in COMPLETED status with a quick
// sign-off action.
export function UnfinishedTasksSection({ tasks, completingId, onComplete }: Props) {
  if (tasks.length === 0) {
    return <p className="panel-ok">✓ 全部任务已签收完成</p>;
  }
  return (
    <ul className="blocker-list">
      {tasks.map((t) => (
        <li key={t.id} className="blocker-row">
          <div className="blocker-main">
            <div className="blocker-title">
              <strong>#{t.id} {formatTaskType(t.task_type)}</strong>
              <StatusBadge value={t.status} />
            </div>
            <div className="blocker-meta">
              <span className="dim">计划 {formatTime(t.planned_start)} ~ 截止 {formatTime(t.deadline)}</span>
              {t.blocker_note && <span className="blocker-note">阻塞原因：{t.blocker_note}</span>}
            </div>
          </div>
          <div className="blocker-side">
            <span className="dim">{formatTaskStatus(t.status)}</span>
            <button className="btn small" disabled={completingId === t.id} onClick={() => onComplete(t.id)}>
              {completingId === t.id ? "签收中…" : "签收完成"}
            </button>
          </div>
        </li>
      ))}
    </ul>
  );
}
