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

type CompleteDispatchRequest struct {
	ActualStartTime   string  `json:"actual_start_time"`
	ActualEndTime     string  `json:"actual_end_time"`
	ActualWorkersCount *int   `json:"actual_workers_count"`
	ActualItemsVolume *float64 `json:"actual_items_volume"`
	Remark            string  `json:"remark"`
}

func GetWorkerDispatches(c *gin.Context) {
	workerIDStr := c.Query("worker_id")
	if workerIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请提供 worker_id 参数"})
		return
	}

	workerID, err := strconv.ParseUint(workerIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "worker_id 格式错误"})
		return
	}

	var allDispatches []models.Dispatch
	if err := database.DB.Preload("Order").Preload("Vehicle").
		Order("scheduled_start_time DESC, id DESC").
		Find(&allDispatches).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "查询派单列表失败: " + err.Error()})
		return
	}

	result := make([]models.Dispatch, 0)
	for _, d := range allDispatches {
		if containsWorkerID(d.WorkerIDs, uint(workerID)) {
			result = append(result, d)
		}
	}

	c.JSON(http.StatusOK, result)
}

func AcceptDispatch(c *gin.Context) {
	id := c.Param("id")

	var dispatch models.Dispatch
	if err := database.DB.First(&dispatch, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "派单不存在"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "查询派单失败: " + err.Error()})
		return
	}

	if dispatch.Status != models.DispatchStatusPending {
		c.JSON(http.StatusBadRequest, gin.H{"error": "当前派单状态不允许确认接单"})
		return
	}

	dispatch.Status = models.DispatchStatusAccepted
	if err := database.DB.Save(&dispatch).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "确认接单失败: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, dispatch)
}

func StartDispatch(c *gin.Context) {
	id := c.Param("id")

	var dispatch models.Dispatch
	if err := database.DB.First(&dispatch, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "派单不存在"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "查询派单失败: " + err.Error()})
		return
	}

	if dispatch.Status != models.DispatchStatusAccepted {
		c.JSON(http.StatusBadRequest, gin.H{"error": "当前派单状态不允许开始作业"})
		return
	}

	now := time.Now()
	dispatch.Status = models.DispatchStatusInProgress
	dispatch.ActualStartTime = &now

	if err := database.DB.Save(&dispatch).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "开始作业失败: " + err.Error()})
		return
	}

	database.DB.Model(&models.Order{}).Where("id = ?", dispatch.OrderID).Update("status", models.OrderStatusInProgress)

	c.JSON(http.StatusOK, dispatch)
}

func parseTimeStr(timeStr string) (*time.Time, error) {
	t, err := time.Parse(time.RFC3339, timeStr)
	if err == nil {
		return &t, nil
	}
	t2, err2 := time.Parse("2006-01-02 15:04:05", timeStr)
	if err2 == nil {
		return &t2, nil
	}
	return nil, err
}

func CompleteDispatch(c *gin.Context) {
	id := c.Param("id")

	var dispatch models.Dispatch
	if err := database.DB.First(&dispatch, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "派单不存在"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "查询派单失败: " + err.Error()})
		return
	}

	if dispatch.Status != models.DispatchStatusInProgress {
		c.JSON(http.StatusBadRequest, gin.H{"error": "当前派单状态不允许完成作业"})
		return
	}

	var req CompleteDispatchRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	dispatch.Status = models.DispatchStatusCompleted

	if req.ActualStartTime != "" {
		t, err := parseTimeStr(req.ActualStartTime)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "实际开始时间格式错误，应为 RFC3339 或 YYYY-MM-DD HH:MM:SS"})
			return
		}
		dispatch.ActualStartTime = t
	} else if dispatch.ActualStartTime == nil {
		now := time.Now()
		dispatch.ActualStartTime = &now
	}

	if req.ActualEndTime != "" {
		t, err := parseTimeStr(req.ActualEndTime)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "实际结束时间格式错误，应为 RFC3339 或 YYYY-MM-DD HH:MM:SS"})
			return
		}
		dispatch.ActualEndTime = t
	} else {
		now := time.Now()
		dispatch.ActualEndTime = &now
	}

	if req.ActualWorkersCount != nil {
		dispatch.ActualWorkersCount = req.ActualWorkersCount
	}
	if req.ActualItemsVolume != nil {
		dispatch.ActualItemsVolume = req.ActualItemsVolume
	}
	if req.Remark != "" {
		dispatch.Remark = req.Remark
	}

	if err := database.DB.Save(&dispatch).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "完成作业失败: " + err.Error()})
		return
	}

	database.DB.Model(&models.Order{}).Where("id = ?", dispatch.OrderID).Update("status", models.OrderStatusCompleted)

	c.JSON(http.StatusOK, dispatch)
}
