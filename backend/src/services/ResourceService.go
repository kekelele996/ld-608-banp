package services

import (
	"errors"
	"fmt"

	"groundTurn/src/constants"
	"groundTurn/src/models"
	"groundTurn/src/repositories"
	"groundTurn/src/types"

	"gorm.io/gorm"
)

type ResourceService struct {
	repo     *repositories.Repository
	conflict *ConflictService
}

func NewResourceService(repo *repositories.Repository, conflict *ConflictService) *ResourceService {
	return &ResourceService{repo: repo, conflict: conflict}
}

func (s *ResourceService) List() ([]models.GroundResource, error) {
	return s.repo.ListResources()
}

// RerbookCandidates returns resources of the booking's task type annotated
// with usability (AVAILABLE + window free) for the swap modal.
func (s *ResourceService) RebookCandidates(bookingID int) ([]types.UsableResource, error) {
	booking, err := s.repo.GetBooking(bookingID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, bizError(constants.BookingNotFound,
				fmt.Sprintf(constants.BookingNotFoundMessage, bookingID))
		}
		return nil, err
	}
	task, err := s.repo.GetTask(booking.TaskID)
	if err != nil {
		return nil, err
	}
	candidates, err := s.repo.ListResourcesByType(task.TaskType)
	if err != nil {
		return nil, err
	}
	out := make([]types.UsableResource, 0, len(candidates))
	for _, res := range candidates {
		item := types.UsableResource{Resource: res, Usable: true}
		if res.ID == booking.ResourceID {
			item.Usable = false
			item.Reason = "当前已预约该资源"
		} else if res.AvailabilityStatus != constants.ResourceStatusAvailable {
			item.Usable = false
			item.Reason = fmt.Sprintf("资源状态 %s", res.AvailabilityStatus)
		} else {
			busy, err := s.conflict.IsResourceFree(res.ID, booking.ID, booking.StartTime, booking.EndTime)
			if err != nil {
				return nil, err
			}
			if busy != nil {
				item.Usable = false
				item.Reason = fmt.Sprintf("时间窗与预约 #%d 冲突", busy.ID)
			}
		}
		out = append(out, item)
	}
	return out, nil
}
