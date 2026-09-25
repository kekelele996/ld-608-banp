package services

import (
	"fmt"
	"time"

	"groundTurn/src/constants"
	"groundTurn/src/models"
	"groundTurn/src/repositories"
	"groundTurn/src/types"
)

// ConflictService evaluates time-overlap conflicts between bookings.
// Two active bookings conflict when they share a resource and their
// [start,end) windows intersect, irrespective of turnaround.
type ConflictService struct{ repo *repositories.Repository }

func NewConflictService(repo *repositories.Repository) *ConflictService {
	return &ConflictService{repo: repo}
}

// WindowsOverlap reports whether [s1,e1) intersects [s2,e2).
// Touching boundaries (e1 == s2) do not conflict.
func WindowsOverlap(s1, e1, s2, e2 time.Time) bool {
	return s1.Before(e2) && s2.Before(e1)
}

// FindConflictsForTurnaround lists every active booking of the turnaround
// that overlaps another active booking on the same resource, alongside
// candidate resources the dispatcher may swap to.
func (s *ConflictService) FindConflictsForTurnaround(turnaroundID int) ([]types.ConflictBooking, error) {
	bookings, err := s.repo.ListActiveBookingsByTurnaround(turnaroundID)
	if err != nil {
		return nil, err
	}

	result := make([]types.ConflictBooking, 0)
	for _, b := range bookings {
		others, err := s.repo.ListActiveBookingsOnResource(b.ResourceID)
		if err != nil {
			return nil, err
		}
		for _, other := range others {
			if other.ID == b.ID || other.TurnaroundID == b.TurnaroundID {
				continue
			}
			if !WindowsOverlap(b.StartTime, b.EndTime, other.StartTime, other.EndTime) {
				continue
			}
			resource, _ := s.repo.GetResource(b.ResourceID)
			task, _ := s.repo.GetTask(b.TaskID)
			choices, _ := s.availableChoices(task.TaskType, b.ID, b.StartTime, b.EndTime)
			result = append(result, types.ConflictBooking{
				Booking:          b,
				Resource:         safeResource(resource),
				Task:             safeTask(task),
				OtherBooking:     other,
				Reason:           s.renderReason(b, other, resource),
				AvailableChoices: choices,
			})
			break // one conflict entry per booking is enough for the panel
		}
	}
	return result, nil
}

// HasConflicts is the fast boolean used by the release re-check.
func (s *ConflictService) HasConflicts(turnaroundID int) (bool, []types.ReleaseViolation, error) {
	conflicts, err := s.FindConflictsForTurnaround(turnaroundID)
	if err != nil {
		return false, nil, err
	}
	violations := make([]types.ReleaseViolation, 0, len(conflicts))
	for _, c := range conflicts {
		violations = append(violations, types.ReleaseViolation{
			Code:       constants.BookingConflict,
			Kind:       "booking",
			BookingID:  c.Booking.ID,
			ResourceID: c.Booking.ResourceID,
			Detail:     c.Reason,
		})
	}
	return len(conflicts) > 0, violations, nil
}

// IsResourceFree reports whether resourceID has no overlapping active booking
// in [start,end), excluding the booking being moved.
func (s *ConflictService) IsResourceFree(resourceID, excludeBookingID int, start, end time.Time) (*models.ResourceBooking, error) {
	overlaps, err := s.repo.OverlappingBookings(resourceID, excludeBookingID, start, end)
	if err != nil {
		return nil, err
	}
	if len(overlaps) > 0 {
		return &overlaps[0], nil
	}
	return nil, nil
}

// availableChoices returns resources of the given task type that are
// AVAILABLE and free across the required window.
func (s *ConflictService) availableChoices(taskType string, excludeBookingID int, start, end time.Time) ([]models.GroundResource, error) {
	candidates, err := s.repo.ListResourcesByType(taskType)
	if err != nil {
		return nil, err
	}
	choices := make([]models.GroundResource, 0)
	for _, res := range candidates {
		if res.AvailabilityStatus != constants.ResourceStatusAvailable {
			continue
		}
		busy, err := s.IsResourceFree(res.ID, excludeBookingID, start, end)
		if err != nil {
			return nil, err
		}
		if busy == nil {
			choices = append(choices, res)
		}
	}
	return choices, nil
}

func (s *ConflictService) renderReason(b, other models.ResourceBooking, resource *models.GroundResource) string {
	code := fmt.Sprintf("资源 #%d", b.ResourceID)
	if resource != nil {
		code = resource.ResourceCode
	}
	return fmt.Sprintf(constants.BookingConflictMessage, b.ID, other.ID, code)
}

func safeResource(r *models.GroundResource) models.GroundResource {
	if r == nil {
		return models.GroundResource{}
	}
	return *r
}

func safeTask(t *models.GroundTask) models.GroundTask {
	if t == nil {
		return models.GroundTask{}
	}
	return *t
}
