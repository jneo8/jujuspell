package view

import (
	"context"
	"sync"

	"github.com/charmbracelet/lipgloss"
	"github.com/jneo8/jujuspell/data"
	"github.com/jneo8/jujuspell/jujuclient"
	"github.com/rs/zerolog/log"
)

var (
	baseStyle = lipgloss.NewStyle()
)

type app struct {
	ui        UIController
	client    jujuclient.Client
	refresher Refresher
}

func NewApp(client jujuclient.Client) App {
	ui := NewUI()
	return &app{
		ui: ui,
		refresher: NewRefresher(
			ui,
			data.NewCTR(client),
		),
		client: client,
	}
}

func (a *app) Exec() error {
	var wg sync.WaitGroup
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	errCh := make(chan error, 2)
	defer close(errCh)

	wg.Add(2)
	go a.refresher.Run(ctx, &wg)
	go a.ui.RunProgram(ctx, cancel, &wg, errCh)

	select {
	case <-ctx.Done():
	case err := <-errCh:
		log.Error().Err(err).Msg("Get error")
		cancel()
	}

	wg.Wait()
	return nil
}
