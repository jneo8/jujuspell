package jujuclient

import (
	"fmt"
	"strconv"

	"github.com/charmbracelet/bubbles/table"
	"github.com/juju/errors"
	"github.com/juju/juju/jujuclient"
	"github.com/rs/zerolog/log"
)

type Client interface {
	GetControllerData() ([]table.Row, []error)
	GetCurrentController() (string, error)
}

type jujuClient struct {
	clientStore jujuclient.ClientStore
}

func NewJujuClient(clientStore jujuclient.ClientStore) (Client, error) {
	return &jujuClient{
		clientStore: clientStore,
	}, nil
}

func (c *jujuClient) GetControllerData() ([]table.Row, []error) {
	allControllers, err := c.clientStore.AllControllers()
	if err != nil {
		return []table.Row{}, []error{err}
	}
	controllerItems, errs := convertControllerDetails(c.clientStore, allControllers)
	if len(errs) != 0 {
		return []table.Row{}, errs
	}

	rows := []table.Row{}
	for ctrlName, ctrl := range controllerItems {
		ha := "none"
		if ctrl.ControllerMachines != nil && ctrl.ControllerMachines.Total > 1 {
			ha = "yes"
		}
		modelName := ctrl.ModelName
		if modelName == "" {
			modelName = "-"
		}
		rows = append(
			rows,
			table.Row{
				ctrlName,
				modelName,
				ctrl.User,
				ctrl.Access,
				fmt.Sprintf("%s%s", ctrl.Cloud, ctrl.CloudRegion),
				strconv.Itoa(*ctrl.ModelCount),
				strconv.Itoa(*ctrl.MachineCount),
				ha,
				ctrl.AgentVersion,
			},
		)
	}
	return rows, []error{}
}

func (c *jujuClient) GetCurrentController() (string, error) {
	currentController, err := c.clientStore.CurrentController()
	if err != nil {
		if errors.Is(err, errors.NotFound) {
			log.Debug().Msg("CurrentController not found")
			return "", nil
		} else {
			return "", err
		}
	}
	return currentController, nil
}
