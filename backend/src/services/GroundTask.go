package services

import (
	"groundTurn/src/models"
	"groundTurn/src/repositories"
)

// ListTasks 地勤任务列表。
func ListTasks() []models.GroundTask {
	repositories.Lock()
	defer repositories.Unlock()
	return repositories.ListTasks()
}
