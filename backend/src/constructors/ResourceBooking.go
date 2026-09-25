package constructors

import (
	"groundTurn/src/constants"
	"groundTurn/src/models"
	"groundTurn/src/types"
)

// NewRebindResult 构造换绑响应：原预约（已释放）+ 新预约，历史可查前后两次。
func NewRebindResult(released, created models.ResourceBooking) types.RebindResult {
	return types.RebindResult{Released: released, Created: created}
}

// NewReboundBooking 由冲突预约派生新预约对象：换绑资源，时段与任务不变。
func NewReboundBooking(origin models.ResourceBooking, newResourceID int) models.ResourceBooking {
	return models.ResourceBooking{
		ResourceID:     newResourceID,
		TurnaroundID:   origin.TurnaroundID,
		TaskID:         origin.TaskID,
		StartTime:      origin.StartTime,
		EndTime:        origin.EndTime,
		BookingStatus:  constants.BookingStatusConfirmed,
		ConflictReason: "",
	}
}
