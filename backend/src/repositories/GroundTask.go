package repositories

import "groundTurn/src/models"

// 以下函数均不加锁，调用方（service 层）必须先 Lock。

func ListTasks() []models.GroundTask {
	rows := make([]models.GroundTask, 0, len(db.tasks))
	for i := 1; i <= len(db.tasks); i++ {
		if t, ok := db.tasks[i]; ok {
			rows = append(rows, t)
		}
	}
	return rows
}

func ListTasksByTurnaround(turnaroundID int) []models.GroundTask {
	rows := []models.GroundTask{}
	for _, t := range db.tasks {
		if t.TurnaroundID == turnaroundID {
			rows = append(rows, t)
		}
	}
	return rows
}

func GetTask(id int) (models.GroundTask, bool) {
	t, ok := db.tasks[id]
	return t, ok
}
