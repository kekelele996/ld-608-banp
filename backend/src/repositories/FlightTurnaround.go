package repositories

import "groundTurn/src/models"

// Lock/Unlock 由 service 层调用，把“校验 + 写入”包成原子操作：
// 放行复核不通过时航班状态和原预约保持不变，换绑时释放与新建同时生效。
func Lock()   { db.mu.Lock() }
func Unlock() { db.mu.Unlock() }

// 以下函数均不加锁，调用方（service 层）必须先 Lock。

func ListTurnarounds() []models.FlightTurnaround {
	rows := make([]models.FlightTurnaround, 0, len(db.turnarounds))
	for i := 1; i <= len(db.turnarounds); i++ {
		if t, ok := db.turnarounds[i]; ok {
			rows = append(rows, t)
		}
	}
	return rows
}

func GetTurnaround(id int) (models.FlightTurnaround, bool) {
	t, ok := db.turnarounds[id]
	return t, ok
}

// UpdateTurnaroundStatus 仅在放行复核全部通过后调用。
func UpdateTurnaroundStatus(id int, status string) (models.FlightTurnaround, bool) {
	t, ok := db.turnarounds[id]
	if !ok {
		return t, false
	}
	t.TurnaroundStatus = status
	db.turnarounds[id] = t
	return t, true
}
