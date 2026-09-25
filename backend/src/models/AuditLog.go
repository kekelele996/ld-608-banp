package models

import "time"

// AuditLog records every write action (release, rebook, dispatch, ...).
type AuditLog struct {
	ID         int       `gorm:"primaryKey;autoIncrement" json:"id"`
	Actor      string    `gorm:"column:actor;size:64;not null" json:"actor"`
	Action     string    `gorm:"column:action;size:64;not null;index" json:"action"`
	TargetType string    `gorm:"column:target_type;size:64;not null" json:"target_type"`
	TargetID   string    `gorm:"column:target_id;size:64" json:"target_id"`
	Detail     string    `gorm:"column:detail;size:1024" json:"detail"`
	CreatedAt  time.Time `gorm:"column:created_at;index" json:"created_at"`
}

func (AuditLog) TableName() string { return "audit_log" }
