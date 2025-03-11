package handlers

import (
	"encoding/json"
	"go-starter/internal/adapters/server/responses"
	"go-starter/internal/adapters/storage/database"
	"go-starter/internal/domain"
	"net/http"
)

// HealthHandler is responsible for handling HTTP requests related to the health status of the database.
type HealthHandler struct{}

// NewHealthHandler initializes and returns a new instance of HealthHandler.
func NewHealthHandler() *HealthHandler {
	return &HealthHandler{}
}

// PostgresHealth godoc
//
//	@Summary		Get database health information
//	@Description	Get database health information
//	@Tags			Health
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	responses.HealthResponse	"Postgres health information"
//	@Failure		500	{object}	responses.ErrorResponse	"Internal server error"
//	@Router			/v1/health/postgres [get]
func (hh *HealthHandler) PostgresHealth(w http.ResponseWriter, _ *http.Request) {
	resp, err := json.Marshal(database.Health())
	if err != nil {
		responses.HandleError(w, domain.ErrInternal)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if _, err := w.Write(resp); err != nil {
		responses.HandleError(w, domain.ErrInternal)
		return
	}
}
