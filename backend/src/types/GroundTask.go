package types

import "groundTurn/src/models"

type TaskStatusRequest struct {
	Status      string `json:"status"`
	BlockerNote string `json:"blocker_note"`
	Actor       string `json:"actor"`
}

// OpenTask is a task plus its active (non-released) booking, if any.
type OpenTask struct {
	Task    models.GroundTask       `json:"task"`
	Booking *models.ResourceBooking `json:"booking,omitempty"`
}
