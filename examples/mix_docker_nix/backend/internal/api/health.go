package api

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

// registerHealthHandlers binds operational liveness checks under the versioned API.
func (a *API) registerHealthHandlers(r *mux.Router) *mux.Router {
	r.HandleFunc("/health", a.getHealth).Methods("GET")

	return r
}

// getHealth is a lightweight liveness endpoint used by Compose and external checks.
func (a *API) getHealth(w http.ResponseWriter, r *http.Request) {
	a.logger.Info("Health handler reached")

	if err := a.redis.Ping(r.Context()); err != nil {
		a.logger.Error("Redis health check failed", zap.Error(err))
		a.writeJSON(w, http.StatusServiceUnavailable, map[string]any{
			"status": "unhealthy",
			"redis":  "unavailable",
		})
		return
	}

	a.writeJSON(w, http.StatusOK, map[string]any{
		"status": "ok",
		"redis":  "ok",
		"time":   time.Now().Format(time.RFC3339),
	})
}

// writeJSON writes the common JSON response shape and records failures that can no longer
// be returned after the HTTP status has been committed.
func (a *API) writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(value); err != nil {
		a.logger.Error("encode HTTP response", zap.Int("status", status), zap.Error(err))
	}
}
