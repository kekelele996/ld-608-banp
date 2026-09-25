package repositories

import "groundTurn/src/models"

// 以下函数均不加锁，调用方（service 层）必须先 Lock。

func ListResources() []models.GroundResource {
	rows := make([]models.GroundResource, 0, len(db.resources))
	for i := 1; i <= len(db.resources); i++ {
		if r, ok := db.resources[i]; ok {
			rows = append(rows, r)
		}
	}
	return rows
}

func GetResource(id int) (models.GroundResource, bool) {
	r, ok := db.resources[id]
	return r, ok
}
