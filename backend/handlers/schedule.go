package handlers

import (
	"net/http"
	"time"

	"moving-schedule-backend/database"
	"moving-schedule-backend/models"

	"github.com/gin-gonic/gin"
)

type CreateWorkerRequest struct {
	Name   string `json:"name" binding:"required"`
	Phone  string `json:"phone" binding:"required"`
	Status string `json:"status"`
}

type CreateVehicleRequest struct {
	PlateNumber    string  `json:"plate_number" binding:"required"`
	VehicleType    string  `json:"vehicle_type" binding:"required"`
	CapacityVolume float64 `json:"capacity_volume" binding:"required"`
	CapacityWeight float64 `json:"capacity_weight" binding:"required"`
	Status         string  `json:"status"`
}

type CreateScheduleRequest struct {
	WorkerID  uint   `json:"worker_id" binding:"required"`
	VehicleID uint   `json:"vehicle_id" binding:"required"`
	WorkDate  string `json:"work_date" binding:"required"`
	Shift     string `json:"shift" binding:"required"`
	Status    string `json:"status"`
}

func GetWorkers(c *gin.Context) {
	var workers []models.Worker
	if err := database.DB.Find(&workers).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "查询师傅列表失败: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, workers)
}

func CreateWorker(c *gin.Context) {
	var req CreateWorkerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	worker := models.Worker{
		Name:  req.Name,
		Phone: req.Phone,
	}

	if req.Status != "" {
		worker.Status = models.WorkerStatus(req.Status)
	} else {
		worker.Status = models.WorkerStatusAvailable
	}

	if err := database.DB.Create(&worker).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "新增师傅失败: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, worker)
}

func GetVehicles(c *gin.Context) {
	var vehicles []models.Vehicle
	if err := database.DB.Find(&vehicles).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "查询车辆列表失败: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, vehicles)
}

func CreateVehicle(c *gin.Context) {
	var req CreateVehicleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	vehicle := models.Vehicle{
		PlateNumber:    req.PlateNumber,
		VehicleType:    req.VehicleType,
		CapacityVolume: req.CapacityVolume,
		CapacityWeight: req.CapacityWeight,
	}

	if req.Status != "" {
		vehicle.Status = models.VehicleStatus(req.Status)
	} else {
		vehicle.Status = models.VehicleStatusAvailable
	}

	if err := database.DB.Create(&vehicle).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "新增车辆失败: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, vehicle)
}

func GetSchedules(c *gin.Context) {
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")

	if startDate == "" || endDate == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请提供 start_date 和 end_date 参数"})
		return
	}

	start, err := time.Parse("2006-01-02", startDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "start_date 格式错误，应为 YYYY-MM-DD"})
		return
	}

	end, err := time.Parse("2006-01-02", endDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "end_date 格式错误，应为 YYYY-MM-DD"})
		return
	}

	var schedules []models.Schedule
	if err := database.DB.Where("work_date BETWEEN ? AND ?", start, end).
		Preload("Worker").
		Preload("Vehicle").
		Order("work_date ASC, id ASC").
		Find(&schedules).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "查询排班列表失败: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, schedules)
}

func CreateSchedule(c *gin.Context) {
	var req CreateScheduleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	workDate, err := time.Parse("2006-01-02", req.WorkDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "日期格式错误，应为 YYYY-MM-DD"})
		return
	}

	schedule := models.Schedule{
		WorkerID:  req.WorkerID,
		VehicleID: req.VehicleID,
		WorkDate: workDate,
		Shift:     req.Shift,
	}

	if req.Status != "" {
		schedule.Status = models.ScheduleStatus(req.Status)
	} else {
		schedule.Status = models.ScheduleStatusScheduled
	}

	if err := database.DB.Create(&schedule).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建排班失败: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, schedule)
}
