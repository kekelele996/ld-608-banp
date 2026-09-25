package types

import (
	"time"

	"groundTurn/src/models"
)

// ReleaseViolation pinpoints exactly which flight/task/resource/delay blocks
// a release. Code comes from constants/errorCodes and Detail is the rendered
// constants/errorMessages text so the UI can show it verbatim.
type ReleaseViolation struct {
	Code       string `json:"code"`
	Kind       string `json:"kind"` // task | delay | booking | turnaround
	TaskID     int    `json:"task_id,omitempty"`
	BookingID  int    `json:"booking_id,omitempty"`
	ResourceID int    `json:"resource_id,omitempty"`
	DelayID    int    `json:"delay_id,omitempty"`
	Detail     string `json:"detail"`
}

// ConflictBooking pairs two overlapping bookings on the same resource.
type ConflictBooking struct {
	Booking          models.ResourceBooking  `json:"booking"`
	Resource         models.GroundResource   `json:"resource"`
	Task             models.GroundTask       `json:"task"`
	OtherBooking     models.ResourceBooking  `json:"other_booking"`
	Reason           string                  `json:"reason"`
	AvailableChoices []models.GroundResource `json:"available_choices"`
}

// ReleaseCheckResult is the flight detail aggregation for the dispatcher.
type ReleaseCheckResult struct {
	Turnaround      models.FlightTurnaround `json:"turnaround"`
	CanRelease      bool                    `json:"can_release"`
	Summary         ReleaseSummary          `json:"summary"`
	UnfinishedTasks []models.GroundTask     `json:"unfinished_tasks"`
	OpenDelays      []models.DelayEvent     `json:"open_delays"`
	Conflicts       []ConflictBooking       `json:"conflicts"`
	Violations      []ReleaseViolation      `json:"violations"`
	CheckedAt       time.Time               `json:"checked_at"`
}

type ReleaseSummary struct {
	TotalTasks       int `json:"total_tasks"`
	CompletedTasks   int `json:"completed_tasks"`
	OpenDelays       int `json:"open_delays"`
	BookingConflicts int `json:"booking_conflicts"`
}

// ReleaseRequest optionally carries the acting dispatcher (header fallback).
type ReleaseRequest struct {
	Actor string `json:"actor"`
}

// ReleaseResponse is returned on a successful release.
type ReleaseResponse struct {
	Turnaround models.FlightTurnaround `json:"turnaround"`
	ReleasedAt time.Time               `json:"released_at"`
	Message    string                  `json:"message"`
}
