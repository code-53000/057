package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"
)

type OrderStatus string

const (
	OrderStatusPending    OrderStatus = "pending"
	OrderStatusDispatched OrderStatus = "dispatched"
	OrderStatusInProgress OrderStatus = "in_progress"
	OrderStatusCompleted  OrderStatus = "completed"
	OrderStatusCancelled  OrderStatus = "cancelled"
)

type Order struct {
	ID                  uint           `gorm:"primaryKey" json:"id"`
	CustomerName        string         `gorm:"column:customer_name;type:varchar(50);not null" json:"customer_name"`
	CustomerPhone       string         `gorm:"column:customer_phone;type:varchar(20);not null" json:"customer_phone"`
	MoveDate            time.Time      `gorm:"column:move_date;not null" json:"move_date"`
	StartAddress        string         `gorm:"column:start_address;type:varchar(255);not null" json:"start_address"`
	StartFloor          int            `gorm:"column:start_floor;not null;default:1" json:"start_floor"`
	StartHasElevator    bool           `gorm:"column:start_has_elevator;not null;default:false" json:"start_has_elevator"`
	EndAddress          string         `gorm:"column:end_address;type:varchar(255);not null" json:"end_address"`
	EndFloor            int            `gorm:"column:end_floor;not null;default:1" json:"end_floor"`
	EndHasElevator      bool           `gorm:"column:end_has_elevator;not null;default:false" json:"end_has_elevator"`
	ItemsVolume         float64        `gorm:"column:items_volume;type:float;not null;default:0" json:"items_volume"`
	ItemsDescription    string         `gorm:"column:items_description;type:text" json:"items_description"`
	EstimatedWorkers    int            `gorm:"column:estimated_workers;not null;default:2" json:"estimated_workers"`
	EstimatedVehicleType string        `gorm:"column:estimated_vehicle_type;type:varchar(20);not null" json:"estimated_vehicle_type"`
	Status              OrderStatus `gorm:"column:status;type:varchar(20);not null;default:pending" json:"status"`
	CreatedAt           time.Time   `json:"created_at"`
	UpdatedAt           time.Time   `json:"updated_at"`
}

func (Order) TableName() string {
	return "orders"
}

type WorkerStatus string

const (
	WorkerStatusAvailable WorkerStatus = "available"
	WorkerStatusOnLeave   WorkerStatus = "on_leave"
)

type Worker struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Name      string         `gorm:"column:name;type:varchar(50);not null" json:"name"`
	Phone     string         `gorm:"column:phone;type:varchar(20);not null;unique" json:"phone"`
	Status    WorkerStatus `gorm:"column:status;type:varchar(20);not null;default:available" json:"status"`
	CreatedAt time.Time    `json:"created_at"`
}

func (Worker) TableName() string {
	return "workers"
}

type VehicleStatus string

const (
	VehicleStatusAvailable   VehicleStatus = "available"
	VehicleStatusMaintenance VehicleStatus = "maintenance"
)

type Vehicle struct {
	ID              uint           `gorm:"primaryKey" json:"id"`
	PlateNumber     string         `gorm:"column:plate_number;type:varchar(20);not null;unique" json:"plate_number"`
	VehicleType     string         `gorm:"column:vehicle_type;type:varchar(20);not null" json:"vehicle_type"`
	CapacityVolume  float64        `gorm:"column:capacity_volume;type:float;not null;default:0" json:"capacity_volume"`
	CapacityWeight  float64        `gorm:"column:capacity_weight;type:float;not null;default:0" json:"capacity_weight"`
	Status          VehicleStatus `gorm:"column:status;type:varchar(20);not null;default:available" json:"status"`
	CreatedAt       time.Time     `json:"created_at"`
}

func (Vehicle) TableName() string {
	return "vehicles"
}

type ScheduleStatus string

const (
	ScheduleStatusScheduled ScheduleStatus = "scheduled"
	ScheduleStatusWorking   ScheduleStatus = "working"
	ScheduleStatusCompleted ScheduleStatus = "completed"
	ScheduleStatusOff       ScheduleStatus = "off"
)

