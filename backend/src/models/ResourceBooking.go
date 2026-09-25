package models

import "time"

// ResourceBooking binds a resource to a task within a turnaround.
//
// Rebind history is modelled as an in-table linked list: when a booking is
// swapped onto another resource, the old row is flipped to RELEASED with
// ReplacedByID pointing at the new row, and the new row carries
// ReplacedBookingID back. Both rows survive so the before/after pair is
// queryable from either side.
type ResourceBooking struct {
	ID                int       `gorm:"primaryKey;autoIncrement" json:"id"`
	ResourceID        int       `gorm:"column:resource_id;not null;index" json:"resource_id"`
	TurnaroundID      int       `gorm:"column:turnaround_id;not null;index" json:"turnaround_id"`
	TaskID            int       `gorm:"column:task_id;not null;index" json:"task_id"`
	StartTime         time.Time `gorm:"column:start_time;not null;index" json:"start_time"`
	EndTime           time.Time `gorm:"column:end_time;not null" json:"end_time"`
	BookingStatus     string    `gorm:"column:booking_status;size:32;not null;index" json:"booking_status"`
	ConflictReason    string    `gorm:"column:conflict_reason;size:512" json:"conflict_reason"`
	ReplacedBookingID *int      `gorm:"column:replaced_booking_id;index" json:"replaced_booking_id"`
	ReplacedByID      *int      `gorm:"column:replaced_by_id;index" json:"replaced_by_id"`
	CreatedAt         time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt         time.Time `gorm:"column:updated_at" json:"updated_at"`
}

func (ResourceBooking) TableName() string { return "resource_booking" }
