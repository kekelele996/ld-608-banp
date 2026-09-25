package types

import "groundTurn/src/models"

// ReleaseBlocker 放行阻碍项：指出具体航班、任务或资源。
type ReleaseBlocker struct {
	Kind    string `json:"kind"`    // TASK / DELAY / BOOKING
	RefID   int    `json:"ref_id"`  // 任务、延误或预约的 id
	Message string `json:"message"` // 面向调度员的具体描述
}

// BookingConflict 预约时间冲突：本航班预约与同资源另一预约重叠。
type BookingConflict struct {
	Booking       models.ResourceBooking `json:"booking"`
	ConflictsWith models.ResourceBooking `json:"conflicts_with"`
	Resource      models.GroundResource  `json:"resource"`
	Reason        string                 `json:"reason"`
}

// ReleaseSummary 航班详情汇总：未完成任务、未关闭延误、预约冲突与放行结论。
type ReleaseSummary struct {
	Turnaround       models.FlightTurnaround `json:"turnaround"`
	UnfinishedTasks  []models.GroundTask     `json:"unfinished_tasks"`
	OpenDelays       []models.DelayEvent     `json:"open_delays"`
	BookingConflicts []BookingConflict       `json:"booking_conflicts"`
	Conclusion       string                  `json:"conclusion"`
	Blockers         []ReleaseBlocker        `json:"blockers"`
}

// ReleaseResult 提交放行的结果。
type ReleaseResult struct {
	OK         bool                    `json:"ok"`
	Turnaround models.FlightTurnaround `json:"turnaround"`
}

// RebindResult 资源换绑结果：原预约（已释放）与新预约，历史可查前后两次。
type RebindResult struct {
	Released models.ResourceBooking `json:"released"`
	Created  models.ResourceBooking `json:"created"`
}

// APIError 统一错误响应，blockers 用于放行复核失败时逐项说明。
type APIError struct {
	Code     string           `json:"code"`
	Message  string           `json:"message"`
	Blockers []ReleaseBlocker `json:"blockers,omitempty"`
}
