package controllers

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"groundTurn/src/services"
)

type DelayController struct{ svc *services.DelayService }

func NewDelayController(svc *services.DelayService) *DelayController {
	return &DelayController{svc: svc}
}

func (ctl *DelayController) List(c *gin.Context) {
	rows, err := ctl.svc.ListAll()
	if err != nil {
		fail(c, err)
		return
	}
	ok(c, rows)
}

// Resolve closes a delay so the flight can pass release re-check.
func (ctl *DelayController) Resolve(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	act := actor(c)
	delay, err := ctl.svc.Resolve(id, act)
	if err != nil {
		fail(c, err)
		return
	}
	ok(c, delay)
}
