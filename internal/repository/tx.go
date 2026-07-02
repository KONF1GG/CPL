package repository

import (
	"context"

	"gorm.io/gorm"
)

type TransactionManager struct {
	db *gorm.DB
}

func NewTransactionManager(db *gorm.DB) *TransactionManager {
	return &TransactionManager{db: db}
}

func (m *TransactionManager) WithinTransaction(
	ctx context.Context,
	fn func(ctx context.Context) error,
) error {
	return m.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return fn(WithTx(ctx, tx))
	})
}
