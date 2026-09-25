package constants

// GroundTaskStatus 地勤任务状态：放行前必须全部 COMPLETED。
var GroundTaskStatus = []string{"PENDING", "IN_PROGRESS", "BLOCKED", "COMPLETED"}

const (
	GroundTaskStatusPending    = "PENDING"
	GroundTaskStatusInProgress = "IN_PROGRESS"
	GroundTaskStatusBlocked    = "BLOCKED"
	GroundTaskStatusCompleted  = "COMPLETED"
)
