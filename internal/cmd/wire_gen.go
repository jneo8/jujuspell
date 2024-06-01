// Code generated by Wire. DO NOT EDIT.

//go:generate go run -mod=mod github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package cmd

import (
	"github.com/jneo8/juju-spell/internal/app"
)

// Injectors from wire.go:

func InitializeExecute() func() {
	logger := app.NewLogger()
	v := GetExecute(logger)
	return v
}
