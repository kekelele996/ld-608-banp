package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"groundTurn/src/constants"
	"groundTurn/src/services"
	"groundTurn/src/types"
)

// ListResourceBooking GET /api/resource-booking?turnaround_id=
// 返回航班全部预约（含已释放），换绑前后两次预约均可查。
func ListResourceBooking(c *gin.Context) {
	turnaroundID, convErr := strconv.Atoi(c.DefaultQuery("turnaround_id", "0"))
	if convErr != nil {
		abortWithError(c, &types.APIError{Code: constants.ValidationFailed, Message: constants.ValidationFailedMessage})
		return
	}
	c.JSON(http.StatusOK, services.ListBookingsByTurnaround(turnaroundID))
}

// ListRebindOptions GET /api/resource-booking/:id/rebind-options
// 冲突任务可换用的可用资源列表。
func ListRebindOptions(c *gin.Context) {
	id, convErr := strconv.Atoi(c.Param("id"))
	if convErr != nil {
		abortWithError(c, &types.APIError{Code: constants.ValidationFailed, Message: constants.ValidationFailedMessage})
		return
	}
	options, err := services.ListRebindOptions(id)
	if err != nil {
		abortWithError(c, err)
		return
	}
	c.JSON(http.StatusOK, options)
}

// RebindResourceBooking POST /api/resource-booking/:id/rebind
// 资源换绑：原预约转为已释放并生成新预约。
func RebindResourceBooking(c *gin.Context) {
	id, convErr := strconv.Atoi(c.Param("id"))
	if convErr != nil {
		abortWithError(c, &types.APIError{Code: constants.ValidationFailed, Message: constants.ValidationFailedMessage})
		return
	}
	var req types.RebindRequest
	if bindErr := c.ShouldBindJSON(&req); bindErr != nil || req.ResourceID <= 0 {
		abortWithError(c, &types.APIError{Code: constants.ValidationFailed, Message: constants.ValidationFailedMessage})
		return
	}
	result, err := services.RebindBooking(id, req.ResourceID)
	if err != nil {
		abortWithError(c, err)
		return
	}
	c.JSON(http.StatusOK, result)
}
