package repositories

import "groundTurn/src/models"

// 以下函数均不加锁，调用方（service 层）必须先 Lock。

func ListDelays() []models.DelayEvent {
	rows := make([]models.DelayEvent, 0, len(db.delays))
	for i := 1; i <= len(db.delays); i++ {
		if d, ok := db.delays[i]; ok {
			rows = append(rows, d)
		}
	}
	return rows
}

func ListDelaysByTurnaround(turnaroundID int) []models.DelayEvent {
	rows := []models.DelayEvent{}
	for _, d := range db.delays {
		if d.TurnaroundID == turnaroundID {
			rows = append(rows, d)
		}
	}
	return rows
}
