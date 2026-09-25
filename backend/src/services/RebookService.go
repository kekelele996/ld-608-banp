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

// RebookService swaps a conflicting booking to another usable resource.
// In one transaction it: validates the target resource (right type,
// AVAILABLE, time-window free), flips the original booking to RELEASED and
// creates a new CONFIRMED booking, linking the two rows in both directions
// so before/after history is queryable from either side.
type RebookService struct {
	repo     *repositories.Repository
	conflict *ConflictService
}

func NewRebookService(repo *repositories.Repository, conflict *ConflictService) *RebookService {
	return &RebookService{repo: repo, conflict: conflict}
}

func (s *RebookService) Rebook(bookingID int, req types.RebookRequest) (*types.RebookResult, error) {
	var result *types.RebookResult

	err := s.repo.WithTx(func(tx *gorm.DB) error {
		txRepo := repositories.New(tx)
		txConflict := NewConflictService(txRepo)

		oldBooking, err := txRepo.GetBooking(bookingID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return bizError(constants.BookingNotFound,
					fmt.Sprintf(constants.BookingNotFoundMessage, bookingID))
			}
			return err
		}
		if oldBooking.BookingStatus == constants.BookingStatusReleased {
			return bizError(constants.BookingNotRebindable,
				fmt.Sprintf(constants.BookingNotRebindableMessage, bookingID, oldBooking.BookingStatus))
		}

		target, err := txRepo.GetResource(req.NewResourceID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return bizError(constants.ResourceNotFound,
					fmt.Sprintf(constants.ResourceNotFoundMessage, req.NewResourceID))
			}
			return err
		}

		task, err := txRepo.GetTask(oldBooking.TaskID)
		if err != nil {
			return err
		}
		if target.ResourceType != task.TaskType {
			return bizError(constants.ResourceTypeMismatch,
				fmt.Sprintf(constants.ResourceTypeMismatchMessage,
					target.ResourceCode, target.ResourceType, task.TaskType))
		}
		if target.AvailabilityStatus != constants.ResourceStatusAvailable {
			return bizError(constants.ResourceUnavailable,
				fmt.Sprintf(constants.ResourceUnavailableMessage,
					target.ResourceCode, target.AvailabilityStatus))
		}

		// Keep the original window when the caller omits new times.
		start, end := oldBooking.StartTime, oldBooking.EndTime
		if !req.StartTime.IsZero() && !req.EndTime.IsZero() {
			start, end = req.StartTime, req.EndTime
		}
		if !start.Before(end) {
			return bizError(constants.ResourceBusy, "换绑时间窗无效：开始时间必须早于结束时间")
		}

		busy, err := txConflict.IsResourceFree(target.ID, bookingID, start, end)
		if err != nil {
			return err
		}
		if busy != nil {
			return bizError(constants.ResourceBusy,
				fmt.Sprintf(constants.ResourceBusyMessage,
					target.ResourceCode,
					start.Format("15:04"), end.Format("15:04"), busy.ID))
		}

		oldResource, err := txRepo.GetResource(oldBooking.ResourceID)
		if err != nil {
			return err
		}

		newBooking := &models.ResourceBooking{
			ResourceID:        target.ID,
			TurnaroundID:      oldBooking.TurnaroundID,
			TaskID:            oldBooking.TaskID,
			StartTime:         start,
			EndTime:           end,
			BookingStatus:     constants.BookingStatusConfirmed,
			ConflictReason:    "",
			ReplacedBookingID: &oldBooking.ID,
		}
		if err := txRepo.CreateBooking(newBooking); err != nil {
			return err
		}

		if err := txRepo.ReleaseBookingChain(oldBooking.ID, newBooking.ID,
			fmt.Sprintf("换绑至资源 %s（新预约 #%d）", target.ResourceCode, newBooking.ID)); err != nil {
			return err
		}
		if err := txRepo.LinkNewBookingBack(newBooking.ID, oldBooking.ID); err != nil {
			return err
		}
		newBooking.ReplacedBookingID = &oldBooking.ID

		// Recalculate statuses: old resource may become AVAILABLE if no other
		// active booking uses it; target becomes BOOKED.
		if err := syncResourceStatus(txRepo, oldResource.ID); err != nil {
			return err
		}
		if err := txRepo.UpdateResourceStatus(target.ID, constants.ResourceStatusBooked); err != nil {
			return err
		}

		if err := txRepo.WriteAuditLog(req.Actor, "RESOURCE_REBOUND",
			"ResourceBooking", fmt.Sprint(newBooking.ID),
			fmt.Sprintf(constants.LogTemplates["RESOURCE_REBOUND"],
				oldBooking.ID, newBooking.ID,
				oldResource.ResourceCode, target.ResourceCode, oldBooking.TaskID)); err != nil {
			return err
		}

		released := *oldBooking
		released.BookingStatus = constants.BookingStatusReleased
		released.ReplacedByID = &newBooking.ID
		result = &types.RebookResult{
			ReleasedBooking: released,
			NewBooking:      *newBooking,
			OldResource:     *oldResource,
			NewResource:     *target,
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return result, nil
}

// History returns the before/after booking chain for a task with related
// entities attached (resource code, task type, flight number).
func (s *RebookService) History(taskID int) ([]types.BookingHistoryItem, error) {
	rows, err := s.repo.ListBookingHistory(taskID)
	if err != nil {
		return nil, err
	}
	items := make([]types.BookingHistoryItem, 0, len(rows))
	for _, b := range rows {
		item := types.BookingHistoryItem{Booking: b}
		if res, err := s.repo.GetResource(b.ResourceID); err == nil {
			item.Resource = *res
		}
		if task, err := s.repo.GetTask(b.TaskID); err == nil {
			item.Task = *task
			if flight, err := s.repo.GetTurnaround(task.TurnaroundID); err == nil {
				item.Flight = *flight
			}
		}
		items = append(items, item)
	}
	return items, nil
}

// syncResourceStatus marks a resource BOOKED while it has any active booking,
// otherwise AVAILABLE. MAINTENANCE/OFFLINE resources are left untouched.
func syncResourceStatus(repo *repositories.Repository, resourceID int) error {
	resource, err := repo.GetResource(resourceID)
	if err != nil {
		return err
	}
	if resource.AvailabilityStatus == constants.ResourceStatusMaintenance ||
		resource.AvailabilityStatus == constants.ResourceStatusOffline {
		return nil
	}
	active, err := repo.ListActiveBookingsOnResource(resourceID)
	if err != nil {
		return err
	}
	desired := constants.ResourceStatusAvailable
	if len(active) > 0 {
		desired = constants.ResourceStatusBooked
	}
	return repo.UpdateResourceStatus(resourceID, desired)
}
