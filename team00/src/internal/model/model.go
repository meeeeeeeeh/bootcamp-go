package model

import "time"

type Anomaly struct {
	SessionId string
	Frequency float64
	Dt        time.Time
}
