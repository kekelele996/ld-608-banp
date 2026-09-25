package constructors

import (
	"time"

	"groundTurn/src/models"
	"groundTurn/src/types"
)

// NewRebookRequest is the only place controllers map raw request fields
// (ISO strings from the client) onto the service DTO.
func NewRebookRequest(resourceID int, startISO, endISO, actor string) types.RebookRequest {
	req := types.RebookRequest{NewResourceID: resourceID, Actor: actor}
	if t, err := time.Parse(time.RFC3339, startISO); err == nil {
		req.StartTime = t
	}
	if t, err := time.Parse(time.RFC3339, endISO); err == nil {
		req.EndTime = t
	}
	return req
}

// NewReleaseViolation builds a violation row for service-layer callers.
func NewReleaseViolation(code, kind, detail string, taskID, bookingID, resourceID, delayID int) types.ReleaseViolation {
	return types.ReleaseViolation{
		Code: code, Kind: kind, Detail: detail,
		TaskID: taskID, BookingID: bookingID, ResourceID: resourceID, DelayID: delayID,
	}
}

// NewTurnaroundResponse wraps a model for list responses.
func NewTurnaroundResponse(m models.FlightTurnaround) models.FlightTurnaround { return m }
