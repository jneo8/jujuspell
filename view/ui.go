package view

import (
	"context"
	"sync"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/google/uuid"
	"github.com/jneo8/jujuspell/model"
	"github.com/rs/zerolog/log"
)

type ui struct {
	model   Model
	program *tea.Program
}

func NewUI() UIController {
	model := getRootModel()
	p := tea.NewProgram(model)
	return &ui{
		model:   model,
		program: p,
	}
}

func (u *ui) RunProgram(
	ctx context.Context,
	cancel context.CancelFunc,
	wg *sync.WaitGroup,
	errCh chan error,
) {
	defer wg.Done()
	go func() {
		defer cancel()
		if _, err := u.program.Run(); err != nil {
			errCh <- err
		}
	}()

	<-ctx.Done()
	u.program.Quit()
	log.Debug().Msg("Quit program")
}

func (u *ui) Send(msg tea.Msg) {
	log.Debug().Msg("ui Send msg")
	u.program.Send(msg)
	log.Debug().Msg("ui finish Send msg")
}

func (u *ui) GetQueryJobs() map[uuid.UUID]model.QueryJob {
	return u.model.GetQueryJobs()
}
