import type { DelayEvent } from "../../types/DelayEvent";
import { DelayTag } from "../common/DelayTag";
import "./BlockerSections.css";

interface Props {
  delays: DelayEvent[];
  resolvingId: number | null;
  onResolve: (id: number) => void;
}

// OpenDelaysSection lists delays without resolved_at with a close action.
export function OpenDelaysSection({ delays, resolvingId, onResolve }: Props) {
  if (delays.length === 0) {
    return <p className="panel-ok">✓ 无未关闭延误</p>;
  }
  return (
    <ul className="blocker-list">
      {delays.map((d) => (
        <li key={d.id} className="blocker-row">
          <div className="blocker-main">
            <div className="blocker-title">
              <DelayTag title="" value="DELAYED" />
              <strong>延误 #{d.id}</strong>
              <span className="delay-minutes">+{d.minutes} 分钟</span>
            </div>
            <div className="blocker-meta">
              <span className="dim">类型 {d.delay_type}</span>
              <span className="dim">责任班组 {d.responsibility_team}</span>
              {d.root_cause && <span className="blocker-note">根因：{d.root_cause}</span>}
            </div>
          </div>
          <div className="blocker-side">
            <button className="btn small" disabled={resolvingId === d.id} onClick={() => onResolve(d.id)}>
              {resolvingId === d.id ? "关闭中…" : "关闭延误"}
            </button>
          </div>
        </li>
      ))}
    </ul>
  );
}
