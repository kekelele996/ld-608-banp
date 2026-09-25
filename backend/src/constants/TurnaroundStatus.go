package constants

var TurnaroundStatus = []string{"ARRIVING", "ON_STAND", "IN_SERVICE", "READY", "DEPARTED", "DELAYED"}

const (
	TurnaroundStatusArriving  = "ARRIVING"
	TurnaroundStatusOnStand   = "ON_STAND"
	TurnaroundStatusInService = "IN_SERVICE"
	TurnaroundStatusReady     = "READY"
	TurnaroundStatusDeparted  = "DEPARTED"
	TurnaroundStatusDelayed   = "DELAYED"
)

// Terminal statuses cannot be submitted for release again.
var TurnaroundTerminalStatuses = []string{TurnaroundStatusReady, TurnaroundStatusDeparted}
