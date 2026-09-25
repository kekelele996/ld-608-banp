import { useMemo } from "react";
import type { ReleaseCheckResult } from "../types/FlightTurnaround";

// useTurnaroundProgress derives completion percentage and blocker counters
// from a release-check aggregation.
export function useTurnaroundProgress(check: ReleaseCheckResult | null) {
  return useMemo(() => {
    const total = check?.summary.total_tasks ?? 0;
    const completed = check?.summary.completed_tasks ?? 0;
    const percent = total === 0 ? 0 : Math.round((completed / total) * 100);
    const blockers =
      (check?.unfinished_tasks.length ?? 0) +
      (check?.open_delays.length ?? 0) +
      (check?.summary.booking_conflicts ?? 0);
    return {
      total,
      completed,
      percent,
      blockers,
      unfinished: check?.unfinished_tasks.length ?? 0,
      openDelays: check?.open_delays.length ?? 0,
      conflicts: check?.summary.booking_conflicts ?? 0
    };
  }, [check]);
}
