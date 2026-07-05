package handler

import (
	"CPL/internal/service"
	"errors"
	"net/http"
)

func writeServiceError(w http.ResponseWriter, err error) {
	status := mapServiceError(err)
	writeError(w, status, err.Error())
}

func mapServiceError(err error) int {
	switch {
	case errors.Is(err, service.ErrVMNotFound),
		errors.Is(err, service.ErrTaskNotFound):
		return http.StatusNotFound
	case errors.Is(err, service.ErrVMNameTaken):
		return http.StatusConflict
	case errors.Is(err, service.ErrInvalidVMConfig):
		return http.StatusBadRequest
	case errors.Is(err, service.ErrVMAlreadyRunning),
		errors.Is(err, service.ErrVMAlreadyStopped),
		errors.Is(err, service.ErrVMNotReady):
		return http.StatusConflict
	default:
		return http.StatusInternalServerError
	}
}
