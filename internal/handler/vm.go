package handler

import (
	"CPL/internal/models"
	"CPL/internal/service"
	"context"
	"net/http"
)

type VMService interface {
	Create(ctx context.Context, input service.CreateVMInput) (*models.VM, error)
	GetByID(ctx context.Context, id uint) (*models.VM, error)
	List(ctx context.Context) ([]models.VM, error)
	Start(ctx context.Context, id uint) (*models.Task, error)
	Stop(ctx context.Context, id uint) (*models.Task, error)
	Delete(ctx context.Context, id uint) (*models.Task, error)
}

type VMHandler struct {
	svc VMService
}

func NewVMHandler(svc VMService) *VMHandler {
	return &VMHandler{svc: svc}
}

func (h *VMHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req CreateVMRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid json")
		return
	}

	vm, err := h.svc.Create(r.Context(), service.CreateVMInput{
		Name:   req.Name,
		CPU:    req.CPU,
		RamMB:  req.RamMB,
		DiskGB: req.DiskGB,
	})

	if err != nil {
		writeServiceError(w, err)
		return
	}

	writeJSON(w, http.StatusAccepted, toVMResponse(vm))
}

func (h *VMHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r.PathValue("id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid id")
		return
	}

	vm, err := h.svc.GetByID(r.Context(), id)
	if err != nil {
		writeServiceError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, toVMResponse(vm))

}

func (h *VMHandler) Start(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r.PathValue("id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid id")
		return
	}

	task, err := h.svc.Start(r.Context(), id)
	if err != nil {
		writeServiceError(w, err)
		return
	}

	writeJSON(w, http.StatusAccepted, toTaskResponse(task))
}

func (h *VMHandler) List(w http.ResponseWriter, r *http.Request) {
	vms, err := h.svc.List(r.Context())
	if err != nil {
		writeServiceError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, toVMResponses(vms))
}

func (h *VMHandler) Stop(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r.PathValue("id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid id")
		return
	}

	task, err := h.svc.Stop(r.Context(), id)
	if err != nil {
		writeServiceError(w, err)
		return
	}

	writeJSON(w, http.StatusAccepted, toTaskResponse(task))
}

func (h *VMHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r.PathValue("id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid id")
		return
	}

	task, err := h.svc.Delete(r.Context(), id)
	if err != nil {
		writeServiceError(w, err)
		return
	}

	writeJSON(w, http.StatusAccepted, toTaskResponse(task))
}
