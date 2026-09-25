import { useCallback, useEffect, useMemo, useState } from "react";
import {
  listTurnarounds,
  getReleaseCheck,
  submitRelease
} from "../api/FlightTurnaround";
import { updateTaskStatus } from "../api/GroundTask";
import { resolveDelay } from "../api/DelayEvent";
import { ApiError } from "../api/client";
import type { FlightTurnaround, ReleaseCheckResult, ConflictBooking } from "../types/FlightTurnaround";
import { useTurnaroundProgress } from "../hooks/useTurnaroundProgress";
import { formatTurnaroundStatus, formatDateTime, formatTime } from "../utils/formatters";
import { StatusBadge } from "../components/common/StatusBadge";
import { ReleaseVerdict } from "../components/release/ReleaseVerdict";
import { UnfinishedTasksSection } from "../components/release/UnfinishedTasksSection";
import { OpenDelaysSection } from "../components/release/OpenDelaysSection";
import { ConflictPanel } from "../components/release/ConflictPanel";
import { RebookModal } from "../components/release/RebookModal";
import "./ReleasePage.css";

// ReleasePage is the dispatcher's turnaround release workbench: pick a flight,
// review the three blocker categories with a verdict, swap conflicting
// resources, and submit — the backend re-verifies before flipping status.
export function ReleasePage() {
  const [flights, setFlights] = useState<FlightTurnaround[]>([]);
  const [selectedId, setSelectedId] = useState<number | null>(null);
  const [check, setCheck] = useState<ReleaseCheckResult | null>(null);
  const [loading, setLoading] = useState(false);
  const [submitting, setSubmitting] = useState(false);
  const [rejectError, setRejectError] = useState<ApiError | null>(null);
  const [toast, setToast] = useState<string | null>(null);
  const [rebookConflict, setRebookConflict] = useState<ConflictBooking | null>(null);
  // Bookings swapped in this session stay marked cleared until refresh.
  const [reboundIds, setReboundIds] = useState<number[]>([]);
  const [completingId, setCompletingId] = useState<number | null>(null);
  const [resolvingId, setResolvingId] = useState<number | null>(null);

  useEffect(() => {
    listTurnarounds()
      .then((rows) => {
        setFlights(rows);
        const first = rows.find((f) => f.turnaround_status !== "READY" && f.turnaround_status !== "DEPARTED");
        setSelectedId((first ?? rows[0])?.id ?? null);
      })
      .catch((e: Error) => setToast("航班列表加载失败：" + e.message));
  }, []);

  const refreshCheck = useCallback((id: number) => {
    setLoading(true);
    setRejectError(null);
    return getReleaseCheck(id)
      .then(setCheck)
      .catch((e: Error) => setToast("放行汇总加载失败：" + e.message))
      .finally(() => setLoading(false));
  }, []);

  useEffect(() => {
    if (selectedId != null) {
      setReboundIds([]);
      void refreshCheck(selectedId);
    }
  }, [selectedId, refreshCheck]);

  const progress = useTurnaroundProgress(check);

  const handleComplete = async (taskId: number) => {
    if (selectedId == null) return;
    setCompletingId(taskId);
    try {
      await updateTaskStatus(taskId, "COMPLETED");
      await refreshCheck(selectedId);
      flash(`任务 #${taskId} 已签收完成`);
    } catch (e) {
      setToast((e as Error).message);
    } finally {
      setCompletingId(null);
    }
  };

  const handleResolveDelay = async (delayId: number) => {
    if (selectedId == null) return;
    setResolvingId(delayId);
    try {
      await resolveDelay(delayId);
      await refreshCheck(selectedId);
      flash(`延误 #${delayId} 已关闭`);
    } catch (e) {
      setToast((e as Error).message);
    } finally {
      setResolvingId(null);
    }
  };

  const handleSubmit = async () => {
    if (selectedId == null) return;
    setSubmitting(true);
    setRejectError(null);
    try {
      const res = await submitRelease(selectedId);
      flash(res.message);
      setFlights((rows) => rows.map((f) => (f.id === selectedId ? res.turnaround : f)));
      await refreshCheck(selectedId);
    } catch (e) {
      if (e instanceof ApiError) {
        setRejectError(e);
        // Refetch the authoritative state (nothing changed server-side).
        await refreshCheck(selectedId);
      } else {
        setToast((e as Error).message);
      }
    } finally {
      setSubmitting(false);
    }
  };

  const flash = (msg: string) => {
    setToast(msg);
    window.setTimeout(() => setToast(null), 3200);
  };

  const visibleConflicts = useMemo(() => check?.conflicts ?? [], [check]);

  return (
    <main className="release-page">
      <header className="release-head">
        <div>
          <p className="eyebrow">ground-turn · 过站放行协同</p>
          <h1>过站放行</h1>
          <p className="subtitle">汇总未完成任务、未关闭延误与预约时间冲突，条件齐备方可放行。</p>
        </div>
        <button className="btn ghost" disabled={selectedId == null || loading} onClick={() => selectedId != null && refreshCheck(selectedId)}>
          {loading ? "核对中…" : "重新核对"}
        </button>
      </header>

      <div className="release-layout">
        <aside className="flight-list panel">
          <h2>航班列表</h2>
          {flights.map((f) => (
            <button
              key={f.id}
              className={"flight-item" + (f.id === selectedId ? " active" : "")}
              onClick={() => setSelectedId(f.id)}
            >
              <div className="flight-item-top">
                <strong>{f.flight_no}</strong>
                <StatusBadge value={f.turnaround_status} />
              </div>
              <span className="dim">{f.aircraft_reg} · 机位 {f.stand_no}</span>
              <span className="dim">{formatTime(f.arrival_time)} → {formatTime(f.departure_time)}</span>
            </button>
          ))}
        </aside>

        <section className="release-detail">
          {check && (
            <>
              <div className="panel flight-summary">
                <div>
                  <h2>{check.turnaround.flight_no}</h2>
                  <span className="dim">{check.turnaround.aircraft_reg} · 机位 {check.turnaround.stand_no} · {formatDateTime(check.turnaround.arrival_time)} 到港</span>
                </div>
                <div className="flight-status">
                  <StatusBadge value={check.turnaround.turnaround_status} />
                  <span className="dim">{formatTurnaroundStatus(check.turnaround.turnaround_status)}</span>
                </div>
                <div className="progress-block">
                  <div className="progress-bar"><i style={{ width: `${progress.percent}%` }} /></div>
                  <span className="dim">任务完成 {progress.completed}/{progress.total}（{progress.percent}%）</span>
                </div>
                <div className="counter-row">
                  <Counter label="未完成任务" value={progress.unfinished} danger={progress.unfinished > 0} />
                  <Counter label="未关闭延误" value={progress.openDelays} danger={progress.openDelays > 0} />
                  <Counter label="预约冲突" value={progress.conflicts} danger={progress.conflicts > 0} />
                </div>
              </div>

              <ReleaseVerdict check={check} submitting={submitting} rejectError={rejectError} onSubmit={handleSubmit} />

              <div className="blocker-grid">
                <section className="panel blocker-panel">
                  <h3>① 未完成任务 <em className="count">{check.unfinished_tasks.length}</em></h3>
                  <UnfinishedTasksSection tasks={check.unfinished_tasks} completingId={completingId} onComplete={handleComplete} />
                </section>

                <section className="panel blocker-panel">
                  <h3>② 未关闭延误 <em className="count">{check.open_delays.length}</em></h3>
                  <OpenDelaysSection delays={check.open_delays} resolvingId={resolvingId} onResolve={handleResolveDelay} />
                </section>
              </div>

              <section className="panel blocker-panel">
                <h3>③ 预约时间冲突 <em className="count">{check.conflicts.length}</em></h3>
                <ConflictPanel
                  conflicts={visibleConflicts}
                  reboundBookingIds={reboundIds}
                  onRebook={setRebookConflict}
                />
              </section>
            </>
          )}
          {!check && loading && <div className="panel">正在汇总航班放行条件…</div>}
        </section>
      </div>

      {rebookConflict && (
        <RebookModal
          conflict={rebookConflict}
          onClose={() => setRebookConflict(null)}
          onRebound={() => {
            setReboundIds((ids) => (ids.includes(rebookConflict.booking.id) ? ids : [...ids, rebookConflict.booking.id]));
            if (selectedId != null) void refreshCheck(selectedId);
          }}
        />
      )}

      {toast && <div className="toast">{toast}</div>}
    </main>
  );
}

function Counter({ label, value, danger }: { label: string; value: number; danger: boolean }) {
  return (
    <div className={"counter " + (danger ? "danger" : "ok")}>
      <strong>{value}</strong>
      <span>{label}</span>
    </div>
  );
}
