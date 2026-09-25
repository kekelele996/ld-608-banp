import { useEffect, useState } from "react";
import { useFlightTurnaroundStore } from "../stores/FlightTurnaroundStore";
import { useTurnaroundReleaseStore } from "../stores/TurnaroundReleaseStore";
import { StatusBadge } from "../components/common/StatusBadge";
import { StatCard } from "../components/common/StatCard";
import { DelayTag } from "../components/common/DelayTag";
import { ConflictBadge } from "../components/common/ConflictBadge";
import { EmptyState } from "../components/common/EmptyState";
import { ReleaseConclusionText, ReleaseBlockerKindText } from "../constants/ReleaseConclusion";
import { GroundTaskStatusText } from "../constants/GroundTaskStatus";
import { BookingStatusText } from "../constants/BookingStatus";
import { formatDate, formatWindow } from "../utils/formatters";
import type { BookingConflict } from "../types/TurnaroundRelease";

export function ReleasePage() {
  const { rows: flights, load: loadFlights } = useFlightTurnaroundStore();
  const { summary, history, rebindOptions, loading, error, notice, loadSummary, release, loadRebindOptions, rebind, clearFeedback } =
    useTurnaroundReleaseStore();
  const [selectedId, setSelectedId] = useState<number>(0);
  const [rebinding, setRebinding] = useState<BookingConflict | null>(null);

  useEffect(() => { loadFlights(); }, [loadFlights]);
  useEffect(() => {
    if (flights.length > 0 && selectedId === 0) setSelectedId(flights[0].id);
  }, [flights, selectedId]);
  useEffect(() => {
    if (selectedId > 0) {
      clearFeedback();
      setRebinding(null);
      loadSummary(selectedId);
    }
  }, [selectedId, loadSummary, clearFeedback]);

  const openRebind = (conflict: BookingConflict) => {
    setRebinding(conflict);
    loadRebindOptions(conflict.booking.id);
  };

  const confirmRebind = async (resourceId: number) => {
    if (!rebinding) return;
    const ok = await rebind(rebinding.booking.id, resourceId, selectedId);
    if (ok) setRebinding(null);
  };

  const submitRelease = async (turnaroundId: number) => {
    const ok = await release(turnaroundId);
    if (ok) loadFlights();
  };

  const t = summary?.turnaround;

  return <main className="page">
    <section className="page-head">
      <div>
        <p className="eyebrow">ground-turn</p>
        <h1>过站放行协同</h1>
      </div>
      {t && <StatusBadge value={t.turnaround_status} />}
    </section>

    <section className="release-layout">
      <div className="panel flight-list">
        <h2>在港航班</h2>
        {flights.map((f) => (
          <button key={f.id} className={"flight-item" + (f.id === selectedId ? " active" : "")} onClick={() => setSelectedId(f.id)}>
            <strong>{f.flight_no}</strong>
            <span>机位 {f.stand_no} · {f.aircraft_reg}</span>
            <StatusBadge value={f.turnaround_status} />
          </button>
        ))}
        {flights.length === 0 && <EmptyState title="暂无在港航班" />}
      </div>

      <div className="release-detail">
        {error && <section className="feedback error">
          <strong>{error.message}</strong>
          {error.blockers && error.blockers.length > 0 && <ul>
            {error.blockers.map((b, i) => <li key={i}><em>[{ReleaseBlockerKindText[b.kind] ?? b.kind}]</em> {b.message}</li>)}
          </ul>}
        </section>}
        {notice && <section className="feedback ok">{notice}</section>}

        {summary && t && <>
          <section className={"conclusion " + summary.conclusion.toLowerCase()}>
            <div>
              <span className="eyebrow">放行结论</span>
              <strong>{ReleaseConclusionText[summary.conclusion]}</strong>
            </div>
            <button className="primary" disabled={loading} onClick={() => submitRelease(t.id)}>
              提交放行
            </button>
          </section>

          <section className="metrics">
            <StatCard label="未完成任务" value={summary.unfinished_tasks.length} />
            <StatCard label="未关闭延误" value={summary.open_delays.length} />
            <StatCard label="预约冲突" value={summary.booking_conflicts.length} />
          </section>

          <section className="panel">
            <h2>未完成任务</h2>
            {summary.unfinished_tasks.length === 0 && <EmptyState title="全部任务已完成" />}
            {summary.unfinished_tasks.map((task) => (
              <article key={task.id} className="row">
                <strong>#{task.id} {task.task_type}</strong>
                <span>截止 {formatDate(task.deadline)}{task.blocker_note ? ` · ${task.blocker_note}` : ""}</span>
                <StatusBadge value={GroundTaskStatusText[task.status as keyof typeof GroundTaskStatusText] ?? task.status} />
              </article>
            ))}
          </section>

          <section className="panel">
            <h2>未关闭延误</h2>
            {summary.open_delays.length === 0 && <EmptyState title="无未关闭延误" />}
            {summary.open_delays.map((d) => (
              <article key={d.id} className="row">
                <DelayTag title={`#${d.id} ${d.delay_type}`} value={`${d.minutes} 分钟`} />
                <span>{d.root_cause} · 责任：{d.responsibility_team}</span>
              </article>
            ))}
          </section>

          <section className="panel">
            <h2>预约时间冲突</h2>
            {summary.booking_conflicts.length === 0 && <EmptyState title="预约无冲突" />}
            {summary.booking_conflicts.map((c) => (
              <article key={c.booking.id} className="row conflict-row">
                <ConflictBadge title={`${c.resource.resource_code} · 预约 #${c.booking.id}`} value={formatWindow(c.booking.start_time, c.booking.end_time)} />
                <span>与预约 #{c.conflicts_with.id}（航班 #{c.conflicts_with.turnaround_id}）时段重叠</span>
                <button onClick={() => openRebind(c)}>换用可用资源</button>
              </article>
            ))}
            {rebinding && <div className="rebind-box">
              <h3>预约 #{rebinding.booking.id} 换绑（{formatWindow(rebinding.booking.start_time, rebinding.booking.end_time)}）</h3>
              {rebindOptions.length === 0 && <EmptyState title="该时段无可用同类资源" />}
              {rebindOptions.map((r) => (
                <button key={r.id} className="rebind-option" disabled={loading} onClick={() => confirmRebind(r.id)}>
                  <strong>{r.resource_code}</strong>
                  <span>{r.location} · {r.owner_team}</span>
                </button>
              ))}
              <button className="link" onClick={() => setRebinding(null)}>取消</button>
            </div>}
          </section>

          <section className="panel">
            <h2>预约历史（含已释放）</h2>
            {history.length === 0 && <EmptyState title="暂无预约记录" />}
            {history.map((b) => (
              <article key={b.id} className="row">
                <strong>预约 #{b.id} · 资源 #{b.resource_id}</strong>
                <span>{formatWindow(b.start_time, b.end_time)}{b.conflict_reason ? ` · ${b.conflict_reason}` : ""}</span>
                <StatusBadge value={BookingStatusText[b.booking_status as keyof typeof BookingStatusText] ?? b.booking_status} />
              </article>
            ))}
          </section>
        </>}
        {!summary && !loading && <EmptyState title="请选择左侧航班查看放行详情" />}
        {loading && <EmptyState title="加载中…" />}
      </div>
    </section>
  </main>;
}
