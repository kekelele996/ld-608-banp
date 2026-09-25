package services

import (
	"groundTurn/src/models"
	"groundTurn/src/repositories"
)

// ListDelays 延误事件列表。
func ListDelays() []models.DelayEvent {
	repositories.Lock()
	defer repositories.Unlock()
	return repositories.ListDelays()
}
