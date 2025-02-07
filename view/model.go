package view

import (
	"sync"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/google/uuid"
	"github.com/jneo8/jujuspell/data"
	"github.com/jneo8/jujuspell/model"
	"github.com/rs/zerolog/log"
)

type rootModel struct {
	currentTab      uuid.UUID
	resourceViewers map[uuid.UUID]ResourceViewer
}

func getRootModel() Model {
	m := rootModel{
		resourceViewers: make(map[uuid.UUID]ResourceViewer),
	}
	return &m
}

func (m *rootModel) Init() tea.Cmd {
	return func() tea.Msg {
		return AddResourceViewerMsg{ResourceType: data.ControllerResourceType}
	}
}

func (m *rootModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	log.Debug().Interface("Msg", msg).Type("type", msg).Msg("Msg")

	var cmd tea.Cmd
	switch msg := msg.(type) {
	case model.RefreshMsg:
		if viewer, ok := m.resourceViewers[msg.GetJobID()]; ok {
			viewer.Refresh(msg)
		}
	case AddResourceViewerMsg:
		m.addResourceViewer(msg)
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			// if m.getCurrentTab().Focused() {
			// 	m.getCurrentTab().Blur()
			// } else {
			// 	m.getCurrentTab().Focus()
			// }
		case "q", "ctrl+c":
			return m, tea.Quit
		case "enter":
			log.Debug().Msg("Enter")
		default:
			if currentTab := m.getCurrentTab(); currentTab != nil {
				cmd = currentTab.Update(msg)
			}
		}
	}

	return m, cmd
}

func (m *rootModel) View() string {
	var views []string
	for _, viewer := range m.resourceViewers {
		views = append(views, viewer.View())
	}
	return baseStyle.Render(lipgloss.JoinVertical(lipgloss.Left, views...))
}

func (m *rootModel) getCurrentTab() ResourceViewer {
	return m.resourceViewers[m.currentTab]
}

func (m *rootModel) GetQueryJobs() map[uuid.UUID]model.QueryJob {
	jobs := make(map[uuid.UUID]model.QueryJob)
	for id, viewer := range m.resourceViewers {
		jobs[id] = viewer.GetQueryJob()
	}
	return jobs
}

func (m *rootModel) addResourceViewer(msg AddResourceViewerMsg) {
	id := uuid.New()
	log.Debug().Str("resource type", string(msg.ResourceType)).Msg("Add resource viewer")
	switch msg.ResourceType {
	case data.ControllerResourceType:
		m.resourceViewers[id] = &controllerResourceViewer{
			ID: id,
			Model: table.New(
				table.WithFocused(true),
			),
			QueryJob: model.QueryJob{
				ID:           id,
				ResourceType: msg.ResourceType,
				Filter:       "",
			},
		}
	}
	m.currentTab = id
}

type AddResourceViewerMsg struct{ ResourceType model.ResourceType }

type ResourceViewer interface {
	Refresh(model.RefreshMsg)
	View() string
	Update(msg tea.Msg) tea.Cmd
	GetQueryJob() model.QueryJob
}

type controllerResourceViewer struct {
	ID                   uuid.UUID
	Model                table.Model
	QueryJob             model.QueryJob
	setCurrentController sync.Once
}

func (viewer *controllerResourceViewer) Refresh(msg model.RefreshMsg) {
	controllerRefreshMsg := msg.(*data.ControllerRefreshMsg)
	viewer.Model.SetColumns(controllerRefreshMsg.Columns)
	viewer.Model.SetRows(controllerRefreshMsg.Rows)
	viewer.setCurrentController.Do(
		func() {
			for i, row := range controllerRefreshMsg.Rows {
				if row[0] == controllerRefreshMsg.CurrentController {
					viewer.Model.SetCursor(i)
				}
			}
		},
	)
}

func (viewer *controllerResourceViewer) View() string {
	return viewer.Model.View()
}

func (viewer *controllerResourceViewer) Update(msg tea.Msg) tea.Cmd {
	m, cmd := viewer.Model.Update(msg)
	viewer.Model = m
	return cmd
}

func (viewer *controllerResourceViewer) GetQueryJob() model.QueryJob {
	return viewer.QueryJob
}
