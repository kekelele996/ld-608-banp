package models

import "time"

// FlightTurnaround is a flight on-stand turnaround being serviced.
type FlightTurnaround struct {
	ID               int       `gorm:"primaryKey;autoIncrement" json:"id"`
	FlightNo         string    `gorm:"column:flight_no;size:32;not null" json:"flight_no"`
	AircraftReg      string    `gorm:"column:aircraft_reg;size:32" json:"aircraft_reg"`
	StandNo          string    `gorm:"column:stand_no;size:32" json:"stand_no"`
	ArrivalTime      time.Time `gorm:"column:arrival_time" json:"arrival_time"`
	DepartureTime    time.Time `gorm:"column:departure_time" json:"departure_time"`
	TurnaroundStatus string    `gorm:"column:turnaround_status;size:32;not null;index" json:"turnaround_status"`
	DelayReason      string    `gorm:"column:delay_reason;size:512" json:"delay_reason"`
	CreatedAt        time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt        time.Time `gorm:"column:updated_at" json:"updated_at"`
}

func (FlightTurnaround) TableName() string { return "flight_turnaround" }
