package handler

import (
	"CPL/internal/models"
	"context"
	"net/http"
)

type TaskService interface {
	GetByID(ctx context.Context, id uint) (*models.Task, error)
	ListByVMID(ctx context.Context, vmID uint) ([]models.Task, error)
}

type TaskHandler struct {
	svc TaskService
}

func NewTaskHandler(svc TaskService) *TaskHandler {
	return &TaskHandler{svc: svc}
}

func (h *TaskHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r.PathValue("id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid id")
		return
	}

	task, err := h.svc.GetByID(r.Context(), id)
	if err != nil {
		writeServiceError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, toTaskResponse(task))
}

func (h *TaskHandler) ListByVMID(w http.ResponseWriter, r *http.Request) {
	vmID, err := parseID(r.PathValue("id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid id")
		return
	}

	tasks, err := h.svc.ListByVMID(r.Context(), vmID)
	if err != nil {
		writeServiceError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, toTaskResponses(tasks))
}
