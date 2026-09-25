package routes

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"groundTurn/src/config"
	"groundTurn/src/controllers"
	"groundTurn/src/middlewares"
	"groundTurn/src/repositories"
	"groundTurn/src/services"

	"gorm.io/gorm"
)

// Start wires repositories -> services -> controllers -> routes, then boots
// Gin. Dependency direction always points inward; no package reaches across.
func Start(addr string) {
	var db *gorm.DB
	var err error
	for attempt := 1; attempt <= 30; attempt++ {
		db, err = config.Connect(config.GetDBConfig())
		if err == nil {
			break
		}
		log.Printf("database init attempt %d/30 failed: %v", attempt, err)
		time.Sleep(2 * time.Second)
	}
	if err != nil {
		log.Fatalf("database init failed after retries: %v", err)
	}
	if err := seedIfEmpty(db); err != nil {
		log.Printf("seed skipped: %v", err)
	}

	repo := repositories.New(db)
	conflictSvc := services.NewConflictService(repo)
	turnaroundSvc := services.NewTurnaroundService(repo)
	releaseSvc := services.NewReleaseService(repo, conflictSvc)
	rebookSvc := services.NewRebookService(repo, conflictSvc)
	taskSvc := services.NewTaskService(repo)
	delaySvc := services.NewDelayService(repo)
	resourceSvc := services.NewResourceService(repo, conflictSvc)

	turnaroundCtl := controllers.NewTurnaroundController(turnaroundSvc, releaseSvc)
	taskCtl := controllers.NewTaskController(taskSvc)
	delayCtl := controllers.NewDelayController(delaySvc)
	resourceCtl := controllers.NewResourceController(resourceSvc, rebookSvc)

	r := gin.New()
	r.Use(gin.Logger(), middlewares.ErrorHandlerMiddleware(),
		middlewares.AuditLogMiddleware(), middlewares.AuthMiddleware())

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok", "service": "ground-turn"})
	})

	api := r.Group("/api")
	{
		// Flight turnarounds + release coordination
		api.GET("/turnarounds", turnaroundCtl.List)
		api.GET("/turnarounds/:id", turnaroundCtl.Detail)
		api.GET("/turnarounds/:id/release-check", turnaroundCtl.ReleaseCheck)
		api.POST("/turnarounds/:id/release",
			middlewares.RbacMiddleware(middlewares.RoleDispatcher, middlewares.RoleSupervisor),
			turnaroundCtl.SubmitRelease)

		// Ground tasks (dispatch / sign-off / block)
		api.GET("/ground-tasks", taskCtl.List)
		api.PATCH("/ground-tasks/:id/status", taskCtl.UpdateStatus)

		// Delays
		api.GET("/delay-events", delayCtl.List)
		api.POST("/delay-events/:id/resolve",
			middlewares.RbacMiddleware(middlewares.RoleDispatcher, middlewares.RoleSupervisor),
			delayCtl.Resolve)

		// Resources + booking rebook
		api.GET("/ground-resources", resourceCtl.List)
		api.GET("/bookings/:bookingId/candidates", resourceCtl.Candidates)
		api.GET("/bookings/task/:taskId/history", resourceCtl.History)
		api.POST("/bookings/:bookingId/rebook",
			middlewares.RbacMiddleware(middlewares.RoleDispatcher, middlewares.RoleResource, middlewares.RoleSupervisor),
			resourceCtl.Rebook)
	}

	log.Printf("ground-turn backend listening on %s", addr)
	if err := r.Run(addr); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}
