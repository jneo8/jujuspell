package data

import (
	"github.com/charmbracelet/bubbles/table"
	"github.com/google/uuid"
	"github.com/jneo8/jujuspell/jujuclient"
	"github.com/jneo8/jujuspell/model"
	"github.com/rs/zerolog/log"
)

const (
	ControllerResourceType model.ResourceType = "Controller"
)

func NewControllerJobWorker(
	client jujuclient.Client,
) JobWorker {
	return &controllerJobWorker{client: client}
}

type controllerJobWorker struct {
	client jujuclient.Client
}

func (c *controllerJobWorker) Fetch(queryJob model.QueryJob) (model.RefreshMsg, error) {
	rows, errs := c.client.GetControllerData()
	if len(errs) != 0 {
		for _, err := range errs {
			log.Warn().Err(err)
		}
	}
	currentController, err := c.client.GetCurrentController()
	if err != nil {
		return nil, err
	}
	return c.getRefreshMsg(
		queryJob.ID,
		rows,
		c.getColumns(),
		currentController,
	), nil
}

func (c *controllerJobWorker) getColumns() []table.Column {
	return []table.Column{
		{Title: "Controller", Width: 15},
		{Title: "Model", Width: 15},
		{Title: "User", Width: 5},
		{Title: "Access", Width: 10},
		{Title: "Cloud/Region", Width: 10},
		{Title: "Models", Width: 5},
		{Title: "Nodes", Width: 5},
		{Title: "HA", Width: 5},
		{Title: "Version", Width: 10},
	}
}

type ControllerRefreshMsg struct {
	*baseRefreshMsg
	Rows              []table.Row
	Columns           []table.Column
	CurrentController string
}

func (c *controllerJobWorker) getRefreshMsg(
	queryJobID uuid.UUID,
	rows []table.Row,
	columns []table.Column,
	currentController string,
) model.RefreshMsg {
	return &ControllerRefreshMsg{
		baseRefreshMsg: &baseRefreshMsg{
			queryJobID: queryJobID,
		},
		Rows:              rows,
		Columns:           columns,
		CurrentController: currentController,
	}
}
