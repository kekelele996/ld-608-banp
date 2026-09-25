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

type TaskService struct{ repo *repositories.Repository }

func NewTaskService(repo *repositories.Repository) *TaskService {
	return &TaskService{repo: repo}
}

func (s *TaskService) ListByTurnaround(turnaroundID int) ([]models.GroundTask, error) {
	return s.repo.ListTasksByTurnaround(turnaroundID)
}

func (s *TaskService) ListAll() ([]models.GroundTask, error) {
	return s.repo.ListAllTasks()
}

// UpdateStatus changes task lifecycle. Completing stamps actual_finish;
// blocking requires a blocker note.
func (s *TaskService) UpdateStatus(id int, req types.TaskStatusRequest) (*models.GroundTask, error) {
	task, err := s.repo.GetTask(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, bizError(constants.TaskNotFound,
				fmt.Sprintf(constants.TaskNotFoundMessage, id))
		}
		return nil, err
	}

	var finish interface{}
	switch req.Status {
	case constants.GroundTaskStatusCompleted:
		now := time.Now()
		finish = now
	case constants.GroundTaskStatusBlocked:
		if req.BlockerNote == "" {
			return nil, bizError(constants.TaskUnfinished, "阻塞任务必须填写阻塞原因")
		}
	}

	if err := s.repo.UpdateTaskStatus(id, req.Status, req.BlockerNote, finish); err != nil {
		return nil, err
	}
	updated, err := s.repo.GetTask(id)
	if err != nil {
		return nil, err
	}

	action := "TASK_DISPATCHED"
	switch req.Status {
	case constants.GroundTaskStatusCompleted:
		action = "TASK_COMPLETED"
	case constants.GroundTaskStatusBlocked:
		action = "TASK_BLOCKED"
	}
	_ = s.repo.WriteAuditLog(req.Actor, action, "GroundTask", fmt.Sprint(id),
		fmt.Sprintf("#%d %s -> %s", id, task.TaskType, req.Status))
	return updated, nil
}
