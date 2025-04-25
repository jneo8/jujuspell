package view

import (
	"context"
	"sync"

	"github.com/charmbracelet/bubbles/help"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/google/uuid"
	"github.com/jneo8/jujuspell/data"
)

type App interface {
	Exec() error
}

type Model interface {
	tea.Model
	QueryJobController
}

type QueryJobController interface {
	GetQueryJobs() map[uuid.UUID]data.QueryJob
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

type ResourceViewer interface {
	GetID() uuid.UUID
	Refresh(data.RefreshMsg)
	View() string
	Update(msg tea.Msg) tea.Cmd
	GetQueryJob() data.QueryJob
	GetKeyMap() help.KeyMap
}
