package services

import (
	"groundTurn/src/models"
	"groundTurn/src/repositories"
)

// ListResources 保障资源列表。
func ListResources() []models.GroundResource {
	repositories.Lock()
	defer repositories.Unlock()
	return repositories.ListResources()
}
