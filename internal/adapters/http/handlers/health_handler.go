package handlers

import (
	"encoding/json"
	"fmt"
	_ "go-starter/internal/adapters/http/responses"
	"go-starter/internal/adapters/storage/postgres"
	"log/slog"
	"net/http"
)

// HealthHandler is responsible for handling HTTP requests related to the health status of the database.
type HealthHandler struct {
}

// NewHealthHandler initializes and returns a new instance of HealthHandler.
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
//	@Success		200	{object}	responses.HealthResponse	"DB information"
//	@Failure		500	{object}	responses.ErrorResponse	"Internal server error"
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