type Schedule struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	WorkerID  uint           `gorm:"column:worker_id;not null" json:"worker_id"`
	VehicleID uint           `gorm:"column:vehicle_id;not null" json:"vehicle_id"`
	WorkDate  time.Time      `gorm:"column:work_date;type:date;not null" json:"work_date"`
	Shift     string         `gorm:"column:shift;type:varchar(20);not null" json:"shift"`
	Status    ScheduleStatus `gorm:"column:status;type:varchar(20);not null;default:scheduled" json:"status"`
	CreatedAt time.Time      `json:"created_at"`
	Worker    Worker         `gorm:"foreignKey:WorkerID" json:"worker,omitempty"`
	Vehicle   Vehicle        `gorm:"foreignKey:VehicleID" json:"vehicle,omitempty"`
}

func (Schedule) TableName() string {
	return "schedules"
}

type DispatchStatus string

const (
	DispatchStatusPending    DispatchStatus = "pending"
	DispatchStatusAccepted   DispatchStatus = "accepted"
	DispatchStatusInProgress DispatchStatus = "in_progress"
	DispatchStatusCompleted  DispatchStatus = "completed"
	DispatchStatusRejected   DispatchStatus = "rejected"
)

type UintArray []uint

func (a UintArray) Value() (driver.Value, error) {
	if a == nil {
		return nil, nil
	}
	return json.Marshal([]uint(a))
}

func (a *UintArray) Scan(value interface{}) error {
	if value == nil {
		*a = nil
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("failed to scan UintArray")
	}
	var result []uint
	if err := json.Unmarshal(bytes, &result); err != nil {
		return err
	}
	*a = UintArray(result)
	return nil
}

type Dispatch struct {
	ID               uint           `gorm:"primaryKey" json:"id"`
	OrderID          uint           `gorm:"column:order_id;not null" json:"order_id"`
	WorkerIDs        UintArray      `gorm:"column:worker_ids;type:json;not null" json:"worker_ids"`
	VehicleID        uint           `gorm:"column:vehicle_id;not null" json:"vehicle_id"`
	ScheduledStartTime time.Time    `gorm:"column:scheduled_start_time;not null" json:"scheduled_start_time"`
	ScheduledEndTime   time.Time    `gorm:"column:scheduled_end_time;not null" json:"scheduled_end_time"`
	Status           DispatchStatus `gorm:"column:status;type:varchar(20);not null;default:pending" json:"status"`
	ActualStartTime  *time.Time     `gorm:"column:actual_start_time" json:"actual_start_time"`
	ActualEndTime    *time.Time     `gorm:"column:actual_end_time" json:"actual_end_time"`
	ActualWorkersCount *int         `gorm:"column:actual_workers_count" json:"actual_workers_count"`
	ActualItemsVolume *float64      `gorm:"column:actual_items_volume;type:float" json:"actual_items_volume"`
	Remark           string         `gorm:"column:remark;type:text" json:"remark"`
	CreatedAt        time.Time      `json:"created_at"`
	UpdatedAt        time.Time      `json:"updated_at"`
	Order            Order          `gorm:"foreignKey:OrderID" json:"order,omitempty"`
	Vehicle          Vehicle        `gorm:"foreignKey:VehicleID" json:"vehicle,omitempty"`
}

func (Dispatch) TableName() string {
	return "dispatches"
}

type ConflictDetail struct {
	Type            string  `json:"type"`
	Description     string  `json:"description"`
	WorkerName      string  `json:"worker_name,omitempty"`
	WorkerID        *uint   `json:"worker_id,omitempty"`
	VehiclePlate    string  `json:"vehicle_plate,omitempty"`
	VehicleID       *uint   `json:"vehicle_id,omitempty"`
	ConflictOrderID *uint   `json:"conflict_order_id,omitempty"`
	OverloadedBy    float64 `json:"overloaded_by,omitempty"`
}

type ConflictCheckResult struct {
	HasConflict bool             `json:"has_conflict"`
	Conflicts   []ConflictDetail `json:"conflicts"`
}
