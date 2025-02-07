package view

import (
	"sync"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/google/uuid"
	"github.com/jneo8/jujuspell/data"
	"github.com/jneo8/jujuspell/model"
)

type controllerResourceViewerKeyMap struct {
	Enter key.Binding
	ESC   key.Binding
}

func (m controllerResourceViewerKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{m.Enter, m.ESC},
	}
}

func (m controllerResourceViewerKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{
		m.Enter, m.ESC,
	}
}

var controllerResourceViewerKeys = controllerResourceViewerKeyMap{
	Enter: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "list models"),
	),
	ESC: key.NewBinding(
		key.WithKeys("esc"),
		key.WithHelp("esc", "focus"),
	),
}

type controllerResourceViewer struct {
	ID                   uuid.UUID
	table                table.Model
	QueryJob             model.QueryJob
	setCurrentController sync.Once
	keyMap               controllerResourceViewerKeyMap
}

func (viewer *controllerResourceViewer) Refresh(msg model.RefreshMsg) {
	controllerRefreshMsg := msg.(*data.ControllerRefreshMsg)
	viewer.table.SetColumns(controllerRefreshMsg.Columns)
	viewer.table.SetRows(controllerRefreshMsg.Rows)
	viewer.setCurrentController.Do(
		func() {
			for i, row := range controllerRefreshMsg.Rows {
				if row[0] == controllerRefreshMsg.CurrentController {
					viewer.table.SetCursor(i)
				}
			}
		},
	)
}

func (viewer *controllerResourceViewer) View() string {
	return viewer.table.View()
}

func (viewer *controllerResourceViewer) Update(msg tea.Msg) tea.Cmd {
	m, cmd := viewer.table.Update(msg)
	viewer.table = m
	return cmd
}

func (viewer *controllerResourceViewer) GetQueryJob() model.QueryJob {
	return viewer.QueryJob
}

func (viewer *controllerResourceViewer) GetKeyMap() help.KeyMap {
	return viewer.keyMap
}
