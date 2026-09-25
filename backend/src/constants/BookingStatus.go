package constants

// BookingStatus enumerates resource booking lifecycle values.
// RELEASED bookings are kept forever for history lookup (old -> new chain).
var BookingStatus = []string{"CONFIRMED", "CONFLICT", "RELEASED"}

const (
	BookingStatusConfirmed = "CONFIRMED"
	BookingStatusConflict  = "CONFLICT"
	BookingStatusReleased  = "RELEASED"
)
