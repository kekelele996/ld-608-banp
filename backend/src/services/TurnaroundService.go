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

type TurnaroundService struct{ repo *repositories.Repository }

func NewTurnaroundService(repo *repositories.Repository) *TurnaroundService {
	return &TurnaroundService{repo: repo}
}

func (s *TurnaroundService) List() ([]models.FlightTurnaround, error) {
	return s.repo.ListTurnarounds()
}

func (s *TurnaroundService) Get(id int) (*models.FlightTurnaround, error) {
	t, err := s.repo.GetTurnaround(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, bizError(constants.TurnaroundNotFound,
				fmt.Sprintf(constants.TurnaroundNotFoundMessage, id))
		}
		return nil, err
	}
	return t, nil
}

// DetailBundle is the full flight detail payload (tasks, delays, bookings).
type DetailBundle struct {
	Turnaround models.FlightTurnaround    `json:"turnaround"`
	Tasks      []models.GroundTask        `json:"tasks"`
	Delays     []models.DelayEvent        `json:"delays"`
	Bookings   []types.BookingHistoryItem `json:"bookings"`
}

func (s *TurnaroundService) Detail(id int) (*DetailBundle, error) {
	if _, err := s.Get(id); err != nil {
		return nil, err
	}
	bundle := &DetailBundle{}
	t, _ := s.repo.GetTurnaround(id)
	bundle.Turnaround = *t

	tasks, err := s.repo.ListTasksByTurnaround(id)
	if err != nil {
		return nil, err
	}
	bundle.Tasks = tasks
	delays, err := s.repo.ListAllDelays()
	if err != nil {
		return nil, err
	}
	for _, d := range delays {
		if d.TurnaroundID == id {
			bundle.Delays = append(bundle.Delays, d)
		}
	}

	activeBookings, err := s.repo.ListActiveBookingsByTurnaround(id)
	if err != nil {
		return nil, err
	}
	for _, b := range activeBookings {
		item := types.BookingHistoryItem{Booking: b}
		if res, err := s.repo.GetResource(b.ResourceID); err == nil {
			item.Resource = *res
		}
		if task, err := s.repo.GetTask(b.TaskID); err == nil {
			item.Task = *task
		}
		item.Flight = bundle.Turnaround
		bundle.Bookings = append(bundle.Bookings, item)
	}
	return bundle, nil
}
