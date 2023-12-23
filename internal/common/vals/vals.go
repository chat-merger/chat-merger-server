package vals

import "time"

const (
	AUTHENTICATE_HEADER = "X-Api-Key"
)

var (
	GracefulShutdownTimeout = 2 * time.Second
)
