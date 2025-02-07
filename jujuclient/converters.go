package jujuclient

import (
	"github.com/juju/errors"
	k8sconstants "github.com/juju/juju/caas/kubernetes/provider/constants"
	"github.com/juju/juju/cmd/juju/common"
	jujucontroller "github.com/juju/juju/cmd/juju/controller"
	"github.com/juju/juju/jujuclient"
	names "github.com/juju/names/v6"
)

// This function is copy from https://github.com/juju/juju/blob/e6501cebd8719b55cfd8b56c386a0bc96104350d/cmd/juju/controller/listcontrollersconverters.go#L55
func convertControllerDetails(clientStore jujuclient.ClientStore, storeControllers map[string]jujuclient.ControllerDetails) (map[string]jujucontroller.ControllerItem, []error) {
	if len(storeControllers) == 0 {
		return nil, nil
	}
	errs := []error{}
	addError := func(msg, controllerName string, err error) {
		errs = append(errs, err)
	}
	controllers := map[string]jujucontroller.ControllerItem{}
	for controllerName, details := range storeControllers {
		serverName := ""
		// The most recently connected-to address
		// is the first in the list
		if len(details.APIEndpoints) > 0 {
			serverName = details.APIEndpoints[0]
		}

		var userName, access string
		accountDetails, err := clientStore.AccountDetails(controllerName)
		if err != nil {
			if !errors.IsNotFound(err) {
				addError("account details", controllerName, err)
				continue
			}
		} else {
			userName = accountDetails.User
			access = accountDetails.LastKnownAccess
		}

		var modelName string
		currentModel, err := clientStore.CurrentModel(controllerName)
		if err != nil {
			if !errors.IsNotFound(err) {
				addError("model", controllerName, err)
			}
		} else {
			modelName = currentModel
			if userName != "" {
				// There's a user loggedin, so display the
				// model name relative to that user.
				if unqualifiedModelName, owner, err := jujuclient.SplitModelName(modelName); err != nil {
					user := names.NewUserTag(userName)
					modelName = common.OwnerQualifiedModelName(unqualifiedModelName, owner, user)
				}
			}
		}

		models, err := clientStore.AllModels(controllerName)
		if err != nil && !errors.IsNotFound(err) {
			addError("models", controllerName, err)
		}
		modelCount := len(models)
		item := jujucontroller.ControllerItem{
			ModelName:      modelName,
			User:           userName,
			Access:         access,
			Server:         serverName,
			APIEndpoints:   details.APIEndpoints,
			ControllerUUID: details.ControllerUUID,
			CACert:         details.CACert,
			Cloud:          details.Cloud,
			CloudRegion:    details.CloudRegion,
			AgentVersion:   details.AgentVersion,
		}
		isCaas := details.CloudType == string(k8sconstants.StorageProviderType)
		if details.MachineCount != nil && *details.MachineCount > 0 {
			if isCaas {
				item.NodeCount = details.MachineCount
			} else {
				item.MachineCount = details.MachineCount
			}
		}
		if modelCount > 0 {
			item.ModelCount = &modelCount
		}
		if details.ControllerMachineCount > 0 {
			if isCaas {
				item.ControllerNodes = &jujucontroller.ControllerMachines{
					Total:  details.ControllerMachineCount,
					Active: details.ActiveControllerMachineCount,
				}
			} else {
				item.ControllerMachines = &jujucontroller.ControllerMachines{
					Total:  details.ControllerMachineCount,
					Active: details.ActiveControllerMachineCount,
				}
			}
		}
		controllers[controllerName] = item
	}
	return controllers, errs
}
