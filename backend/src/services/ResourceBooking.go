package services

import (
	"fmt"
	"log"

	"groundTurn/src/constants"
	"groundTurn/src/constructors"
	"groundTurn/src/models"
	"groundTurn/src/repositories"
	"groundTurn/src/types"
	"groundTurn/src/utils"
)

// ListBookingsByTurnaround 航班预约历史（含已释放），可查换绑前后两次预约。
func ListBookingsByTurnaround(turnaroundID int) []models.ResourceBooking {
	repositories.Lock()
	defer repositories.Unlock()
	return repositories.ListBookingsByTurnaround(turnaroundID)
}

// FindBookingConflictsLocked 检测航班生效预约与同资源其他生效预约的时间重叠，
// 调用方必须已持有仓储锁。
func FindBookingConflictsLocked(turnaroundID int) []types.BookingConflict {
	conflicts := []types.BookingConflict{}
	all := repositories.ListBookings()
	for _, b := range repositories.ListBookingsByTurnaround(turnaroundID) {
		if !constants.IsActiveBookingStatus(b.BookingStatus) {
			continue
		}
		for _, other := range all {
			if other.ID == b.ID || other.ResourceID != b.ResourceID || !constants.IsActiveBookingStatus(other.BookingStatus) {
				continue
			}
			if utils.Overlaps(b.StartTime, b.EndTime, other.StartTime, other.EndTime) {
				resource, _ := repositories.GetResource(b.ResourceID)
				conflicts = append(conflicts, types.BookingConflict{
					Booking:       b,
					ConflictsWith: other,
					Resource:      resource,
					Reason:        fmt.Sprintf("%s~%s", b.StartTime, b.EndTime),
				})
			}
		}
	}
	return conflicts
}

// ListRebindOptions 换绑候选资源：同类型、AVAILABLE 且在预约时段内无冲突。
func ListRebindOptions(bookingID int) ([]models.GroundResource, *types.APIError) {
	repositories.Lock()
	defer repositories.Unlock()
	booking, ok := repositories.GetBooking(bookingID)
	if !ok || !constants.IsActiveBookingStatus(booking.BookingStatus) {
		return nil, &types.APIError{Code: constants.BookingNotFound, Message: constants.BookingNotFoundMessage}
	}
	current, _ := repositories.GetResource(booking.ResourceID)
	options := []models.GroundResource{}
	for _, r := range repositories.ListResources() {
		if r.ID == booking.ResourceID || r.ResourceType != current.ResourceType || r.AvailabilityStatus != "AVAILABLE" {
			continue
		}
		if resourceFreeLocked(r.ID, booking.StartTime, booking.EndTime) {
			options = append(options, r)
		}
	}
	return options, nil
}

// RebindBooking 资源换绑：校验通过后原预约转为已释放并生成新预约，
// 两步在同一锁内完成，历史记录可查前后两次预约。
func RebindBooking(bookingID int, newResourceID int) (types.RebindResult, *types.APIError) {
	repositories.Lock()
	defer repositories.Unlock()
	booking, ok := repositories.GetBooking(bookingID)
	if !ok || !constants.IsActiveBookingStatus(booking.BookingStatus) {
		return types.RebindResult{}, &types.APIError{Code: constants.BookingNotFound, Message: constants.BookingNotFoundMessage}
	}
	resource, ok := repositories.GetResource(newResourceID)
	if !ok {
		return types.RebindResult{}, &types.APIError{Code: constants.ResourceNotFound, Message: constants.ResourceNotFoundMessage}
	}
	if resource.AvailabilityStatus != "AVAILABLE" {
		return types.RebindResult{}, &types.APIError{
			Code:    constants.ResourceUnavailable,
			Message: fmt.Sprintf("%s：%s 当前状态 %s", constants.ResourceUnavailableMessage, resource.ResourceCode, resource.AvailabilityStatus),
		}
	}
	if !resourceFreeLocked(newResourceID, booking.StartTime, booking.EndTime) {
		return types.RebindResult{}, &types.APIError{
			Code:    constants.RebindConflict,
			Message: fmt.Sprintf("%s：%s 在 %s~%s 已有生效预约", constants.RebindConflictMessage, resource.ResourceCode, booking.StartTime, booking.EndTime),
		}
	}
	released, _ := repositories.UpdateBookingStatus(booking.ID, constants.BookingStatusReleased,
		fmt.Sprintf("已换绑至资源 %s", resource.ResourceCode))
	created := repositories.CreateBooking(constructors.NewReboundBooking(booking, newResourceID))
	log.Printf("%s: booking=%d resource=%d", constants.LogTemplates["ResourceBooking"][4], released.ID, released.ResourceID)
	log.Printf("%s: booking=%d resource=%d", constants.LogTemplates["ResourceBooking"][5], created.ID, created.ResourceID)
	return constructors.NewRebindResult(released, created), nil
}

// resourceFreeLocked 资源在指定时段内是否没有生效预约，调用方必须已持有仓储锁。
func resourceFreeLocked(resourceID int, start, end string) bool {
	for _, b := range repositories.ListBookings() {
		if b.ResourceID == resourceID && constants.IsActiveBookingStatus(b.BookingStatus) && utils.Overlaps(start, end, b.StartTime, b.EndTime) {
			return false
		}
	}
	return true
}
