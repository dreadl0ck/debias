package debias

import "time"

type Stats struct {
	FileName string
	BytesIn int
	BytesOut int
	Duration time.Duration
}
