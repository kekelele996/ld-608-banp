package constants

// BookingStatus 资源预约状态：换绑后原预约转为 RELEASED，新预约为 CONFIRMED。
var BookingStatus = []string{"PENDING", "CONFIRMED", "RELEASED", "CANCELLED"}

const (
	BookingStatusPending   = "PENDING"
	BookingStatusConfirmed = "CONFIRMED"
	BookingStatusReleased  = "RELEASED"
	BookingStatusCancelled = "CANCELLED"
)

// IsActiveBookingStatus 预约是否仍占用资源（释放/取消后不再参与冲突检测）。
func IsActiveBookingStatus(status string) bool {
	return status == BookingStatusPending || status == BookingStatusConfirmed
}
