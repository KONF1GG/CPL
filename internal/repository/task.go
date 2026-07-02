package repository

import (
	"CPL/internal/models"
	"context"

	"gorm.io/gorm"
)

type TaskRepository struct {
	db *gorm.DB
}

func NewTaskRepo(db *gorm.DB) *TaskRepository {
	return &TaskRepository{db: db}
}

func (r *TaskRepository) Create(ctx context.Context, task *models.Task) error {
	err := DBFromContext(ctx, r.db).WithContext(ctx).Create(task).Error
	return mapGORMError(err)
}

func (r *TaskRepository) GetByID(ctx context.Context, id uint) (*models.Task, error) {
	var task models.Task
	err := DBFromContext(ctx, r.db).WithContext(ctx).First(&task, id).Error
	if err != nil {
		return nil, mapGORMError(err)
	}
	return &task, nil
}

func (r *TaskRepository) GetAll(ctx context.Context) ([]models.Task, error) {
	var tasks []models.Task
	err := DBFromContext(ctx, r.db).WithContext(ctx).Order("created_at DESC").Find(&tasks).Error
	return tasks, mapGORMError(err)
}

func (r *TaskRepository) ListByVMID(ctx context.Context, vmID uint) ([]models.Task, error) {
	var tasks []models.Task
	err := DBFromContext(ctx, r.db).WithContext(ctx).
		Where("vm_id = ?", vmID).
		Order("created_at DESC").
		Find(&tasks).Error
	return tasks, mapGORMError(err)
}

func (r *TaskRepository) Update(ctx context.Context, task *models.Task) error {
	result := DBFromContext(ctx, r.db).WithContext(ctx).Save(task)
	if result.Error != nil {
		return mapGORMError(result.Error)
	}
	if result.RowsAffected == 0 {
		return ErrNotFound
	}
	return nil
}
