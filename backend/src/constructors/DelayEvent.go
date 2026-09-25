package constructors

import "groundTurn/src/models"

// NewDelayEvent builds a default open delay for the reporting form.
func NewDelayEvent(turnaroundID int, delayType string, minutes int, cause, team string) models.DelayEvent {
	return models.DelayEvent{
		TurnaroundID:       turnaroundID,
		DelayType:          delayType,
		Minutes:            minutes,
		RootCause:          cause,
		ResponsibilityTeam: team,
	}
}
