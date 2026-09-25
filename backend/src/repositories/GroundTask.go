package repositories

import "groundTurn/src/models"

func (r *Repository) ListTasksByTurnaround(turnaroundID int) ([]models.GroundTask, error) {
	var rows []models.GroundTask
	err := r.db.Where("turnaround_id = ?", turnaroundID).
		Order("deadline asc").Find(&rows).Error
	return rows, err
}

func (r *Repository) ListAllTasks() ([]models.GroundTask, error) {
	var rows []models.GroundTask
	err := r.db.Order("id asc").Find(&rows).Error
	return rows, err
}

func (r *Repository) GetTask(id int) (*models.GroundTask, error) {
	var row models.GroundTask
	if err := r.db.First(&row, id).Error; err != nil {
		return nil, err
	}
	return &row, nil
}

func (r *Repository) UpdateTaskStatus(id int, status, blockerNote string, actualFinish interface{}) error {
	updates := map[string]interface{}{
		"status":       status,
		"blocker_note": blockerNote,
	}
	if actualFinish != nil {
		updates["actual_finish"] = actualFinish
	}
	return r.db.Model(&models.GroundTask{}).Where("id = ?", id).Updates(updates).Error
}

func (r *Repository) CreateTask(t *models.GroundTask) error {
	return r.db.Create(t).Error
}
