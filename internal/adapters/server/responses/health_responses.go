package responses

// HealthResponse represents the structure of the health check response.
type HealthResponse struct {
	Idle              string `json:"idle" example:"1"`
	InUse             string `json:"in_use" example:"0"`
	MaxIdleClosed     string `json:"max_idle_closed" example:"0"`
	MaxLifetimeClosed string `json:"max_lifetime_closed" example:"0"`
	Message           string `json:"message" example:"It's healthy'"`
	OpenConnections   string `json:"open_connections" example:"1"`
	Status            string `json:"status" example:"up"`
	WaitCount         string `json:"wait_count" example:"0"`
	WaitDuration      string `json:"wait_duration" example:"0s"`
}
