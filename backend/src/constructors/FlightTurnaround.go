package constructors

import (
	"time"

	"groundTurn/src/models"
)

// NewFlightTurnaround builds a default turnaround for registration forms.
func NewFlightTurnaround(flightNo, aircraftReg, standNo string, arrival, departure time.Time) models.FlightTurnaround {
	return models.FlightTurnaround{
		FlightNo:         flightNo,
		AircraftReg:      aircraftReg,
		StandNo:          standNo,
		ArrivalTime:      arrival,
		DepartureTime:    departure,
		TurnaroundStatus: "ON_STAND",
	}
}
