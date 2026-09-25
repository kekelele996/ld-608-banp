import { create } from "zustand";
import { fetchReleaseSummary, submitRelease } from "../api/TurnaroundRelease";
import { fetchRebindOptions, listResourceBookingByTurnaround, rebindResourceBooking } from "../api/ResourceBooking";
import { LOG_TEMPLATES } from "../constants/logTemplates";
import type { ApiError, ReleaseSummary } from "../types/TurnaroundRelease";
import type { ResourceBooking } from "../types/ResourceBooking";
import type { GroundResource } from "../types/GroundResource";

type State = {
  summary: ReleaseSummary | null;
  history: ResourceBooking[];
  rebindOptions: GroundResource[];
  loading: boolean;
  error: ApiError | null;
  notice: string;
  loadSummary: (turnaroundId: number) => Promise<void>;
  release: (turnaroundId: number) => Promise<boolean>;
  loadRebindOptions: (bookingId: number) => Promise<void>;
  rebind: (bookingId: number, resourceId: number, turnaroundId: number) => Promise<boolean>;
  clearFeedback: () => void;
};

export const useTurnaroundReleaseStore = create<State>((set, get) => ({
  summary: null,
  history: [],
  rebindOptions: [],
  loading: false,
  error: null,
  notice: "",
  async loadSummary(turnaroundId) {
    set({ loading: true, error: null });
    try {
      const [summary, history] = await Promise.all([
        fetchReleaseSummary(turnaroundId),
        listResourceBookingByTurnaround(turnaroundId)
      ]);
      set({ summary, history, loading: false });
    } catch (err) {
      set({ error: err as ApiError, summary: null, history: [], loading: false });
    }
  },
  async release(turnaroundId) {
    set({ loading: true, error: null, notice: "" });
    try {
      const result = await submitRelease(turnaroundId);
      console.info(LOG_TEMPLATES.FlightTurnaround[4], result.turnaround.flight_no);
      set({ notice: `航班 ${result.turnaround.flight_no} 放行成功，状态已置为 ${result.turnaround.turnaround_status}`, loading: false });
      await get().loadSummary(turnaroundId);
      return true;
    } catch (err) {
      // 复核未通过：后端已指出具体航班、任务或资源，航班状态与预约保持不变。
      console.warn(LOG_TEMPLATES.FlightTurnaround[5], turnaroundId);
      set({ error: err as ApiError, loading: false });
      return false;
    }
  },
  async loadRebindOptions(bookingId) {
    try {
      set({ rebindOptions: await fetchRebindOptions(bookingId) });
    } catch (err) {
      set({ error: err as ApiError, rebindOptions: [] });
    }
  },
  async rebind(bookingId, resourceId, turnaroundId) {
    set({ loading: true, error: null, notice: "" });
    try {
      const result = await rebindResourceBooking(bookingId, resourceId);
      console.info(LOG_TEMPLATES.ResourceBooking[4], result.released.id);
      console.info(LOG_TEMPLATES.ResourceBooking[5], result.created.id);
      set({ notice: `预约 #${result.released.id} 已释放，新预约 #${result.created.id} 已生效`, loading: false });
      await get().loadSummary(turnaroundId);
      return true;
    } catch (err) {
      set({ error: err as ApiError, loading: false });
      return false;
    }
  },
  clearFeedback() {
    set({ error: null, notice: "" });
  }
}));
