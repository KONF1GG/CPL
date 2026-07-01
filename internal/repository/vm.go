package repository

import (
	"CPL/internal/models"
	"context"

	"gorm.io/gorm"
)

type VMRepository struct {
	db *gorm.DB
}

func NewVMRepository(db *gorm.DB) VMRepository {
	return VMRepository{db: db}

}

func (r *VMRepository) Create(ctx context.Context, vm *models.VM) error {
	err := r.db.WithContext(ctx).Create(vm).Error
	return mapGORMError(err)
}

func (r *VMRepository) Update(ctx context.Context, vm *models.VM) error {
	result := r.db.WithContext(ctx).Save(vm)
	if result.Error != nil {
		return mapGORMError(result.Error)
	}
	if result.RowsAffected == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *VMRepository) Delete(ctx context.Context, id uint) error {
	result := r.db.WithContext(ctx).Delete(&models.VM{}, id)
	if result.Error != nil {
		return mapGORMError(result.Error)
	}
	if result.RowsAffected == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *VMRepository) GetByID(ctx context.Context, id uint) (*models.VM, error) {
	var vm models.VM
	err := r.db.WithContext(ctx).First(&vm, id).Error
	if err != nil {
		return nil, mapGORMError(err)
	}
	return &vm, nil
}

func (r *VMRepository) GetAll(ctx context.Context) ([]models.VM, error) {
	var vms []models.VM
	err := r.db.WithContext(ctx).Find(&vms).Error
	return vms, mapGORMError(err)
}
