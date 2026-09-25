package models

// ResourceBooking 资源预约：连接航班、任务和资源。
type ResourceBooking struct {
	ID             int    `json:"id"`
	ResourceID     int    `json:"resource_id"`
	TurnaroundID   int    `json:"turnaround_id"`
	TaskID         int    `json:"task_id"`
	StartTime      string `json:"start_time"`
	EndTime        string `json:"end_time"`
	BookingStatus  string `json:"booking_status"`
	ConflictReason string `json:"conflict_reason"`
}
