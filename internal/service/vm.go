package service

import (
	"CPL/internal/models"
	"context"
	"fmt"
)

type VMRepository interface {
	Create(ctx context.Context, vm *models.VM) error
	GetByID(ctx context.Context, id uint) (*models.VM, error)
	GetAll(ctx context.Context) ([]models.VM, error)
	Update(ctx context.Context, vm *models.VM) error
	Delete(ctx context.Context, id uint) error
}

type TxManager interface {
	WithinTransaction(ctx context.Context, fn func(ctx context.Context) error) error
}

type CreateVMInput struct {
	Name   string
	CPU    int
	RamMB  int
	DiskGB int
}

type VMService struct {
	vmRepo   VMRepository
	taskRepo TaskRepository
	tx       TxManager
}

func NewVMService(vmRepo VMRepository, taskRepo TaskRepository, tx TxManager) *VMService {
	return &VMService{
		vmRepo:   vmRepo,
		taskRepo: taskRepo,
		tx:       tx,
	}
}

func (s *VMService) Create(ctx context.Context, input CreateVMInput) (*models.VM, error) {
	if err := validateCreateInput(input); err != nil {
		return nil, err
	}

	var result *models.VM

	err := s.tx.WithinTransaction(ctx, func(txCtx context.Context) error {
		vm := &models.VM{
			Name:   input.Name,
			Status: models.VMStatusPending,
			CPU:    input.CPU,
			RamMB:  input.RamMB,
			DiskGB: input.DiskGB,
		}

		if err := s.vmRepo.Create(txCtx, vm); err != nil {
			return mapVMError(err)
		}

		task := &models.Task{
			Type:   models.TaskTypeProvision,
			Status: models.TaskStatusPending,
			VMID:   vm.ID,
		}
		if err := s.taskRepo.Create(txCtx, task); err != nil {
			return mapTaskError(err)
		}

		vm.TaskID = task.ID
		if err := s.vmRepo.Update(txCtx, vm); err != nil {
			return mapVMError(err)
		}

		result = vm
		return nil
	})
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (s *VMService) GetByID(ctx context.Context, id uint) (*models.VM, error) {
	vm, err := s.vmRepo.GetByID(ctx, id)
	if err != nil {
		return nil, mapVMError(err)
	}
	return vm, nil
}

func (s *VMService) List(ctx context.Context) ([]models.VM, error) {
	vms, err := s.vmRepo.GetAll(ctx)
	if err != nil {
		return nil, mapVMError(err)
	}
	return vms, nil
}

func (s *VMService) Start(ctx context.Context, id uint) (*models.Task, error) {
	vm, err := s.vmRepo.GetByID(ctx, id)
	if err != nil {
		return nil, mapVMError(err)
	}

	switch vm.Status {
	case models.VMStatusRunning:
		return nil, ErrVMAlreadyRunning
	case models.VMStatusPending:
		return nil, ErrVMNotReady
	}

	return s.enqueueTask(ctx, vm, models.TaskTypeStart)
}

func (s *VMService) Stop(ctx context.Context, id uint) (*models.Task, error) {
	vm, err := s.vmRepo.GetByID(ctx, id)
	if err != nil {
		return nil, mapVMError(err)
	}

	switch vm.Status {
	case models.VMStatusStopped:
		return nil, ErrVMAlreadyStopped
	case models.VMStatusPending:
		return nil, ErrVMNotReady
	}

	return s.enqueueTask(ctx, vm, models.TaskTypeStop)
}

func (s *VMService) Delete(ctx context.Context, id uint) (*models.Task, error) {
	vm, err := s.vmRepo.GetByID(ctx, id)
	if err != nil {
		return nil, mapVMError(err)
	}

	return s.enqueueTask(ctx, vm, models.TaskTypeDelete)
}

func (s *VMService) enqueueTask(ctx context.Context, vm *models.VM, taskType models.TaskType) (*models.Task, error) {
	var task *models.Task

	err := s.tx.WithinTransaction(ctx, func(txCtx context.Context) error {
		t := &models.Task{
			Type:   taskType,
			Status: models.TaskStatusPending,
			VMID:   vm.ID,
		}
		if err := s.taskRepo.Create(txCtx, t); err != nil {
			return mapTaskError(err)
		}

		vm.TaskID = t.ID
		if err := s.vmRepo.Update(txCtx, vm); err != nil {
			return mapVMError(err)
		}

		task = t
		return nil
	})
	if err != nil {
		return nil, err
	}

	return task, nil
}

func validateCreateInput(input CreateVMInput) error {
	if input.Name == "" {
		return fmt.Errorf("%w: name is required", ErrInvalidVMConfig)
	}
	if input.CPU <= 0 || input.RamMB <= 0 || input.DiskGB <= 0 {
		return fmt.Errorf("%w: cpu, ram_mb, disk_gb must be positive", ErrInvalidVMConfig)
	}
	return nil
}
