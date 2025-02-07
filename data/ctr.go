package data

import (
	"github.com/google/uuid"
	"github.com/jneo8/jujuspell/jujuclient"
	"github.com/jneo8/jujuspell/model"
	"github.com/rs/zerolog/log"
)

type ctr struct {
	workers map[uuid.UUID]JobWorker
	client  jujuclient.Client
}

func NewCTR(client jujuclient.Client) CTR {
	return &ctr{
		client:  client,
		workers: make(map[uuid.UUID]JobWorker),
	}
}

func (c *ctr) GetJobWorker(job model.QueryJob) JobWorker {
	defer log.Debug().Str("job id", job.ID.String()).Msg("Get job worker")
	// If worker already exists
	if worker, ok := c.workers[job.ID]; ok {
		return worker
	}

	// Register new worker
	switch job.ResourceType {
	case ControllerResourceType:
		worker := NewControllerJobWorker(c.client)
		c.workers[job.ID] = worker
		return worker
	}
	return nil
}
