package services

import (
	"errors"
	"fmt"
	"time"

	"groundTurn/src/constants"
	"groundTurn/src/models"
	"groundTurn/src/repositories"

	"gorm.io/gorm"
)

type DelayService struct{ repo *repositories.Repository }

func NewDelayService(repo *repositories.Repository) *DelayService {
	return &DelayService{repo: repo}
}

func (s *DelayService) ListAll() ([]models.DelayEvent, error) {
	return s.repo.ListAllDelays()
}

// Resolve closes an open delay (sets resolved_at), removing it from the
// release blocker list.
func (s *DelayService) Resolve(id int, actor string) (*models.DelayEvent, error) {
	delay, err := s.repo.GetDelay(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, bizError(constants.DelayNotClosed, fmt.Sprintf("延误 #%d 不存在", id))
		}
		return nil, err
	}
	now := time.Now()
	if err := s.repo.ResolveDelay(id, now); err != nil {
		return nil, err
	}
	updated, _ := s.repo.GetDelay(id)
	_ = s.repo.WriteAuditLog(actor, "DELAY_RESOLVED", "DelayEvent", fmt.Sprint(id),
		fmt.Sprintf(constants.LogTemplates["DELAY_RESOLVED"], id, fmt.Sprint(delay.TurnaroundID)))
	return updated, nil
}
