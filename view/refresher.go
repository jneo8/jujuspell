package view

import (
	"context"
	"sync"
	"time"

	"github.com/jneo8/jujuspell/data"
	"github.com/rs/zerolog/log"
)

type Refresher interface {
	Run(ctx context.Context, wg *sync.WaitGroup)
}

func NewRefresher(
	uiController UIController,
	ctr data.CTR,
) Refresher {
	return &refresher{
		uiController: uiController,
		ctr:          ctr,
	}
}

type refresher struct {
	uiController UIController
	ctr          data.CTR
}

func (r *refresher) Run(
	ctx context.Context, wg *sync.WaitGroup,
) {
	defer wg.Done()
	for {
		log.Debug().Msg("Refresher loop")
		select {
		case <-ctx.Done():
			return
		default:
			for _, job := range r.uiController.GetQueryJobs() {
				worker := r.ctr.GetJobWorker(job)
				if worker == nil {
					log.Warn().Str("resource type", string(job.ResourceType)).Msg("Unknown resource type")
					continue
				}
				msg, err := worker.Fetch(job)
				if err != nil {
					log.Error().Err(err)
					continue
				}
				r.uiController.Send(msg)
			}
		}
		time.Sleep(3 * time.Second)
	}
}
