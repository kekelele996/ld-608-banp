package constructors

import "groundTurn/src/models"

// NewGroundResource builds a default AVAILABLE resource for the ledger form.
func NewGroundResource(code, resourceType, location, ownerTeam string) models.GroundResource {
	return models.GroundResource{
		ResourceCode:       code,
		ResourceType:       resourceType,
		Location:           location,
		OwnerTeam:          ownerTeam,
		AvailabilityStatus: "AVAILABLE",
	}
}
