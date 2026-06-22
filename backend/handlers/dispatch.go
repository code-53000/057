package handlers

import (
	"net/http"
	"strconv"
	"time"

	"moving-schedule-backend/database"
	"moving-schedule-backend/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type CreateDispatchRequest struct {
	OrderID            uint   `json:"order_id" binding:"required"`
	WorkerIDs          []uint `json:"worker_ids" binding:"required,min=1"`
	VehicleID          uint   `json:"vehicle_id" binding:"required"`
	ScheduledStartTime string `json:"scheduled_start_time" binding:"required"`
	ScheduledEndTime   string `json:"scheduled_end_time" binding:"required"`
}

func isTimeOverlap(start1, end1, start2, end2 time.Time) bool {
	return start1.Before(end2) && start2.Before(end1)
}

func containsWorkerID(workerIDs models.UintArray, workerID uint) bool {
	for _, id := range workerIDs {
		if id == workerID {
			return true
		}
	}
	return false
}

func checkWorkerScheduleConflict(workerID uint, date time.Time) *models.ConflictDetail {
	var schedule models.Schedule
	startOfDay := date.Truncate(24 * time.Hour)
	endOfDay := startOfDay.Add(24 * time.Hour).Add(-time.Second)

	err := database.DB.Where("worker_id = ? AND work_date >= ? AND work_date <= ? AND status != ?",
		workerID, startOfDay, endOfDay, models.ScheduleStatusOff).
		First(&schedule).Error
	if err == gorm.ErrRecordNotFound {
		var worker models.Worker
		database.DB.First(&worker, workerID)
		return &models.ConflictDetail{
			Type:        "worker_not_scheduled",
			Description: "师傅当天未排班在岗",
			WorkerID:    &workerID,
			WorkerName:  worker.Name,
		}
	}
	return nil
}

func checkWorkerTimeConflict(workerID uint, startTime, endTime time.Time, excludeDispatchID *uint) *models.ConflictDetail {
	var dispatches []models.Dispatch
	query := database.DB.Where("status != ? AND status != ?",
		models.DispatchStatusRejected,
		models.DispatchStatusCompleted)

	if excludeDispatchID != nil {
		query = query.Where("id != ?", *excludeDispatchID)
	}

	query.Find(&dispatches)

	for _, d := range dispatches {
		if containsWorkerID(d.WorkerIDs, workerID) && isTimeOverlap(startTime, endTime, d.ScheduledStartTime, d.ScheduledEndTime) {
			var worker models.Worker
			database.DB.First(&worker, workerID)
			conflictOrderID := d.OrderID
			return &models.ConflictDetail{
				Type:            "worker_time_conflict",
				Description:     "师傅该时间段已有其他派单",
				WorkerID:        &workerID,
				WorkerName:      worker.Name,
				ConflictOrderID: &conflictOrderID,
			}
		}
	}
	return nil
}

func checkVehicleTimeConflict(vehicleID uint, startTime, endTime time.Time, excludeDispatchID *uint) *models.ConflictDetail {
	var dispatches []models.Dispatch
	query := database.DB.Where("vehicle_id = ? AND status != ? AND status != ?",
		vehicleID,
		models.DispatchStatusRejected,
		models.DispatchStatusCompleted)

	if excludeDispatchID != nil {
		query = query.Where("id != ?", *excludeDispatchID)
	}

	query.Find(&dispatches)

	for _, d := range dispatches {
		if isTimeOverlap(startTime, endTime, d.ScheduledStartTime, d.ScheduledEndTime) {
			var vehicle models.Vehicle
			database.DB.First(&vehicle, vehicleID)
			conflictOrderID := d.OrderID
			return &models.ConflictDetail{
				Type:            "vehicle_time_conflict",
				Description:     "车辆该时间段已有其他派单",
				VehicleID:       &vehicleID,
				VehiclePlate:    vehicle.PlateNumber,
				ConflictOrderID: &conflictOrderID,
			}
		}
	}
	return nil
}

func checkCapacityOverload(vehicleID uint, itemsVolume float64) *models.ConflictDetail {
	var vehicle models.Vehicle
	if err := database.DB.First(&vehicle, vehicleID).Error; err != nil {
		return &models.ConflictDetail{
			Type:        "vehicle_not_found",
			Description: "车辆不存在",
			VehicleID:   &vehicleID,
		}
	}

	if itemsVolume > vehicle.CapacityVolume {
		return &models.ConflictDetail{
			Type:         "capacity_overload",
			Description:  "订单物品方量超过车辆容量",
			VehicleID:    &vehicleID,
			VehiclePlate: vehicle.PlateNumber,
			OverloadedBy: itemsVolume - vehicle.CapacityVolume,
		}
	}
	return nil
}

