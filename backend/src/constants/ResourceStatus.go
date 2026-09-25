package constants

var ResourceStatus = []string{"AVAILABLE", "BOOKED", "MAINTENANCE", "OFFLINE"}

const (
	ResourceStatusAvailable   = "AVAILABLE"
	ResourceStatusBooked      = "BOOKED"
	ResourceStatusMaintenance = "MAINTENANCE"
	ResourceStatusOffline     = "OFFLINE"
)

// UsableForRebook marks resource statuses that may receive a swapped booking.
var UsableForRebook = []string{ResourceStatusAvailable}
