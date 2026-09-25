package models

// FlightTurnaround 航班过站：拥有多个 GroundTask 和 ResourceBooking。
type FlightTurnaround struct {
	ID               int    `json:"id"`
	FlightNo         string `json:"flight_no"`
	AircraftReg      string `json:"aircraft_reg"`
	StandNo          string `json:"stand_no"`
	ArrivalTime      string `json:"arrival_time"`
	DepartureTime    string `json:"departure_time"`
	TurnaroundStatus string `json:"turnaround_status"`
	DelayReason      string `json:"delay_reason"`
}
