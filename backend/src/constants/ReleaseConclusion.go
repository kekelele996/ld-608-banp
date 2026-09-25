package constants

// ReleaseConclusion 放行结论：汇总页与提交放行共用同一套判定。
var ReleaseConclusion = []string{"RELEASABLE", "BLOCKED"}

const (
	ReleaseConclusionReleasable = "RELEASABLE"
	ReleaseConclusionBlocked    = "BLOCKED"
)

// 放行复核的三类条件标识，用于 blocker 的 kind 字段。
const (
	ReleaseBlockerKindTask    = "TASK"
	ReleaseBlockerKindDelay   = "DELAY"
	ReleaseBlockerKindBooking = "BOOKING"
)
