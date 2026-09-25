package models

import "time"

// GroundTask is a single servicing task belonging to a turnaround.
type GroundTask struct {
	ID           int        `gorm:"primaryKey;autoIncrement" json:"id"`
	TurnaroundID int        `gorm:"column:turnaround_id;not null;index" json:"turnaround_id"`
	TaskType     string     `gorm:"column:task_type;size:32;not null" json:"task_type"`
	TeamID       int        `gorm:"column:team_id" json:"team_id"`
	PlannedStart time.Time  `gorm:"column:planned_start" json:"planned_start"`
	Deadline     time.Time  `gorm:"column:deadline" json:"deadline"`
	ActualFinish *time.Time `gorm:"column:actual_finish" json:"actual_finish"`
	Status       string     `gorm:"column:status;size:32;not null;index" json:"status"`
	BlockerNote  string     `gorm:"column:blocker_note;size:512" json:"blocker_note"`
	CreatedAt    time.Time  `gorm:"column:created_at" json:"created_at"`
	UpdatedAt    time.Time  `gorm:"column:updated_at" json:"updated_at"`
}

func (GroundTask) TableName() string { return "ground_task" }
