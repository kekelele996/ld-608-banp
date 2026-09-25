package main

import (
	"groundTurn/src/config"
	"groundTurn/src/routes"
)

func main() {
	routes.Start(":" + config.ServerPort())
}
