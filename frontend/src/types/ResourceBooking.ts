export interface ResourceBooking {
  id: number;
  resource_id: number;
  turnaround_id: number;
  task_id: number;
  start_time: string;
  end_time: string;
  booking_status: string;
  conflict_reason: string;
  replaced_booking_id: number | null;
  replaced_by_id: number | null;
  created_at?: string;
}

export interface RebookResult {
  released_booking: ResourceBooking;
  new_booking: ResourceBooking;
  old_resource: { id: number; resource_code: string; resource_type: string };
  new_resource: GroundResourceLite;
}

export interface BookingHistoryItem {
  booking: ResourceBooking;
  resource: GroundResourceLite;
  task: { id: number; task_type: string; status: string };
  flight: { id: number; flight_no: string };
}

export interface RebookCandidate {
  resource: GroundResourceLite;
  usable: boolean;
  reason?: string;
}

interface GroundResourceLite {
  id: number;
  resource_code: string;
  resource_type: string;
  location: string;
  availability_status: string;
  owner_team: string;
}
