package main

import (
	"log"

	"moving-schedule-backend/config"
	"moving-schedule-backend/database"
	"moving-schedule-backend/handlers"
	"moving-schedule-backend/middleware"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("加载配置失败: %v", err)
	}

	if err := database.Init(cfg); err != nil {
		log.Fatalf("初始化数据库失败: %v", err)
	}

	r := gin.Default()

	r.Use(middleware.CORS())

	api := r.Group("/api")
	{
		orders := api.Group("/orders")
		{
			orders.POST("", handlers.CreateOrder)
			orders.GET("", handlers.GetOrders)
			orders.GET("/:id", handlers.GetOrder)
			orders.PUT("/:id", handlers.UpdateOrder)
			orders.DELETE("/:id", handlers.DeleteOrder)
		}

		schedules := api.Group("/schedules")
		{
			schedules.GET("", handlers.GetSchedules)
			schedules.POST("", handlers.CreateSchedule)
		}

		workers := api.Group("/workers")
		{
			workers.GET("", handlers.GetWorkers)
			workers.POST("", handlers.CreateWorker)
			workers.GET("/:id/dispatches", handlers.GetWorkerDispatchesByID)
		}

		vehicles := api.Group("/vehicles")
		{
			vehicles.GET("", handlers.GetVehicles)
			vehicles.POST("", handlers.CreateVehicle)
		}

		dispatches := api.Group("/dispatches")
		{
			dispatches.POST("", handlers.CreateDispatch)
			dispatches.GET("/check", handlers.CheckDispatch)
			dispatches.GET("", handlers.GetDispatches)
			dispatches.GET("/:id", handlers.GetDispatch)
			dispatches.PUT("/:id/accept", handlers.AcceptDispatch)
			dispatches.PUT("/:id/start", handlers.StartDispatch)
			dispatches.PUT("/:id/complete", handlers.CompleteDispatch)
		}

		worker := api.Group("/worker")
		{
			worker.GET("/my-dispatches", handlers.GetWorkerDispatches)
			worker.POST("/dispatches/:id/accept", handlers.AcceptDispatch)
			worker.POST("/dispatches/:id/start", handlers.StartDispatch)
			worker.POST("/dispatches/:id/complete", handlers.CompleteDispatch)
		}
	}

	log.Println("服务器启动在 :8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("启动服务器失败: %v", err)
	}
}
