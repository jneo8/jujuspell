package view

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/google/uuid"
	"github.com/jneo8/jujuspell/data"
	"github.com/rs/zerolog/log"
)

type AddResourceViewerMsg struct{ ResourceType data.ResourceType }

type rootModel struct {
	currentTab      uuid.UUID
	resourceViewers map[uuid.UUID]ResourceViewer
	help            help.Model
}

type rootModelKeyMap struct {
	Help key.Binding
	Quit key.Binding
}

var rootModelKeys = rootModelKeyMap{
	Help: key.NewBinding(
		key.WithKeys("?"),
		key.WithHelp("?", "toggle help"),
	),
	Quit: key.NewBinding(
		key.WithKeys("q", "ctrl+c"),
		key.WithHelp("q / ctrl+c", "quit"),
	),
}

func (m rootModelKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{m.Help, m.Quit},
	}
}

func (m rootModelKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{
		m.Help, m.Quit,
	}
}

func getRootModel() Model {
	m := rootModel{
		resourceViewers: make(map[uuid.UUID]ResourceViewer),
		help:            help.New(),
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
	currentTab := m.getCurrentTab()
	log.Debug().Interface("CurrentTab", currentTab).Msg("CurrentTab")

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, rootModelKeys.Quit):
			return m, tea.Quit
		case key.Matches(msg, rootModelKeys.Help):
			m.help.ShowAll = !m.help.ShowAll
		default:
			if currentTab != nil {
				cmd = currentTab.Update(msg)
			}
		}
	case AddResourceViewerMsg:
		m.addResourceViewer(msg)
	case data.RefreshMsg:
		if currentTab != nil {
			currentTab.Refresh(msg)
		}
	default:
		if currentTab != nil {
			cmd = currentTab.Update(msg)
		}
	}

	return m, cmd
}

func (m *rootModel) View() string {
	var views []string
	for _, viewer := range m.resourceViewers {
		views = append(views, viewer.View())
	}
	helpView := m.getHelpView()

	return baseStyle.Render(
		lipgloss.JoinVertical(
			lipgloss.Top,
			[]string{
				lipgloss.JoinVertical(lipgloss.Left, views...),
				helpView,
			}...,
		),
	)
}

func (m *rootModel) getHelpView() (helpView string) {
	currentTab := m.getCurrentTab()
	if !m.help.ShowAll {
		keyBindings := rootModelKeys.ShortHelp()
		if currentTab != nil {
			keyBindings = append(
				currentTab.GetKeyMap().ShortHelp(),
				keyBindings...,
			)
		}
		helpView = m.help.ShortHelpView(keyBindings)
	} else {
		keyBindings := rootModelKeys.FullHelp()
		if currentTab != nil {
			keyBindings = append(
				currentTab.GetKeyMap().FullHelp(),
				keyBindings...,
			)
		}
		helpView = m.help.FullHelpView(keyBindings)
	}
	return helpView
}

func (m *rootModel) getCurrentTab() ResourceViewer {
	return m.resourceViewers[m.currentTab]
}

func (m *rootModel) GetQueryJobs() map[uuid.UUID]data.QueryJob {
	jobs := make(map[uuid.UUID]data.QueryJob)
	for id, viewer := range m.resourceViewers {
		jobs[id] = viewer.GetQueryJob()
	}
	return jobs
}

func (m *rootModel) addResourceViewer(msg AddResourceViewerMsg) {
	log.Debug().Str("resource type", string(msg.ResourceType)).Msg("Add resource viewer")
	switch msg.ResourceType {
	case data.ControllerResourceType:
		viewer := newControllerResourceViewer()
		m.resourceViewers[viewer.GetID()] = viewer
		m.currentTab = viewer.GetID()
	}
}
