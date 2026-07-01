package repository

import (
	"errors"

	"gorm.io/gorm"
)

var (
	ErrNotFound     = errors.New("record not found")
	ErrDuplicateKey = errors.New("duplicate key")
)

func mapGORMError(err error) error {
	switch err {
	case nil:
		return nil
	case gorm.ErrRecordNotFound:
		return ErrNotFound
	case gorm.ErrDuplicatedKey:
		return ErrDuplicateKey
	}

	return err
}
