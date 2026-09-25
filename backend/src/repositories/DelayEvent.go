package repositories

import "groundTurn/src/models"

// ListOpenDelaysByTurnaround returns delays without resolved_at.
func (r *Repository) ListOpenDelaysByTurnaround(turnaroundID int) ([]models.DelayEvent, error) {
	var rows []models.DelayEvent
	err := r.db.Where("turnaround_id = ? AND resolved_at IS NULL", turnaroundID).
		Order("id asc").Find(&rows).Error
	return rows, err
}

func (r *Repository) ListAllDelays() ([]models.DelayEvent, error) {
	var rows []models.DelayEvent
	err := r.db.Order("id asc").Find(&rows).Error
	return rows, err
}

func (r *Repository) GetDelay(id int) (*models.DelayEvent, error) {
	var row models.DelayEvent
	if err := r.db.First(&row, id).Error; err != nil {
		return nil, err
	}
	return &row, nil
}

func (r *Repository) ResolveDelay(id int, resolvedAt interface{}) error {
	return r.db.Model(&models.DelayEvent{}).Where("id = ?", id).
		Update("resolved_at", resolvedAt).Error
}

func (r *Repository) CreateDelay(d *models.DelayEvent) error {
	return r.db.Create(d).Error
}
