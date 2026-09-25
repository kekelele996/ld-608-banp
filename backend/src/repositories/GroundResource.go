package repositories

import "groundTurn/src/models"

func (r *Repository) ListResources() ([]models.GroundResource, error) {
	var rows []models.GroundResource
	err := r.db.Order("id asc").Find(&rows).Error
	return rows, err
}

func (r *Repository) ListResourcesByType(resourceType string) ([]models.GroundResource, error) {
	var rows []models.GroundResource
	err := r.db.Where("resource_type = ?", resourceType).
		Order("id asc").Find(&rows).Error
	return rows, err
}

func (r *Repository) GetResource(id int) (*models.GroundResource, error) {
	var row models.GroundResource
	if err := r.db.First(&row, id).Error; err != nil {
		return nil, err
	}
	return &row, nil
}

func (r *Repository) UpdateResourceStatus(id int, status string) error {
	return r.db.Model(&models.GroundResource{}).
		Where("id = ?", id).
		Update("availability_status", status).Error
}

func (r *Repository) CreateResource(res *models.GroundResource) error {
	return r.db.Create(res).Error
}
