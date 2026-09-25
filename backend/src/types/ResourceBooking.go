package types

import (
	"time"

	"groundTurn/src/models"
)

// RebookRequest swaps an existing (non-released) booking onto NewResourceID.
type RebookRequest struct {
	NewResourceID int       `json:"new_resource_id" binding:"required"`
	StartTime     time.Time `json:"new_start_time"`
	EndTime       time.Time `json:"new_end_time"`
	Actor         string    `json:"actor"`
}

// RebookResult exposes the before/after booking pair for history display.
type RebookResult struct {
	ReleasedBooking models.ResourceBooking `json:"released_booking"`
	NewBooking      models.ResourceBooking `json:"new_booking"`
	OldResource     models.GroundResource  `json:"old_resource"`
	NewResource     models.GroundResource  `json:"new_resource"`
}

// BookingHistoryItem decorates a booking row with related entity names.
type BookingHistoryItem struct {
	Booking  models.ResourceBooking  `json:"booking"`
	Resource models.GroundResource   `json:"resource"`
	Task     models.GroundTask       `json:"task"`
	Flight   models.FlightTurnaround `json:"flight"`
}
