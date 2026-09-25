package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"groundTurn/src/controllers"
)

func Start(addr string) {
	r := gin.Default()
	r.GET("/health", func(c *gin.Context) { c.JSON(http.StatusOK, gin.H{"status": "ok", "service": "ground-turn"}) })

	r.GET("/api/flight-turnaround", controllers.ListFlightTurnaround)
	r.GET("/api/flight-turnaround/:id/release-summary", controllers.GetReleaseSummary)
	r.POST("/api/flight-turnaround/:id/release", controllers.ReleaseFlightTurnaround)

	r.GET("/api/ground-task", controllers.ListGroundTask)
	r.GET("/api/ground-resource", controllers.ListGroundResource)

	r.GET("/api/resource-booking", controllers.ListResourceBooking)
	r.GET("/api/resource-booking/:id/rebind-options", controllers.ListRebindOptions)
	r.POST("/api/resource-booking/:id/rebind", controllers.RebindResourceBooking)

	r.GET("/api/delay-event", controllers.ListDelayEvent)

	r.Run(addr)
}
