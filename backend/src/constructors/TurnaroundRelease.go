package constructors

import (
	"fmt"

	"groundTurn/src/constants"
	"groundTurn/src/models"
	"groundTurn/src/types"
)

// NewReleaseSummary 构造航班放行汇总响应对象，页面不得自行拼装。
func NewReleaseSummary(
	turnaround models.FlightTurnaround,
	unfinished []models.GroundTask,
	openDelays []models.DelayEvent,
	conflicts []types.BookingConflict,
	blockers []types.ReleaseBlocker,
) types.ReleaseSummary {
	conclusion := constants.ReleaseConclusionReleasable
	if len(blockers) > 0 {
		conclusion = constants.ReleaseConclusionBlocked
	}
	return types.ReleaseSummary{
		Turnaround:       turnaround,
		UnfinishedTasks:  unfinished,
		OpenDelays:       openDelays,
		BookingConflicts: conflicts,
		Conclusion:       conclusion,
		Blockers:         blockers,
	}
}

// NewTaskBlocker 未完成任务阻碍项：指出具体航班与任务。
func NewTaskBlocker(t models.GroundTask, flightNo string) types.ReleaseBlocker {
	return types.ReleaseBlocker{
		Kind:    constants.ReleaseBlockerKindTask,
		RefID:   t.ID,
		Message: fmt.Sprintf("航班 %s 任务 #%d（%s）未完成，当前状态 %s", flightNo, t.ID, t.TaskType, t.Status),
	}
}

// NewDelayBlocker 未关闭延误阻碍项：指出具体航班与延误事件。
func NewDelayBlocker(d models.DelayEvent, flightNo string) types.ReleaseBlocker {
	return types.ReleaseBlocker{
		Kind:    constants.ReleaseBlockerKindDelay,
		RefID:   d.ID,
		Message: fmt.Sprintf("航班 %s 延误事件 #%d（%s，%d 分钟）未关闭", flightNo, d.ID, d.DelayType, d.Minutes),
	}
}

// NewBookingBlocker 预约冲突阻碍项：指出具体资源与两次预约。
func NewBookingBlocker(c types.BookingConflict) types.ReleaseBlocker {
	return types.ReleaseBlocker{
		Kind:    constants.ReleaseBlockerKindBooking,
		RefID:   c.Booking.ID,
		Message: fmt.Sprintf("资源 %s 预约 #%d 与预约 #%d 时段重叠（%s）", c.Resource.ResourceCode, c.Booking.ID, c.ConflictsWith.ID, c.Reason),
	}
}
