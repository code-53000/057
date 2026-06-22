package handlers

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"moving-schedule-backend/database"
	"moving-schedule-backend/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type CreateOrderRequest struct {
	CustomerName         string  `json:"customer_name" binding:"required"`
	CustomerPhone        string  `json:"customer_phone" binding:"required"`
	MoveDate             string  `json:"move_date" binding:"required"`
	StartAddress         string  `json:"start_address" binding:"required"`
	StartFloor           int     `json:"start_floor"`
	StartHasElevator     bool    `json:"start_has_elevator"`
	EndAddress           string  `json:"end_address" binding:"required"`
	EndFloor             int     `json:"end_floor"`
	EndHasElevator       bool    `json:"end_has_elevator"`
	ItemsVolume          float64 `json:"items_volume"`
	ItemsDescription     string  `json:"items_description"`
	EstimatedWorkers     int     `json:"estimated_workers"`
	EstimatedVehicleType string  `json:"estimated_vehicle_type" binding:"required"`
}

type UpdateOrderRequest struct {
	CustomerName         string   `json:"customer_name"`
	CustomerPhone        string   `json:"customer_phone"`
	MoveDate             string   `json:"move_date"`
	StartAddress         string   `json:"start_address"`
	StartFloor           *int     `json:"start_floor"`
	StartHasElevator     *bool    `json:"start_has_elevator"`
	EndAddress           string   `json:"end_address"`
	EndFloor             *int     `json:"end_floor"`
	EndHasElevator       *bool    `json:"end_has_elevator"`
	ItemsVolume          *float64 `json:"items_volume"`
	ItemsDescription     string   `json:"items_description"`
	EstimatedWorkers     *int     `json:"estimated_workers"`
	EstimatedVehicleType string   `json:"estimated_vehicle_type"`
	Status               string   `json:"status"`
}

func parseDateFlexible(dateStr string) (time.Time, error) {
	formats := []string{
		"2006-01-02",
		"2006-01-02 15:04:05",
		time.RFC3339,
		"2006/01/02",
		"2006-01-02T15:04:05.000Z",
		"2006-01-02T15:04:05Z",
	}
	for _, format := range formats {
		t, err := time.ParseInLocation(format, dateStr, time.Local)
		if err == nil {
			return t, nil
		}
	}
	return time.Time{}, nil
}

func parseDateFromAny(v interface{}) (time.Time, error) {
	var dateStr string
	switch val := v.(type) {
	case string:
		dateStr = val
	default:
		return time.Time{}, nil
	}

	if dateStr == "" {
		return time.Time{}, nil
	}

	formats := []string{
		"2006-01-02",
		"2006-01-02 15:04:05",
		time.RFC3339,
		"2006/01/02",
		"2006-01-02T15:04:05.000Z",
		"2006-01-02T15:04:05Z",
		time.RFC3339Nano,
	}
	for _, format := range formats {
		t, err := time.ParseInLocation(format, dateStr, time.Local)
		if err == nil {
			year, month, day := t.Date()
			return time.Date(year, month, day, 0, 0, 0, 0, time.Local), nil
		}
	}
	return time.Time{}, nil
}

func CreateOrder(c *gin.Context) {
	var req CreateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	moveDate, err := parseDateFlexible(req.MoveDate)
	if err != nil || moveDate.IsZero() {
		c.JSON(http.StatusBadRequest, gin.H{"error": "日期格式错误，应为 YYYY-MM-DD"})
		return
	}

	order := models.Order{
		CustomerName:         req.CustomerName,
		CustomerPhone:        req.CustomerPhone,
		MoveDate:             moveDate,
		StartAddress:         req.StartAddress,
		StartFloor:           req.StartFloor,
		StartHasElevator:     req.StartHasElevator,
		EndAddress:           req.EndAddress,
		EndFloor:             req.EndFloor,
		EndHasElevator:       req.EndHasElevator,
		ItemsVolume:          req.ItemsVolume,
		ItemsDescription:     req.ItemsDescription,
		EstimatedWorkers:     req.EstimatedWorkers,
		EstimatedVehicleType: req.EstimatedVehicleType,
		Status:               models.OrderStatusPending,
	}

	if req.StartFloor == 0 {
		order.StartFloor = 1
	}
	if req.EndFloor == 0 {
		order.EndFloor = 1
	}
	if req.EstimatedWorkers == 0 {
		order.EstimatedWorkers = 2
	}

	if err := database.DB.Create(&order).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建订单失败: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, order)
}

