import type { FlightTurnaround } from "./FlightTurnaround";
import type { GroundTask } from "./GroundTask";
import type { GroundResource } from "./GroundResource";
import type { ResourceBooking } from "./ResourceBooking";
import type { DelayEvent } from "./DelayEvent";

export interface ReleaseBlocker {
  kind: "TASK" | "DELAY" | "BOOKING";
  ref_id: number;
  message: string;
}

export interface BookingConflict {
  booking: ResourceBooking;
  conflicts_with: ResourceBooking;
  resource: GroundResource;
  reason: string;
}

export interface ReleaseSummary {
  turnaround: FlightTurnaround;
  unfinished_tasks: GroundTask[];
  open_delays: DelayEvent[];
  booking_conflicts: BookingConflict[];
  conclusion: "RELEASABLE" | "BLOCKED";
  blockers: ReleaseBlocker[];
}

export interface ReleaseResult {
  ok: boolean;
  turnaround: FlightTurnaround;
}

export interface RebindResult {
  released: ResourceBooking;
  created: ResourceBooking;
}

export interface ApiError {
  code: string;
  message: string;
  blockers?: ReleaseBlocker[];
}
