package constants

// Error codes are referenced by services, controllers and the release-check
// response payloads; keep one constant per failure so callers can branch.
const (
	AuthRequired = "AUTH_REQUIRED"

	TurnaroundNotFound      = "TURNAROUND_NOT_FOUND"
	TurnaroundNotReleasable = "TURNAROUND_NOT_RELEASABLE"

	TaskUnfinished       = "TASK_UNFINISHED"
	TaskNotFound         = "TASK_NOT_FOUND"
	DelayNotClosed       = "DELAY_NOT_CLOSED"
	BookingConflict      = "BOOKING_CONFLICT"
	BookingNotFound      = "BOOKING_NOT_FOUND"
	BookingNotRebindable = "BOOKING_NOT_REBINDABLE"
	ResourceUnavailable  = "RESOURCE_UNAVAILABLE"
	ResourceNotFound     = "RESOURCE_NOT_FOUND"
	ResourceTypeMismatch = "RESOURCE_TYPE_MISMATCH"
	ResourceBusy         = "RESOURCE_BUSY"

	ReleaseRejected = "RELEASE_REJECTED"
)