func CheckDispatchConflicts(orderID uint, workerIDs []uint, vehicleID uint, startTime, endTime time.Time, excludeDispatchID *uint) models.ConflictCheckResult {
	result := models.ConflictCheckResult{
		HasConflict: false,
		Conflicts:   make([]models.ConflictDetail, 0),
	}

	var order models.Order
	if err := database.DB.First(&order, orderID).Error; err != nil {
		result.HasConflict = true
		result.Conflicts = append(result.Conflicts, models.ConflictDetail{
			Type:        "order_not_found",
			Description: "订单不存在",
		})
		return result
	}

	for _, workerID := range workerIDs {
		if conflict := checkWorkerScheduleConflict(workerID, startTime); conflict != nil {
			result.HasConflict = true
			result.Conflicts = append(result.Conflicts, *conflict)
		}
	}

	for _, workerID := range workerIDs {
		if conflict := checkWorkerTimeConflict(workerID, startTime, endTime, excludeDispatchID); conflict != nil {
			result.HasConflict = true
			result.Conflicts = append(result.Conflicts, *conflict)
		}
	}

	if conflict := checkVehicleTimeConflict(vehicleID, startTime, endTime, excludeDispatchID); conflict != nil {
		result.HasConflict = true
		result.Conflicts = append(result.Conflicts, *conflict)
	}

	if conflict := checkCapacityOverload(vehicleID, order.ItemsVolume); conflict != nil {
		result.HasConflict = true
		result.Conflicts = append(result.Conflicts, *conflict)
	}

	return result
}

func CheckDispatch(c *gin.Context) {
	orderIDStr := c.Query("order_id")
	workerIDsStr := c.Query("worker_ids")
	vehicleIDStr := c.Query("vehicle_id")
	scheduledStartTimeStr := c.Query("scheduled_start_time")
	scheduledEndTimeStr := c.Query("scheduled_end_time")

	if orderIDStr == "" || workerIDsStr == "" || vehicleIDStr == "" || scheduledStartTimeStr == "" || scheduledEndTimeStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "缺少必要参数: order_id, worker_ids, vehicle_id, scheduled_start_time, scheduled_end_time"})
		return
	}

	orderID, err := strconv.ParseUint(orderIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "order_id 格式错误"})
		return
	}

	workerIDs, err := parseWorkerIDsQuery(workerIDsStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "worker_ids 格式错误"})
		return
	}

	vehicleID, err := strconv.ParseUint(vehicleIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "vehicle_id 格式错误"})
		return
	}

	startTime, err := time.Parse("2006-01-02 15:04:05", scheduledStartTimeStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "scheduled_start_time 格式错误，应为 YYYY-MM-DD HH:MM:SS"})
		return
	}

	endTime, err := time.Parse("2006-01-02 15:04:05", scheduledEndTimeStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "scheduled_end_time 格式错误，应为 YYYY-MM-DD HH:MM:SS"})
		return
	}

	result := CheckDispatchConflicts(uint(orderID), workerIDs, uint(vehicleID), startTime, endTime, nil)
	c.JSON(http.StatusOK, result)
}

func parseWorkerIDsQuery(workerIDsStr string) ([]uint, error) {
	var result []uint
	var current string
	for _, ch := range workerIDsStr {
		if ch == ',' {
			if current != "" {
				id, err := strconv.ParseUint(current, 10, 32)
				if err != nil {
					return nil, err
				}
				result = append(result, uint(id))
				current = ""
			}
		} else {
			current += string(ch)
		}
	}
	if current != "" {
		id, err := strconv.ParseUint(current, 10, 32)
		if err != nil {
			return nil, err
		}
		result = append(result, uint(id))
	}
	return result, nil
}

func CreateDispatch(c *gin.Context) {
	var req CreateDispatchRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	startTime, err := time.Parse("2006-01-02 15:04:05", req.ScheduledStartTime)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "scheduled_start_time 格式错误，应为 YYYY-MM-DD HH:MM:SS"})
		return
	}

	endTime, err := time.Parse("2006-01-02 15:04:05", req.ScheduledEndTime)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "scheduled_end_time 格式错误，应为 YYYY-MM-DD HH:MM:SS"})
		return
	}

	checkResult := CheckDispatchConflicts(req.OrderID, req.WorkerIDs, req.VehicleID, startTime, endTime, nil)
	if checkResult.HasConflict {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":     "派单冲突校验失败",
			"conflicts": checkResult.Conflicts,
		})
		return
	}

	dispatch := models.Dispatch{
		OrderID:            req.OrderID,
		WorkerIDs:          models.UintArray(req.WorkerIDs),
		VehicleID:          req.VehicleID,
		ScheduledStartTime: startTime,
		ScheduledEndTime:   endTime,
		Status:             models.DispatchStatusPending,
	}

	if err := database.DB.Create(&dispatch).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建派单失败: " + err.Error()})
		return
	}

	database.DB.Model(&models.Order{}).Where("id = ?", req.OrderID).Update("status", models.OrderStatusDispatched)

	c.JSON(http.StatusCreated, dispatch)
}

func GetDispatches(c *gin.Context) {
	var dispatches []models.Dispatch
	if err := database.DB.Preload("Order").Preload("Vehicle").
		Order("scheduled_start_time DESC, id DESC").
		Find(&dispatches).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "查询派单列表失败: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, dispatches)
}

func GetDispatch(c *gin.Context) {
	id := c.Param("id")

	var dispatch models.Dispatch
	if err := database.DB.Preload("Order").Preload("Vehicle").First(&dispatch, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "派单不存在"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "查询派单失败: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, dispatch)
}
