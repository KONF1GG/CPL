package handler

import (
	"CPL/internal/service"
	"errors"
	"net/http"
)

func writeServiceError(w http.ResponseWriter, err error) {
	status, message := mapServiceError(err)
	writeError(w, status, message)
}

func mapServiceError(err error) (int, string) {
	switch {
	case errors.Is(err, service.ErrVMNotFound):
		return http.StatusNotFound, "vm not found"
	case errors.Is(err, service.ErrTaskNotFound):
		return http.StatusNotFound, "task not found"
	case errors.Is(err, service.ErrVMNameTaken):
		return http.StatusConflict, "vm name already exists"
	case errors.Is(err, service.ErrInvalidVMConfig):
		return http.StatusBadRequest, err.Error()
	case errors.Is(err, service.ErrVMAlreadyRunning):
		return http.StatusConflict, "vm already running"
	case errors.Is(err, service.ErrVMAlreadyStopped):
		return http.StatusConflict, "vm already stopped"
	case errors.Is(err, service.ErrVMNotReady):
		return http.StatusConflict, "vm is not ready"
	default:
		return http.StatusInternalServerError, "internal server error"
	}
}
