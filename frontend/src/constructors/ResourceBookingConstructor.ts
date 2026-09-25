import type { ResourceBooking } from "../types/ResourceBooking";

// Default RELEASED-sentinel-free booking shape used by forms and optimistic
// rows; pages and stores must not hand-write this structure.
export const createDefaultResourceBooking = (overrides: Partial<ResourceBooking> = {}): ResourceBooking => ({
  id: 0,
  resource_id: 0,
  turnaround_id: 0,
  task_id: 0,
  start_time: "",
  end_time: "",
  booking_status: "CONFIRMED",
  conflict_reason: "",
  replaced_booking_id: null,
  replaced_by_id: null,
  ...overrides
});

export const createResourceBookingForm = createDefaultResourceBooking;
export const createResourceBookingResponse = createDefaultResourceBooking;
