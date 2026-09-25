export interface FlightTurnaround {
  id: number;
  flight_no: string;
  aircraft_reg: string;
  stand_no: string;
  arrival_time: string;
  departure_time: string;
  turnaround_status: string;
  delay_reason: string;
}

// ---- Release coordination payloads ----

export type ViolationKind = "task" | "delay" | "booking" | "turnaround";

export interface ReleaseViolation {
  code: string;
  kind: ViolationKind;
  task_id?: number;
  booking_id?: number;
  resource_id?: number;
  delay_id?: number;
  detail: string;
}

export interface ConflictBooking {
  booking: import("./ResourceBooking").ResourceBooking;
  resource: import("./GroundResource").GroundResource;
  task: import("./GroundTask").GroundTask;
  other_booking: import("./ResourceBooking").ResourceBooking;
  reason: string;
  available_choices: import("./GroundResource").GroundResource[];
}

export interface ReleaseSummary {
  total_tasks: number;
  completed_tasks: number;
  open_delays: number;
  booking_conflicts: number;
}

export interface ReleaseCheckResult {
  turnaround: FlightTurnaround;
  can_release: boolean;
  summary: ReleaseSummary;
  unfinished_tasks: import("./GroundTask").GroundTask[];
  open_delays: import("./DelayEvent").DelayEvent[];
  conflicts: ConflictBooking[];
  violations: ReleaseViolation[];
  checked_at: string;
}

export interface ReleaseResponse {
  turnaround: FlightTurnaround;
  released_at: string;
  message: string;
}

export interface ApiErrorBody {
  ok: false;
  code: string;
  message: string;
  violations?: ReleaseViolation[];
  violation?: ReleaseViolation;
}
