package client

import (
	"team-00/internal/model"
)

type DbRepo interface {
	SaveAnomaly(anomaly model.Anomaly) error
}
