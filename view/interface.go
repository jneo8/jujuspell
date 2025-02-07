package view

import (
	"context"
	"sync"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/google/uuid"
	"github.com/jneo8/jujuspell/model"
)

type App interface {
	Exec() error
}

type Model interface {
	tea.Model
	QueryJobController
}

type QueryJobController interface {
	GetQueryJobs() map[uuid.UUID]model.QueryJob
}

type teaProgramWrapper interface {
	Send(msg tea.Msg)
}

type UIController interface {
	teaProgramWrapper
	QueryJobController
	RunProgram(
		ctx context.Context,
		cancel context.CancelFunc,
		wg *sync.WaitGroup,
		errCh chan error,
	)
}
