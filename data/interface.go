package data

import (
	"context"
	"sync"

	"github.com/google/uuid"
)

type CTR interface {
	AddJob(job QueryJob)
	Run(ctx context.Context, wg *sync.WaitGroup)
}

// 1 to 1 mapping to QueryJob
type DataFetcher interface {
	Fetch(queryJob QueryJob) (RefreshMsg, error)
}

type Worker interface {
	Run()
	Stop()
}

type JobWorker interface {
	DataFetcher
	Worker
}

type RefreshMsg interface {
	GetJobID() uuid.UUID
}

type baseRefreshMsg struct {
	queryJobID uuid.UUID
}

func (msg *baseRefreshMsg) GetJobID() uuid.UUID {
	return msg.queryJobID
}

type ResourceType string
type Filter string

type QueryJob struct {
	ID           uuid.UUID
	ResourceType ResourceType
	Filter       Filter
}
