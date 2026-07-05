package handler

import (
	"CPL/internal/models"
	"encoding/json"
	"net/http"
	"strconv"
	"time"
)

type VMResponse struct {
	ID        uint   `json:"id"`
	Name      string `json:"name"`
	Status    string `json:"status"`
	CPU       int    `json:"cpu"`
	RamMB     int    `json:"ram_mb"`
	DiskGB    int    `json:"disk_gb"`
	TaskID    uint   `json:"task_id,omitempty"`
	CreatedAt string `json:"created_at"`
}
type TaskResponse struct {
	ID     uint   `json:"id"`
	Type   string `json:"type"`
	Status string `json:"status"`
	VMID   uint   `json:"vm_id"`
}
type ErrorResponse struct {
	Error string `json:"error"`
}

func toVMResponse(vm *models.VM) VMResponse {
	return VMResponse{
		ID:        vm.ID,
		Name:      vm.Name,
		Status:    string(vm.Status),
		CPU:       vm.CPU,
		RamMB:     vm.RamMB,
		DiskGB:    vm.DiskGB,
		TaskID:    vm.TaskID,
		CreatedAt: vm.CreatedAt.Format(time.RFC3339),
	}
}

func toVMResponses(vms []models.VM) []VMResponse {
	out := make([]VMResponse, len(vms))
	for i := range vms {
		out[i] = toVMResponse(&vms[i])
	}
	return out
}

func toTaskResponse(task *models.Task) TaskResponse {
	return TaskResponse{
		ID:     task.ID,
		Type:   string(task.Type),
		Status: string(task.Status),
		VMID:   task.VMID,
	}
}

func toTaskResponses(tasks []models.Task) []TaskResponse {
	out := make([]TaskResponse, len(tasks))
	for i := range tasks {
		out[i] = toTaskResponse(&tasks[i])
	}
	return out
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func decodeJSON(r *http.Request, dst any) error {
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	return dec.Decode(dst)
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, ErrorResponse{Error: msg})
}

func parseID(s string) (uint, error) {
	id, err := strconv.ParseUint(s, 10, 64)
	return uint(id), err
}
