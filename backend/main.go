package main

import (
	"os"

	"groundTurn/src/routes"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}
	routes.Start(":" + port)
}
