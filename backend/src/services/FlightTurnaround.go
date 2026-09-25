package services

import (
	"log"

	"groundTurn/src/constants"
	"groundTurn/src/constructors"
	"groundTurn/src/models"
	"groundTurn/src/repositories"
	"groundTurn/src/types"
)

// ListTurnarounds 航班过站列表。
func ListTurnarounds() []models.FlightTurnaround {
	repositories.Lock()
	defer repositories.Unlock()
	return repositories.ListTurnarounds()
}

// BuildReleaseSummary 航班详情汇总：未完成任务、未关闭延误、预约时间冲突与放行结论。
func BuildReleaseSummary(turnaroundID int) (types.ReleaseSummary, *types.APIError) {
	repositories.Lock()
	defer repositories.Unlock()
	turnaround, ok := repositories.GetTurnaround(turnaroundID)
	if !ok {
		return types.ReleaseSummary{}, &types.APIError{Code: constants.TurnaroundNotFound, Message: constants.TurnaroundNotFoundMessage}
	}
	unfinished, openDelays, conflicts, blockers := collectReleaseBlockers(turnaround)
	return constructors.NewReleaseSummary(turnaround, unfinished, openDelays, conflicts, blockers), nil
}

// ReleaseTurnaround 提交放行：后端重新核对三类条件，任一不满足则指出具体
// 航班、任务或资源，航班状态和原预约保持不变；全部通过才置为 READY。
func ReleaseTurnaround(turnaroundID int) (types.ReleaseResult, *types.APIError) {
	repositories.Lock()
	defer repositories.Unlock()
	turnaround, ok := repositories.GetTurnaround(turnaroundID)
	if !ok {
		return types.ReleaseResult{}, &types.APIError{Code: constants.TurnaroundNotFound, Message: constants.TurnaroundNotFoundMessage}
	}
	_, _, _, blockers := collectReleaseBlockers(turnaround)
	if len(blockers) > 0 {
		log.Printf("%s: flight=%s blockers=%d", constants.LogTemplates["FlightTurnaround"][5], turnaround.FlightNo, len(blockers))
		return types.ReleaseResult{}, &types.APIError{
			Code:     constants.ReleaseCheckFailed,
			Message:  constants.ReleaseCheckFailedMessage,
			Blockers: blockers,
		}
	}
	updated, _ := repositories.UpdateTurnaroundStatus(turnaroundID, "READY")
	log.Printf("%s: flight=%s status=READY", constants.LogTemplates["FlightTurnaround"][4], turnaround.FlightNo)
	return types.ReleaseResult{OK: true, Turnaround: updated}, nil
}

// collectReleaseBlockers 汇总三类放行条件，调用方必须已持有仓储锁。
func collectReleaseBlockers(turnaround models.FlightTurnaround) (
	unfinished []models.GroundTask,
	openDelays []models.DelayEvent,
	conflicts []types.BookingConflict,
	blockers []types.ReleaseBlocker,
) {
	unfinished = []models.GroundTask{}
	openDelays = []models.DelayEvent{}
	blockers = []types.ReleaseBlocker{}

	for _, t := range repositories.ListTasksByTurnaround(turnaround.ID) {
		if t.Status != constants.GroundTaskStatusCompleted {
			unfinished = append(unfinished, t)
			blockers = append(blockers, constructors.NewTaskBlocker(t, turnaround.FlightNo))
		}
	}
	for _, d := range repositories.ListDelaysByTurnaround(turnaround.ID) {
		if d.ResolvedAt == "" {
			openDelays = append(openDelays, d)
			blockers = append(blockers, constructors.NewDelayBlocker(d, turnaround.FlightNo))
		}
	}
	conflicts = FindBookingConflictsLocked(turnaround.ID)
	for _, c := range conflicts {
		blockers = append(blockers, constructors.NewBookingBlocker(c))
	}
	return unfinished, openDelays, conflicts, blockers
}
