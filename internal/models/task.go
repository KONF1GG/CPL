package models

import "time"

type TaskType string

const (
	TaskTypeProvision TaskType = "provision"
	TaskTypeStart     TaskType = "start"
	TaskTypeStop      TaskType = "stop"
	TaskTypeDelete    TaskType = "delete"
)

type TaskStatus string

const (
	TaskStatusPending   TaskStatus = "pending"
	TaskStatusRunning   TaskStatus = "running"
	TaskStatusCompleted TaskStatus = "completed"
	TaskStatusFailed    TaskStatus = "failed"
)

type Task struct {
	ID        uint       `gorm:"primaryKey" json:"id"`
	Type      TaskType   `gorm:"not null;size:50" json:"type"`
	Status    TaskStatus `gorm:"not null;default:pending;size:50" json:"status"`
	VMID      uint       `gorm:"not null;index" json:"vm_id"`
	VM        *VM        `gorm:"foreignKey:VMID" json:"-"`
	CreatedAt time.Time  `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time  `gorm:"autoUpdateTime" json:"updated_at"`
}
