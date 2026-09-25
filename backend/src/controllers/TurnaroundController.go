package controllers

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"groundTurn/src/services"
	"groundTurn/src/types"
)

type TurnaroundController struct {
	svc     *services.TurnaroundService
	release *services.ReleaseService
}

func NewTurnaroundController(svc *services.TurnaroundService, release *services.ReleaseService) *TurnaroundController {
	return &TurnaroundController{svc: svc, release: release}
}

// List returns all turnarounds for the release worklist.
func (ctl *TurnaroundController) List(c *gin.Context) {
	rows, err := ctl.svc.List()
	if err != nil {
		fail(c, err)
		return
	}
	ok(c, rows)
}

// Detail returns tasks/delays/bookings for one flight.
func (ctl *TurnaroundController) Detail(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	bundle, err := ctl.svc.Detail(id)
	if err != nil {
		fail(c, err)
		return
	}
	ok(c, bundle)
}

// ReleaseCheck aggregates the three blocker categories (read-only).
func (ctl *TurnaroundController) ReleaseCheck(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	result, err := ctl.release.Check(id)
	if err != nil {
		fail(c, err)
		return
	}
	ok(c, result)
}

// SubmitRelease re-verifies server-side and flips status only when all pass.
func (ctl *TurnaroundController) SubmitRelease(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var req types.ReleaseRequest
	_ = c.ShouldBindJSON(&req)
	if req.Actor == "" {
		req.Actor = actor(c)
	}
	result, err := ctl.release.SubmitRelease(id, req.Actor)
	if err != nil {
		fail(c, err)
		return
	}
	ok(c, result)
}
