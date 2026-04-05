package models

import "time"

type ResponseStats struct {
	Count           int
	CountCompleted  int
	PerCompleted    float64
	AvgTimeComplete time.Duration
}

const (
	KeyFormatALl  = "app:stats:all"
	KeyFormatUser = "app:stats:user:" // + {id}
	KeyFormatTime = "app:stats:time:" // + {time_start|time_end}
)
