package services

import (
	"errors"
	"fmt"
	"time"

	"groundTurn/src/constants"
	"groundTurn/src/models"
	"groundTurn/src/repositories"
	"groundTurn/src/types"

	"gorm.io/gorm"
)

// ReleaseService implements the turnaround release coordination workflow:
// it aggregates the three blocking conditions and re-verifies them inside a
// transaction at submission time, so a rejected release never mutates the
// flight status nor any original booking.
type ReleaseService struct {
	repo     *repositories.Repository
	conflict *ConflictService
}

func NewReleaseService(repo *repositories.Repository, conflict *ConflictService) *ReleaseService {
	return &ReleaseService{repo: repo, conflict: conflict}
}

// Check aggregates unfinished tasks, open delays and booking time conflicts
// for the flight detail page. It is read-only and never mutates anything.
func (s *ReleaseService) Check(turnaroundID int) (*types.ReleaseCheckResult, error) {
	turnaround, err := s.repo.GetTurnaround(turnaroundID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, bizError(constants.TurnaroundNotFound,
				fmt.Sprintf(constants.TurnaroundNotFoundMessage, turnaroundID))
		}
		return nil, err
	}

	tasks, err := s.repo.ListTasksByTurnaround(turnaroundID)
	if err != nil {
		return nil, err
	}
	unfinished := make([]models.GroundTask, 0)
	completed := 0
	for _, t := range tasks {
		if t.Status == constants.GroundTaskStatusCompleted {
			completed++
		} else {
			unfinished = append(unfinished, t)
		}
	}

	openDelays, err := s.repo.ListOpenDelaysByTurnaround(turnaroundID)
	if err != nil {
		return nil, err
	}

	conflicts, err := s.conflict.FindConflictsForTurnaround(turnaroundID)
	if err != nil {
		return nil, err
	}

	violations := s.collectViolations(turnaround, unfinished, openDelays, conflicts)

	return &types.ReleaseCheckResult{
		Turnaround:      *turnaround,
		CanRelease:      len(violations) == 0,
		Summary:         types.ReleaseSummary{TotalTasks: len(tasks), CompletedTasks: completed, OpenDelays: len(openDelays), BookingConflicts: len(conflicts)},
		UnfinishedTasks: unfinished,
		OpenDelays:      openDelays,
		Conflicts:       conflicts,
		Violations:      violations,
		CheckedAt:       time.Now(),
	}, nil
}

