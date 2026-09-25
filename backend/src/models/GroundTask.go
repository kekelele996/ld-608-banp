package models

// GroundTask 地勤任务：属于过站航班并消耗资源。
type GroundTask struct {
	ID           int    `json:"id"`
	TurnaroundID int    `json:"turnaround_id"`
	TaskType     string `json:"task_type"`
	TeamID       int    `json:"team_id"`
	PlannedStart string `json:"planned_start"`
	Deadline     string `json:"deadline"`
	ActualFinish string `json:"actual_finish"`
	Status       string `json:"status"`
	BlockerNote  string `json:"blocker_note"`
}
