package data

import (
	"context"
	"sync"

	"github.com/google/uuid"
	"github.com/jneo8/jujuspell/jujuclient"
)

type ctr struct {
	workers map[uuid.UUID]JobWorker
	client  jujuclient.Client
	jobStatus []
}

func NewCTR(client jujuclient.Client) CTR {
	return &ctr{
		client:  client,
		workers: make(map[uuid.UUID]JobWorker),
	}
}

func (c *ctr) AddJob(job QueryJob) {
	// If worker already exists
	if _, ok := c.workers[job.ID]; ok {
		return
	}

	// Register new worker
	switch job.ResourceType {
	case ControllerResourceType:
		worker := NewControllerJobWorker(c.client)
		c.workers[job.ID] = worker
	}

}

func (c *ctr) Run(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()
}