func GetOrders(c *gin.Context) {
	status := c.Query("status")
	date := c.Query("date")
	moveDate := c.Query("move_date")

	pageStr := c.DefaultQuery("page", "1")
	sizeStr := c.DefaultQuery("size", "10")

	page, _ := strconv.Atoi(pageStr)
	size, _ := strconv.Atoi(sizeStr)
	if page < 1 {
		page = 1
	}
	if size < 1 {
		size = 10
	}
	offset := (page - 1) * size

	query := database.DB.Model(&models.Order{})

	if status != "" {
		query = query.Where("status = ?", status)
	}

	if date != "" {
		query = query.Where("DATE(move_date) = ?", date)
	} else if moveDate != "" {
		query = query.Where("DATE(move_date) = ?", strings.Split(moveDate, "T")[0])
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "查询订单总数失败: " + err.Error()})
		return
	}

	var orders []models.Order
	if err := query.Order("move_date DESC, id DESC").
		Limit(size).
		Offset(offset).
		Find(&orders).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "查询订单列表失败: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"list":  orders,
		"total": total,
		"page":  page,
		"size":  size,
	})
}

func GetOrder(c *gin.Context) {
	id := c.Param("id")

	var order models.Order
	if err := database.DB.First(&order, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "订单不存在"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "查询订单失败: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, order)
}

func UpdateOrder(c *gin.Context) {
	id := c.Param("id")

	var order models.Order
	if err := database.DB.First(&order, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "订单不存在"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "查询订单失败: " + err.Error()})
		return
	}

	var req UpdateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.CustomerName != "" {
		order.CustomerName = req.CustomerName
	}
	if req.CustomerPhone != "" {
		order.CustomerPhone = req.CustomerPhone
	}
	if req.MoveDate != "" {
		moveDate, err := parseDateFlexible(req.MoveDate)
		if err != nil || moveDate.IsZero() {
			c.JSON(http.StatusBadRequest, gin.H{"error": "日期格式错误，应为 YYYY-MM-DD"})
			return
		}
		order.MoveDate = moveDate
	}
	if req.StartAddress != "" {
		order.StartAddress = req.StartAddress
	}
	if req.StartFloor != nil {
		order.StartFloor = *req.StartFloor
	}
	if req.StartHasElevator != nil {
		order.StartHasElevator = *req.StartHasElevator
	}
	if req.EndAddress != "" {
		order.EndAddress = req.EndAddress
	}
	if req.EndFloor != nil {
		order.EndFloor = *req.EndFloor
	}
	if req.EndHasElevator != nil {
		order.EndHasElevator = *req.EndHasElevator
	}
	if req.ItemsVolume != nil {
		order.ItemsVolume = *req.ItemsVolume
	}
	if req.ItemsDescription != "" {
		order.ItemsDescription = req.ItemsDescription
	}
	if req.EstimatedWorkers != nil {
		order.EstimatedWorkers = *req.EstimatedWorkers
	}
	if req.EstimatedVehicleType != "" {
		order.EstimatedVehicleType = req.EstimatedVehicleType
	}
	if req.Status != "" {
		order.Status = models.OrderStatus(req.Status)
	}

	if err := database.DB.Save(&order).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新订单失败: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, order)
}

func DeleteOrder(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的订单ID"})
		return
	}

	result := database.DB.Model(&models.Order{}).Where("id = ?", id).Update("status", models.OrderStatusCancelled)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "取消订单失败: " + result.Error.Error()})
		return
	}

	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "订单不存在"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "订单已取消"})
}
