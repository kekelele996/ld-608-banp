import { useCallback, useEffect, useState } from "react";
import { listRebookCandidates, rebookBooking } from "../api/ResourceBooking";
import type { RebookCandidate } from "../types/ResourceBooking";

// useResourceConflict loads the swap candidates for a conflicting booking and
// performs the rebook mutation, exposing transient loading/error state for
// the swap modal.
export function useResourceConflict(bookingId: number | null) {
  const [candidates, setCandidates] = useState<RebookCandidate[]>([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const load = useCallback((id: number) => {
    setLoading(true);
    setError(null);
    listRebookCandidates(id)
      .then(setCandidates)
      .catch((e: Error) => setError(e.message))
      .finally(() => setLoading(false));
  }, []);

  useEffect(() => {
    if (bookingId == null) {
      setCandidates([]);
      return;
    }
    load(bookingId);
  }, [bookingId, load]);

  const swap = useCallback(
    async (newResourceId: number) => {
      if (bookingId == null) throw new Error("未指定预约");
      return rebookBooking(bookingId, { new_resource_id: newResourceId });
    },
    [bookingId]
  );

  return { candidates, loading, error, swap, reload: () => bookingId != null && load(bookingId) };
}