// SubmitRelease re-checks all three conditions against fresh data inside a
// transaction. If anything fails it returns a BusinessError whose message
// names the exact flight/task/resource/delay; nothing is committed.
func (s *ReleaseService) SubmitRelease(turnaroundID int, actor string) (*types.ReleaseResponse, error) {
	var response *types.ReleaseResponse

	err := s.repo.WithTx(func(tx *gorm.DB) error {
		txRepo := repositories.New(tx)
		txConflict := NewConflictService(txRepo)

		turnaround, err := txRepo.GetTurnaround(turnaroundID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return bizError(constants.TurnaroundNotFound,
					fmt.Sprintf(constants.TurnaroundNotFoundMessage, turnaroundID))
			}
			return err
		}

		for _, terminal := range constants.TurnaroundTerminalStatuses {
			if turnaround.TurnaroundStatus == terminal {
				return bizError(constants.TurnaroundNotReleasable,
					fmt.Sprintf(constants.TurnaroundNotReleasableMessage,
						turnaround.FlightNo, turnaround.TurnaroundStatus))
			}
		}

		tasks, err := txRepo.ListTasksByTurnaround(turnaroundID)
		if err != nil {
			return err
		}
		unfinished := make([]models.GroundTask, 0)
		for _, t := range tasks {
			if t.Status != constants.GroundTaskStatusCompleted {
				unfinished = append(unfinished, t)
			}
		}

		openDelays, err := txRepo.ListOpenDelaysByTurnaround(turnaroundID)
		if err != nil {
			return err
		}

		conflicts, err := txConflict.FindConflictsForTurnaround(turnaroundID)
		if err != nil {
			return err
		}

		violations := s.collectViolations(turnaround, unfinished, openDelays, conflicts)
		if len(violations) > 0 {
			// Returning the error rolls the whole transaction back, so the
			// flight status and original bookings stay untouched. The audit
			// line for the rejected attempt is written after rollback by the
			// caller (see SubmitRelease), via the outer repository.
			return &BusinessError{
				Code:       constants.ReleaseRejected,
				Message:    fmt.Sprintf(constants.ReleaseRejectedMessage, len(violations)),
				Violations: violations,
			}
		}

		if err := txRepo.UpdateTurnaroundStatus(turnaroundID, constants.TurnaroundStatusReady); err != nil {
			return err
		}
		updated, err := txRepo.GetTurnaround(turnaroundID)
		if err != nil {
			return err
		}
		if err := txRepo.WriteAuditLog(actor, "TURNAROUND_RELEASED",
			"FlightTurnaround", fmt.Sprint(turnaroundID),
			fmt.Sprintf(constants.LogTemplates["TURNAROUND_RELEASED"],
				updated.FlightNo, turnaround.TurnaroundStatus, updated.TurnaroundStatus)); err != nil {
			return err
		}
		response = &types.ReleaseResponse{
			Turnaround: *updated,
			ReleasedAt: time.Now(),
			Message:    fmt.Sprintf("航班 %s 放行通过", updated.FlightNo),
		}
		return nil
	})
	if err != nil {
		// The transaction has rolled back; persist the rejection audit line
		// outside the tx so dispatch history keeps the failed attempt while
		// flight status and original bookings remain unchanged.
		var bizErr *BusinessError
		if errors.As(err, &bizErr) && bizErr.Code == constants.ReleaseRejected {
			flightLabel := fmt.Sprint(turnaroundID)
			if t, _ := s.repo.GetTurnaround(turnaroundID); t != nil {
				flightLabel = t.FlightNo
			}
			_ = s.repo.WriteAuditLog(actor, "RELEASE_REJECTED",
				"FlightTurnaround", fmt.Sprint(turnaroundID),
				fmt.Sprintf(constants.LogTemplates["RELEASE_REJECTED"],
					flightLabel, len(bizErr.Violations),
					s.countKind(bizErr.Violations, "task"),
					s.countKind(bizErr.Violations, "delay"),
					s.countKind(bizErr.Violations, "booking")))
		}
		return nil, err
	}
	return response, nil
}

func (s *ReleaseService) countKind(vs []types.ReleaseViolation, kind string) int {
	n := 0
	for _, v := range vs {
		if v.Kind == kind {
			n++
		}
	}
	return n
}

// collectViolations renders the exact blocking reasons in a fixed order:
// tasks -> delays -> booking conflicts.
func (s *ReleaseService) collectViolations(
	turnaround *models.FlightTurnaround,
	unfinished []models.GroundTask,
	openDelays []models.DelayEvent,
	conflicts []types.ConflictBooking,
) []types.ReleaseViolation {
	violations := make([]types.ReleaseViolation, 0)

	for _, t := range unfinished {
		violations = append(violations, types.ReleaseViolation{
			Code:   constants.TaskUnfinished,
			Kind:   "task",
			TaskID: t.ID,
			Detail: fmt.Sprintf(constants.TaskUnfinishedMessage, t.ID, t.TaskType, t.Status),
		})
	}
	for _, d := range openDelays {
		violations = append(violations, types.ReleaseViolation{
			Code:    constants.DelayNotClosed,
			Kind:    "delay",
			DelayID: d.ID,
			Detail: fmt.Sprintf(constants.DelayNotClosedMessage,
				turnaround.FlightNo, d.ID, d.DelayType, d.Minutes),
		})
	}
	for _, c := range conflicts {
		violations = append(violations, types.ReleaseViolation{
			Code:       constants.BookingConflict,
			Kind:       "booking",
			BookingID:  c.Booking.ID,
			ResourceID: c.Booking.ResourceID,
			Detail:     c.Reason,
		})
	}
	return violations
}
