package controllers

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"groundTurn/src/constructors"
	"groundTurn/src/services"
)

type ResourceController struct {
	svc    *services.ResourceService
	rebook *services.RebookService
}

func NewResourceController(svc *services.ResourceService, rebook *services.RebookService) *ResourceController {
	return &ResourceController{svc: svc, rebook: rebook}
}

func (ctl *ResourceController) List(c *gin.Context) {
	rows, err := ctl.svc.List()
	if err != nil {
		fail(c, err)
		return
	}
	ok(c, rows)
}

// Candidates lists usable resources for swapping a specific booking.
func (ctl *ResourceController) Candidates(c *gin.Context) {
	bookingID, _ := strconv.Atoi(c.Param("bookingId"))
	rows, err := ctl.svc.RebookCandidates(bookingID)
	if err != nil {
		fail(c, err)
		return
	}
	ok(c, rows)
}

// History exposes before/after booking chain for a task.
func (ctl *ResourceController) History(c *gin.Context) {
	taskID, _ := strconv.Atoi(c.Param("taskId"))
	rows, err := ctl.rebook.History(taskID)
	if err != nil {
		fail(c, err)
		return
	}
	ok(c, rows)
}

// Rebook performs the transactional swap for a booking.
func (ctl *ResourceController) Rebook(c *gin.Context) {
	bookingID, _ := strconv.Atoi(c.Param("bookingId"))
	var req struct {
		NewResourceID int    `json:"new_resource_id"`
		StartTime     string `json:"new_start_time"`
		EndTime       string `json:"new_end_time"`
		Actor         string `json:"actor"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		fail(c, err)
		return
	}
	if req.Actor == "" {
		req.Actor = actor(c)
	}
	payload := constructors.NewRebookRequest(req.NewResourceID, req.StartTime, req.EndTime, req.Actor)
	result, err := ctl.rebook.Rebook(bookingID, payload)
	if err != nil {
		fail(c, err)
		return
	}
	ok(c, result)
}
