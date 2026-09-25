package models

import "time"

// DelayEvent is an unresolved/resolved delay attributed to a team.
// ResolvedAt == nil means the delay is still open and blocks release.
type DelayEvent struct {
	ID                 int        `gorm:"primaryKey;autoIncrement" json:"id"`
	TurnaroundID       int        `gorm:"column:turnaround_id;not null;index" json:"turnaround_id"`
	DelayType          string     `gorm:"column:delay_type;size:32" json:"delay_type"`
	Minutes            int        `gorm:"column:minutes" json:"minutes"`
	RootCause          string     `gorm:"column:root_cause;size:512" json:"root_cause"`
	ResponsibilityTeam string     `gorm:"column:responsibility_team;size:64" json:"responsibility_team"`
	ResolvedAt         *time.Time `gorm:"column:resolved_at;index" json:"resolved_at"`
	CreatedAt          time.Time  `gorm:"column:created_at" json:"created_at"`
	UpdatedAt          time.Time  `gorm:"column:updated_at" json:"updated_at"`
}

func (DelayEvent) TableName() string { return "delay_event" }
