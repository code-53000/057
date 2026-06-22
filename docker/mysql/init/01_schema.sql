-- 搬家公司数据库建表脚本
-- 创建时间: 2026-06-22

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- 订单表
-- ----------------------------
DROP TABLE IF EXISTS `orders`;
CREATE TABLE `orders` (
  `id` BIGINT NOT NULL AUTO_INCREMENT COMMENT '订单ID',
  `customer_name` VARCHAR(50) NOT NULL COMMENT '客户姓名',
  `customer_phone` VARCHAR(20) NOT NULL COMMENT '客户电话',
  `move_date` DATETIME NOT NULL COMMENT '搬家日期',
  `start_address` VARCHAR(255) NOT NULL COMMENT '起始地址',
  `start_floor` INT NOT NULL DEFAULT 1 COMMENT '起始楼层',
  `start_has_elevator` TINYINT(1) NOT NULL DEFAULT 0 COMMENT '起始地址是否有电梯',
  `end_address` VARCHAR(255) NOT NULL COMMENT '目的地地址',
  `end_floor` INT NOT NULL DEFAULT 1 COMMENT '目的地楼层',
  `end_has_elevator` TINYINT(1) NOT NULL DEFAULT 0 COMMENT '目的地是否有电梯',
  `items_volume` FLOAT NOT NULL DEFAULT 0 COMMENT '物品体积(立方米)',
  `items_description` TEXT COMMENT '物品描述',
  `estimated_workers` INT NOT NULL DEFAULT 2 COMMENT '预估需要师傅人数',
  `estimated_vehicle_type` VARCHAR(20) NOT NULL COMMENT '预估车辆类型(小面/金杯/厢货)',
  `status` VARCHAR(20) NOT NULL DEFAULT 'pending' COMMENT '订单状态(pending待派单/dispatched已派单/in_progress进行中/completed已完成/cancelled已取消)',
  `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  KEY `idx_status` (`status`),
  KEY `idx_move_date` (`move_date`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='订单表';

-- ----------------------------
-- 师傅表
-- ----------------------------
DROP TABLE IF EXISTS `workers`;
CREATE TABLE `workers` (
  `id` BIGINT NOT NULL AUTO_INCREMENT COMMENT '师傅ID',
  `name` VARCHAR(50) NOT NULL COMMENT '师傅姓名',
  `phone` VARCHAR(20) NOT NULL COMMENT '师傅电话',
  `status` VARCHAR(20) NOT NULL DEFAULT 'available' COMMENT '状态(on_leave休假/available在岗)',
  `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_phone` (`phone`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='师傅表';

-- ----------------------------
-- 车辆表
-- ----------------------------
DROP TABLE IF EXISTS `vehicles`;
CREATE TABLE `vehicles` (
  `id` BIGINT NOT NULL AUTO_INCREMENT COMMENT '车辆ID',
  `plate_number` VARCHAR(20) NOT NULL COMMENT '车牌号',
  `vehicle_type` VARCHAR(20) NOT NULL COMMENT '车辆类型(小面/金杯/厢货)',
  `capacity_volume` FLOAT NOT NULL DEFAULT 0 COMMENT '容量体积(立方米)',
  `capacity_weight` FLOAT NOT NULL DEFAULT 0 COMMENT '载重(吨)',
  `status` VARCHAR(20) NOT NULL DEFAULT 'available' COMMENT '状态(available可用/maintenance维修中)',
  `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_plate_number` (`plate_number`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='车辆表';

-- ----------------------------
-- 排班表
-- ----------------------------
DROP TABLE IF EXISTS `schedules`;
CREATE TABLE `schedules` (
  `id` BIGINT NOT NULL AUTO_INCREMENT COMMENT '排班ID',
  `worker_id` BIGINT NOT NULL COMMENT '师傅ID',
  `vehicle_id` BIGINT NOT NULL COMMENT '车辆ID',
  `work_date` DATE NOT NULL COMMENT '工作日期',
  `shift` VARCHAR(20) NOT NULL COMMENT '班次(morning早班/afternoon午班/full_day全天)',
  `status` VARCHAR(20) NOT NULL DEFAULT 'scheduled' COMMENT '状态(scheduled已排班/working工作中/completed已完成/off休息)',
  `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  PRIMARY KEY (`id`),
  KEY `idx_worker_date` (`worker_id`, `work_date`),
  KEY `idx_vehicle_date` (`vehicle_id`, `work_date`),
  KEY `idx_work_date` (`work_date`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='排班表';

-- ----------------------------
-- 派单表
-- ----------------------------
DROP TABLE IF EXISTS `dispatches`;
CREATE TABLE `dispatches` (
  `id` BIGINT NOT NULL AUTO_INCREMENT COMMENT '派单ID',
  `order_id` BIGINT NOT NULL COMMENT '订单ID',
  `worker_ids` JSON NOT NULL COMMENT '分配的师傅ID列表(JSON数组)',
  `vehicle_id` BIGINT NOT NULL COMMENT '分配的车辆ID',
  `scheduled_start_time` DATETIME NOT NULL COMMENT '计划开始时间',
  `scheduled_end_time` DATETIME NOT NULL COMMENT '计划结束时间',
  `status` VARCHAR(20) NOT NULL DEFAULT 'pending' COMMENT '派单状态(pending待确认/accepted已接单/in_progress进行中/completed已完成)',
  `actual_start_time` DATETIME COMMENT '实际开始时间',
  `actual_end_time` DATETIME COMMENT '实际结束时间',
  `actual_workers_count` INT COMMENT '实际参与师傅人数',
  `actual_items_volume` FLOAT COMMENT '实际物品体积(立方米)',
  `remark` TEXT COMMENT '备注',
  `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  KEY `idx_order_id` (`order_id`),
  KEY `idx_status` (`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='派单表';

SET FOREIGN_KEY_CHECKS = 1;
