package controllers

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"groundTurn/src/services"
	"groundTurn/src/types"
)

type TaskController struct{ svc *services.TaskService }

func NewTaskController(svc *services.TaskService) *TaskController {
	return &TaskController{svc: svc}
}

func (ctl *TaskController) List(c *gin.Context) {
	if raw := c.Query("turnaround_id"); raw != "" {
		id, _ := strconv.Atoi(raw)
		rows, err := ctl.svc.ListByTurnaround(id)
		if err != nil {
			fail(c, err)
			return
		}
		ok(c, rows)
		return
	}
	rows, err := ctl.svc.ListAll()
	if err != nil {
		fail(c, err)
		return
	}
	ok(c, rows)
}

// UpdateStatus handles dispatch / block / complete (sign-off).
func (ctl *TaskController) UpdateStatus(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var req types.TaskStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		fail(c, err)
		return
	}
	if req.Actor == "" {
		req.Actor = actor(c)
	}
	task, err := ctl.svc.UpdateStatus(id, req)
	if err != nil {
		fail(c, err)
		return
	}
	ok(c, task)
}
