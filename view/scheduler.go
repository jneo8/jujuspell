package view

import (
	"context"
	"sync"
	"time"

	"github.com/jneo8/jujuspell/data"
	"github.com/rs/zerolog/log"
)

type Scheduler interface {
	Run(ctx context.Context, wg *sync.WaitGroup)
}

func NewScheduler(
	uiController UIController,
	ctr data.CTR,
) Scheduler {
	return &scheduler{
		uiController: uiController,
		ctr:          ctr,
	}
}

type scheduler struct {
	uiController UIController
	ctr          data.CTR
}

func (s *scheduler) Run(
	ctx context.Context, wg *sync.WaitGroup,
) {
	defer wg.Done()
	for {
		log.Debug().Msg("Scheduler loop")
		select {
		case <-ctx.Done():
			return
		default:
			for _, job := range s.uiController.GetQueryJobs() {
				s.ctr.AddJob(job)
				// worker := s.ctr.GetJobWorker(job)
				// if worker == nil {
				// 	log.Warn().Str("resource type", string(job.ResourceType)).Msg("Unknown resource type")
				// 	continue
				// }
				// msg, err := worker.Fetch(job)
				// if err != nil {
				// 	log.Error().Err(err)
				// 	continue
				// }
				// s.uiController.Send(msg)
			}
		}
		time.Sleep(3 * time.Second)
	}
}
