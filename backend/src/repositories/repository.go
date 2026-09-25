package repositories

import (
	"time"

	"groundTurn/src/models"

	"gorm.io/gorm"
)

// Repository is the single data-access entry point. Services may either use
// it directly or call WithTx for transactional workflows (release/rebook).
type Repository struct {
	db *gorm.DB
}

func New(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) DB() *gorm.DB { return r.db }

// WithTx runs fn inside a transaction; returning an error rolls everything
// back so a rejected release leaves statuses and bookings untouched.
func (r *Repository) WithTx(fn func(tx *gorm.DB) error) error {
	return r.db.Transaction(fn)
}

func (r *Repository) WriteAuditLog(actor, action, targetType, targetID, detail string) error {
	return r.db.Create(&models.AuditLog{
		Actor:      actor,
		Action:     action,
		TargetType: targetType,
		TargetID:   targetID,
		Detail:     detail,
		CreatedAt:  time.Now(),
	}).Error
}

func (r *Repository) ListAuditLogs(limit int) ([]models.AuditLog, error) {
	var rows []models.AuditLog
	err := r.db.Order("created_at desc, id desc").Limit(limit).Find(&rows).Error
	return rows, err
}
