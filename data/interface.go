package data

import (
	"github.com/google/uuid"
	"github.com/jneo8/jujuspell/model"
)

type CTR interface {
	GetJobWorker(job model.QueryJob) JobWorker
}

// 1 to 1 mapping to QueryJob
type JobWorker interface {
	Fetch(queryJob model.QueryJob) (model.RefreshMsg, error)
}

type baseRefreshMsg struct {
	queryJobID uuid.UUID
}

func (msg *baseRefreshMsg) GetJobID() uuid.UUID {
	return msg.queryJobID
}
