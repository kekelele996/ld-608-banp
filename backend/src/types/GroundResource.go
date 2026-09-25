package types

import "groundTurn/src/models"

// UsableResource is a rebook candidate with its overlap verdict precomputed.
type UsableResource struct {
	Resource models.GroundResource `json:"resource"`
	Usable   bool                  `json:"usable"`
	Reason   string                `json:"reason,omitempty"`
}
