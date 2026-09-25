import { useEffect, useState } from "react";
import { useResourceConflict } from "../../hooks/useResourceConflict";
import { listBookingHistory } from "../../api/ResourceBooking";
import type { BookingHistoryItem } from "../../types/ResourceBooking";
import type { ConflictBooking } from "../../types/FlightTurnaround";
import { formatBookingStatus, formatDateTime, formatTime } from "../../utils/formatters";
import { StatusBadge } from "../common/StatusBadge";
import "./RebookModal.css";

interface Props {
  conflict: ConflictBooking;
  onClose: () => void;
  onRebound: () => void;
}

// RebookModal lets the dispatcher move a conflicting booking onto an
// AVAILABLE, time-window-free resource. After a successful swap it shows the
// before (RELEASED) / after (CONFIRMED) booking pair queried by task id.
export function RebookModal({ conflict, onClose, onRebound }: Props) {
  const { candidates, loading, error, swap } = useResourceConflict(conflict.booking.id);
  const [selected, setSelected] = useState<number | null>(null);
  const [submitting, setSubmitting] = useState(false);
  const [submitError, setSubmitError] = useState<string | null>(null);
  const [history, setHistory] = useState<BookingHistoryItem[] | null>(null);

  useEffect(() => {
    listBookingHistory(conflict.task.id).then(setHistory).catch(() => undefined);
  }, [conflict.task.id]);

  const usable = candidates.filter((c) => c.usable);

  const confirm = async () => {
    if (selected == null) return;
    setSubmitting(true);
    setSubmitError(null);
    try {
      await swap(selected);
      const rows = await listBookingHistory(conflict.task.id);
      setHistory(rows);
      onRebound();
    } catch (e) {
      setSubmitError((e as Error).message);
    } finally {
      setSubmitting(false);
    }
  };

  return (
    <div className="modal-mask" onClick={onClose}>
      <div className="modal" onClick={(e) => e.stopPropagation()}>
        <header className="modal-head">
          <h3>冲突任务换用资源</h3>
          <button className="modal-x" onClick={onClose} aria-label="关闭">×</button>
        </header>

        <div className="modal-current">
          <div>
            <span className="dim">冲突预约</span>
            <strong>#{conflict.booking.id}</strong>
            <StatusBadge value={conflict.booking.booking_status} />
          </div>
          <div>
            <span className="dim">现用资源</span>
            <strong>{conflict.resource.resource_code}</strong>
            <span className="dim">{conflict.resource.location}</span>
          </div>
          <div>
            <span className="dim">时间窗</span>
            <strong>{formatTime(conflict.booking.start_time)} ~ {formatTime(conflict.booking.end_time)}</strong>
          </div>
        </div>
        <p className="modal-reason">{conflict.reason}</p>

        <h4 className="modal-section">可用资源（同类型且时间窗空闲）</h4>
        {loading && <p className="dim">加载候选资源…</p>}
        {!loading && usable.length === 0 && (
          <p className="modal-empty">当前没有可换用的资源，请先释放占用或维护中的同类型资源。</p>
        )}
        <div className="candidate-grid">
          {candidates.map((c) => (
            <button
              key={c.resource.id}
              disabled={!c.usable}
              className={"candidate " + (selected === c.resource.id ? "picked" : "") + (c.usable ? "" : " disabled")}
              onClick={() => setSelected(c.resource.id)}
              title={c.reason}
            >
              <strong>{c.resource.resource_code}</strong>
              <span>{c.resource.location}</span>
              <StatusBadge value={c.resource.availability_status} />
              {!c.usable && <em className="candidate-reason">{c.reason}</em>}
            </button>
          ))}
        </div>

        {submitError && <p className="modal-error">换绑失败：{submitError}</p>}
        {error && <p className="modal-error">候选加载失败：{error}</p>}

        <h4 className="modal-section">预约历史（前后两次）</h4>
        <ul className="history-list">
          {(history ?? []).map((h) => (
            <li key={h.booking.id} className={h.booking.booking_status === "RELEASED" ? "history-old" : "history-new"}>
              <StatusBadge value={h.booking.booking_status} />
              <span className="history-id">#{h.booking.id}</span>
              <strong>{h.resource.resource_code}</strong>
              <span className="dim">{formatTime(h.booking.start_time)}~{formatTime(h.booking.end_time)}</span>
              <span className="dim">{formatDateTime(h.booking.created_at)}</span>
              {h.booking.booking_status === "RELEASED" && h.booking.replaced_by_id != null && (
                <em className="chain">→ 新预约 #{h.booking.replaced_by_id}</em>
              )}
              {h.booking.booking_status !== "RELEASED" && h.booking.replaced_booking_id != null && (
                <em className="chain">← 原预约 #{h.booking.replaced_booking_id}</em>
              )}
            </li>
          ))}
        </ul>

        <footer className="modal-foot">
          <button className="btn ghost" onClick={onClose}>取消</button>
          <button className="btn primary" disabled={selected == null || submitting} onClick={confirm}>
            {submitting ? "换绑中…" : "确认换用该资源"}
          </button>
        </footer>
      </div>
    </div>
  );
}
