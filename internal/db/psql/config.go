package psql

import "time"

type Config struct {
	ConnectionString string
	Timeout          time.Duration
}
