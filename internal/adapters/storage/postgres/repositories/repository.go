package repositories

import (
	"time"
)

var (
	// QueryTimeoutDuration specifies the timeout duration for database queries
	QueryTimeoutDuration = time.Second * 5
)
