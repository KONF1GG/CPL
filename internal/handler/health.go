package handler

import (
	"context"
	"net/http"
	"time"
)

type HealthChecker interface {
	PingContext(ctx context.Context) error
}

func Health(checker HealthChecker) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
		defer cancel()

		if err := checker.PingContext(ctx); err != nil {
			writeError(w, http.StatusServiceUnavailable, "database unavailable")
			return
		}

		writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
	}
}
