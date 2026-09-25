import { create } from "zustand";
import { listTurnarounds, getReleaseCheck } from "../api/FlightTurnaround";
import type { FlightTurnaround, ReleaseCheckResult } from "../types/FlightTurnaround";

// ReleaseStore backs the release workbench. The page calls load() once and
// refresh(id) after every mutation (sign-off, delay close, rebook).
type State = {
  flights: FlightTurnaround[];
  selectedId: number | null;
  check: ReleaseCheckResult | null;
  loading: boolean;
  load: () => Promise<void>;
  select: (id: number) => Promise<void>;
  refresh: () => Promise<void>;
};

export const useReleaseStore = create<State>((set, get) => ({
  flights: [],
  selectedId: null,
  check: null,
  loading: false,
  async load() {
    const flights = await listTurnarounds();
    const first = flights.find((f) => f.turnaround_status !== "READY" && f.turnaround_status !== "DEPARTED") ?? flights[0];
    set({ flights, selectedId: first?.id ?? null });
    if (first) {
      set({ loading: true });
      const check = await getReleaseCheck(first.id);
      set({ check, loading: false });
    }
  },
  async select(id) {
    set({ selectedId: id, loading: true, check: null });
    const check = await getReleaseCheck(id);
    set({ check, loading: false });
  },
  async refresh() {
    const id = get().selectedId;
    if (id == null) return;
    set({ loading: true });
    const check = await getReleaseCheck(id);
    const flights = await listTurnarounds();
    set({ check, flights, loading: false });
  }
}));
