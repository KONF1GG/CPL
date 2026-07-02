package service

import (
	"CPL/internal/models"
	"context"
)

type TaskRepository interface {
	Create(ctx context.Context, task *models.Task) error
	GetByID(ctx context.Context, id uint) (*models.Task, error)
	Update(ctx context.Context, task *models.Task) error
	ListByVMID(ctx context.Context, vmID uint) ([]models.Task, error)
}

type TaskService struct {
	repo TaskRepository
}

func NewTaskService(repo TaskRepository) *TaskService {
	return &TaskService{repo: repo}
}

func (s *TaskService) GetByID(ctx context.Context, id uint) (*models.Task, error) {
	task, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, mapTaskError(err)
	}
	return task, nil
}

func (s *TaskService) ListByVMID(ctx context.Context, vmID uint) ([]models.Task, error) {
	tasks, err := s.repo.ListByVMID(ctx, vmID)
	if err != nil {
		return nil, mapTaskError(err)
	}
	return tasks, nil
}
