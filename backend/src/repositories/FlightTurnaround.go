package repositories

import "groundTurn/src/models"

func (r *Repository) ListTurnarounds() ([]models.FlightTurnaround, error) {
	var rows []models.FlightTurnaround
	err := r.db.Order("arrival_time asc").Find(&rows).Error
	return rows, err
}

func (r *Repository) GetTurnaround(id int) (*models.FlightTurnaround, error) {
	var row models.FlightTurnaround
	if err := r.db.First(&row, id).Error; err != nil {
		return nil, err
	}
	return &row, nil
}

func (r *Repository) UpdateTurnaroundStatus(id int, status string) error {
	return r.db.Model(&models.FlightTurnaround{}).
		Where("id = ?", id).
		Update("turnaround_status", status).Error
}

func (r *Repository) CreateTurnaround(t *models.FlightTurnaround) error {
	return r.db.Create(t).Error
}
