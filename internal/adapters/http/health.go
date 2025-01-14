package http

import (
	"encoding/json"
	"fmt"
	"go-starter/internal/adapters/storage/postgres"
	"log/slog"
	"net/http"
)

// HealthHandler represents the HTTP handler for database health
type HealthHandler struct {
}

// NewHealthHandler creates a new HealthHandler instance
func NewHealthHandler() *HealthHandler {
	return &HealthHandler{}
}

// Health godoc
//
//	@Summary		Get database health information
//	@Description	Get database health information
//	@Tags			Health
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	healthResponse	"DB information"
//	@Failure		500	{object}	errorResponse	"Internal server error"
//	@Router			/v1/health [get]
func (hh *HealthHandler) Health(w http.ResponseWriter, _ *http.Request) {
	resp, err := json.Marshal(postgres.Health())
	if err != nil {
		http.Error(w, "failed to marshal health check response", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if _, err := w.Write(resp); err != nil {
		errMsg := fmt.Sprintf("failed to write health check response: %s", err)
		slog.Error(errMsg)
		http.Error(w, errMsg, http.StatusInternalServerError)
	}
}
