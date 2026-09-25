package repositories

import (
	"time"

	"groundTurn/src/models"
)

// ListActiveBookingsByTurnaround returns non-released bookings (CONFIRMED or
// CONFLICT) for a turnaround.
func (r *Repository) ListActiveBookingsByTurnaround(turnaroundID int) ([]models.ResourceBooking, error) {
	var rows []models.ResourceBooking
	err := r.db.Where("turnaround_id = ? AND booking_status <> ?",
		turnaroundID, "RELEASED").
		Order("start_time asc").Find(&rows).Error
	return rows, err
}

func (r *Repository) GetBooking(id int) (*models.ResourceBooking, error) {
	var row models.ResourceBooking
	if err := r.db.First(&row, id).Error; err != nil {
		return nil, err
	}
	return &row, nil
}

// ListActiveBookingsOnResource returns non-released bookings overlapping
// other turnarounds; used for conflict detection and rebook feasibility.
func (r *Repository) ListActiveBookingsOnResource(resourceID int) ([]models.ResourceBooking, error) {
	var rows []models.ResourceBooking
	err := r.db.Where("resource_id = ? AND booking_status <> ?",
		resourceID, "RELEASED").
		Order("start_time asc").Find(&rows).Error
	return rows, err
}

// OverlappingBookings returns active bookings of the resource whose
// [start,end) intersects the given window, excluding excludeBookingID.
func (r *Repository) OverlappingBookings(resourceID, excludeBookingID int, start, end time.Time) ([]models.ResourceBooking, error) {
	var rows []models.ResourceBooking
	err := r.db.Where(`
			resource_id = ?
			AND booking_status <> ?
			AND id <> ?
			AND start_time < ?
			AND end_time > ?`,
		resourceID, "RELEASED", excludeBookingID, end, start).
		Find(&rows).Error
	return rows, err
}

func (r *Repository) CreateBooking(b *models.ResourceBooking) error {
	return r.db.Create(b).Error
}

// UpdateBookingStatusChain flips an old booking to RELEASED while recording
// the pointer to its replacement. Called inside the rebook transaction.
func (r *Repository) ReleaseBookingChain(id, replacedByID int, conflictReason string) error {
	return r.db.Model(&models.ResourceBooking{}).Where("id = ?", id).
		Updates(map[string]interface{}{
			"booking_status":  "RELEASED",
			"replaced_by_id":  replacedByID,
			"conflict_reason": conflictReason,
		}).Error
}

// LinkNewBookingBack fills replaced_booking_id on the freshly created row.
func (r *Repository) LinkNewBookingBack(newID, oldID int) error {
	return r.db.Model(&models.ResourceBooking{}).Where("id = ?", newID).
		Update("replaced_booking_id", oldID).Error
}

func (r *Repository) MarkBookingConflict(id int, reason string) error {
	return r.db.Model(&models.ResourceBooking{}).Where("id = ?", id).
		Updates(map[string]interface{}{
			"booking_status":  "CONFLICT",
			"conflict_reason": reason,
		}).Error
}

// ListBookingHistory returns the before/after chain for one task: the
// released predecessors and the current active booking, newest first.
func (r *Repository) ListBookingHistory(taskID int) ([]models.ResourceBooking, error) {
	var rows []models.ResourceBooking
	err := r.db.Where("task_id = ?", taskID).
		Order("created_at desc, id desc").Find(&rows).Error
	return rows, err
}
