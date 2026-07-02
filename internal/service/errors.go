package service

import (
	"CPL/internal/repository"
	"errors"
	"fmt"
)

var (
	ErrVMNotFound       = errors.New("vm not found")
	ErrVMNameTaken      = errors.New("vm name already exists")
	ErrInvalidVMConfig  = errors.New("invalid vm configuration")
	ErrVMAlreadyRunning = errors.New("vm already running")
	ErrVMAlreadyStopped = errors.New("vm already stopped")
	ErrVMNotReady       = errors.New("vm is not ready")
	ErrTaskNotFound     = errors.New("task not found")
)

func mapVMError(err error) error {
	switch {
	case errors.Is(err, repository.ErrNotFound):
		return ErrVMNotFound
	case errors.Is(err, repository.ErrDuplicateKey):
		return ErrVMNameTaken
	default:
		return fmt.Errorf("vm: %w", err)
	}
}

func mapTaskError(err error) error {
	switch {
	case errors.Is(err, repository.ErrNotFound):
		return ErrTaskNotFound
	default:
		return fmt.Errorf("task: %w", err)
	}
}
