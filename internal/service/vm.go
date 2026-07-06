package service

import (
	"CPL/internal/models"
	"context"
	"fmt"
)

type VMRepository interface {
	Create(ctx context.Context, vm *models.VM) error
	GetByID(ctx context.Context, id uint) (*models.VM, error)
	GetByIDForUpdate(ctx context.Context, id uint) (*models.VM, error)
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
	return s.enqueueTask(ctx, id, models.TaskTypeStart, validateStart)
}

func (s *VMService) Stop(ctx context.Context, id uint) (*models.Task, error) {
	return s.enqueueTask(ctx, id, models.TaskTypeStop, validateStop)
}

func (s *VMService) Delete(ctx context.Context, id uint) (*models.Task, error) {
	return s.enqueueTask(ctx, id, models.TaskTypeDelete, validateDelete)
}

func (s *VMService) enqueueTask(
	ctx context.Context,
	vmID uint,
	taskType models.TaskType,
	validate func(*models.VM) error,
) (*models.Task, error) {
	var task *models.Task

	err := s.tx.WithinTransaction(ctx, func(txCtx context.Context) error {
		vm, err := s.vmRepo.GetByIDForUpdate(txCtx, vmID)
		if err != nil {
			return mapVMError(err)
		}
		if err := validate(vm); err != nil {
			return err
		}

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

func validateStart(vm *models.VM) error {
	switch vm.Status {
	case models.VMStatusRunning:
		return ErrVMAlreadyRunning
	case models.VMStatusPending:
		return ErrVMNotReady
	}
	return nil
}

func validateStop(vm *models.VM) error {
	switch vm.Status {
	case models.VMStatusStopped:
		return ErrVMAlreadyStopped
	case models.VMStatusPending:
		return ErrVMNotReady
	}
	return nil
}

func validateDelete(vm *models.VM) error {
	if vm.Status == models.VMStatusPending {
		return ErrVMNotReady
	}
	return nil
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
