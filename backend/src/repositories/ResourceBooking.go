package repositories

import "groundTurn/src/models"

// 以下函数均不加锁，调用方（service 层）必须先 Lock。

func ListBookings() []models.ResourceBooking {
	rows := make([]models.ResourceBooking, 0, len(db.bookings))
	for i := 1; i <= len(db.bookings); i++ {
		if b, ok := db.bookings[i]; ok {
			rows = append(rows, b)
		}
	}
	return rows
}

// ListBookingsByTurnaround 返回航班全部预约（含已释放），历史记录可查前后两次预约。
func ListBookingsByTurnaround(turnaroundID int) []models.ResourceBooking {
	rows := []models.ResourceBooking{}
	for i := 1; i <= len(db.bookings); i++ {
		if b, ok := db.bookings[i]; ok && b.TurnaroundID == turnaroundID {
			rows = append(rows, b)
		}
	}
	return rows
}

func GetBooking(id int) (models.ResourceBooking, bool) {
	b, ok := db.bookings[id]
	return b, ok
}

// UpdateBookingStatus 换绑时把原预约转为 RELEASED。
func UpdateBookingStatus(id int, status string, conflictReason string) (models.ResourceBooking, bool) {
	b, ok := db.bookings[id]
	if !ok {
		return b, false
	}
	b.BookingStatus = status
	b.ConflictReason = conflictReason
	db.bookings[id] = b
	return b, true
}

// CreateBooking 换绑时生成新预约，id 自增。
func CreateBooking(b models.ResourceBooking) models.ResourceBooking {
	b.ID = db.nextBooking
	db.nextBooking++
	db.bookings[b.ID] = b
	return b
}
