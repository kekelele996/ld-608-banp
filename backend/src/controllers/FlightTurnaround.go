package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"groundTurn/src/constants"
	"groundTurn/src/services"
	"groundTurn/src/types"
)

// statusOf 错误码到 HTTP 状态码的映射，controller 层自行包装异常。
func statusOf(err *types.APIError) int {
	switch err.Code {
	case constants.TurnaroundNotFound, constants.TaskNotFound, constants.ResourceNotFound, constants.BookingNotFound:
		return http.StatusNotFound
	case constants.ReleaseCheckFailed, constants.ResourceUnavailable, constants.RebindConflict:
		return http.StatusConflict
	default:
		return http.StatusBadRequest
	}
}

func abortWithError(c *gin.Context, err *types.APIError) {
	c.JSON(statusOf(err), gin.H{"error": err})
}

// ListFlightTurnaround GET /api/flight-turnaround
func ListFlightTurnaround(c *gin.Context) {
	c.JSON(http.StatusOK, services.ListTurnarounds())
}

// GetReleaseSummary GET /api/flight-turnaround/:id/release-summary
// 航班详情汇总：未完成任务、未关闭延误、预约冲突与放行结论。
func GetReleaseSummary(c *gin.Context) {
	id, convErr := strconv.Atoi(c.Param("id"))
	if convErr != nil {
		abortWithError(c, &types.APIError{Code: constants.ValidationFailed, Message: constants.ValidationFailedMessage})
		return
	}
	summary, err := services.BuildReleaseSummary(id)
	if err != nil {
		abortWithError(c, err)
		return
	}
	c.JSON(http.StatusOK, summary)
}

// ReleaseFlightTurnaround POST /api/flight-turnaround/:id/release
// 提交放行：后端重新核对三类条件，不满足则 409 并逐项指出，数据保持不变。
func ReleaseFlightTurnaround(c *gin.Context) {
	id, convErr := strconv.Atoi(c.Param("id"))
	if convErr != nil {
		abortWithError(c, &types.APIError{Code: constants.ValidationFailed, Message: constants.ValidationFailedMessage})
		return
	}
	result, err := services.ReleaseTurnaround(id)
	if err != nil {
		abortWithError(c, err)
		return
	}
	c.JSON(http.StatusOK, result)
}
