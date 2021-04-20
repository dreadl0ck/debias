package debias

import "time"

type Stats struct {
	FileName string
	BytesIn  int64
	BytesOut int64
	Duration time.Duration
}
