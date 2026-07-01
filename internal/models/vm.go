package models

import (
	"time"

	"gorm.io/gorm"
)

type VMStatus string

const (
	VMStatusPending = "pending"
	VMStatusStopped = "stopped"
	VMStatusRunning = "running"
)

type VM struct {
	ID        uint           `gorm:"PrimaryKey" json:"id"`
	Name      string         `gorm:"not null;size:255;uniqueIndex:idx_vms_name_active,where:deleted_at IS NULL" json:"name"`
	Status    VMStatus       `gorm:"not null;default:pending;size:50" json:"status"`
	CPU       int            `gorm:"not null" json:"cpu"`
	RamMB     int            `gorm:"not null" json:"ram_mb"`
	DiskGB    int            `gorm:"not null" json:"disk_gb"`
	TaskID    uint           `gorm:"index" json:"task_id,omitempty"`
	Task      *Task          `gorm:"foreignKey:TaskID" json:"-"`
	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}
