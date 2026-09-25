package constructors

import (
	"time"

	"groundTurn/src/models"
)

// NewGroundTask builds a default pending task for dispatch forms.
func NewGroundTask(turnaroundID, teamID int, taskType string, start, deadline time.Time) models.GroundTask {
	return models.GroundTask{
		TurnaroundID: turnaroundID,
		TeamID:       teamID,
		TaskType:     taskType,
		PlannedStart: start,
		Deadline:     deadline,
		Status:       "PENDING",
	}
}
