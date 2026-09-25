package constants

// GroundTaskStatus enumerates the lifecycle of a ground task.
// A task is considered finished only when it equals COMPLETED.
var GroundTaskStatus = []string{"PENDING", "IN_PROGRESS", "BLOCKED", "COMPLETED"}

const (
	GroundTaskStatusPending    = "PENDING"
	GroundTaskStatusInProgress = "IN_PROGRESS"
	GroundTaskStatusBlocked    = "BLOCKED"
	GroundTaskStatusCompleted  = "COMPLETED"
)
