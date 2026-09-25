package types

import "groundTurn/src/models"

type DelayResolveRequest struct {
	Actor string `json:"actor"`
}

type DelayResponse struct {
	Delay models.DelayEvent `json:"delay"`
}
