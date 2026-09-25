import type { ReleaseCheckResult, ReleaseViolation } from "../../types/FlightTurnaround";
import { ApiError } from "../../api/client";
import "./ReleaseVerdict.css";

interface Props {
  check: ReleaseCheckResult | null;
  submitting: boolean;
  rejectError: ApiError | null;
  onSubmit: () => void;
}

// ReleaseVerdict renders the go/no-go conclusion. A server rejection (409)
// surfaces the exact violations returned by the backend re-check — the UI
// never decides release on its own.
export function ReleaseVerdict({ check, submitting, rejectError, onSubmit }: Props) {
  if (!check) return null;

  const serverViolations: ReleaseViolation[] = rejectError?.violations ?? [];

  return (
    <section className={"verdict " + (check.can_release ? "pass" : "block")}>
      <div className="verdict-head">
        <span className="verdict-stamp">{check.can_release ? "可放行" : "不可放行"}</span>
        <p className="verdict-summary">
          {check.can_release
            ? "三类条件均已满足：任务全部完成、延误已关闭、预约无冲突，可以提交放行。"
            : `共 ${check.violations.length} 项阻塞条件未解除，提交后后端将驳回且航班状态与原预约保持不变。`}
        </p>
      </div>

      {serverViolations.length > 0 && (
        <div className="verdict-reject">
          <strong>后端核对未通过：{rejectError?.message}</strong>
          <ul>
            {serverViolations.map((v, i) => (
              <li key={i} className={"violation kind-" + v.kind}>
                <span className="violation-kind">{kindLabel(v.kind)}</span>
                {v.detail}
              </li>
            ))}
          </ul>
        </div>
      )}

      <div className="verdict-actions">
        <button className="btn primary lg" disabled={!check.can_release || submitting} onClick={onSubmit}>
          {submitting ? "后端核对中…" : check.can_release ? "提交放行" : "条件不足，禁止放行"}
        </button>
        <span className="dim">核对时间 {new Date(check.checked_at).toLocaleTimeString("zh-CN", { hour12: false })}</span>
      </div>
    </section>
  );
}

function kindLabel(kind: string) {
  return { task: "任务", delay: "延误", booking: "预约", turnaround: "航班" }[kind] ?? kind;
}
