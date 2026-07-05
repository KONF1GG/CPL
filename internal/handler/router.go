package handler

import "net/http"

func NewRouter(vmHandler *VMHandler, taskHandler *TaskHandler) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
	})

	mux.HandleFunc("POST /vms", vmHandler.Create)
	mux.HandleFunc("GET /vms", vmHandler.List)
	mux.HandleFunc("GET /vms/{id}", vmHandler.GetByID)
	mux.HandleFunc("POST /vms/{id}/start", vmHandler.Start)
	mux.HandleFunc("POST /vms/{id}/stop", vmHandler.Stop)
	mux.HandleFunc("DELETE /vms/{id}", vmHandler.Delete)

	mux.HandleFunc("GET /tasks/{id}", taskHandler.GetByID)
	mux.HandleFunc("GET /vms/{id}/tasks", taskHandler.ListByVMID)

	return mux
}
