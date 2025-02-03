package cmd

import (
	"io"
	"os"
	"time"

	"github.com/jneo8/jujuspell/config"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

func init() {
}

func configureLogger(logFile io.Writer) {
	consoleWriter := zerolog.ConsoleWriter{
		Out:        logFile,
		TimeFormat: time.RFC3339,
	}

	log.Logger = zerolog.New(consoleWriter).With().Timestamp().Logger()
}

const (
	appName      = config.AppName
	shortAppDesc = "A graphical CLI for your Juju cluster management."
	longAppDesc  = "JujuSpell is a CLI to view and manage your juju clusters."
)

var (
	rootCmd = &cobra.Command{
		Use:   appName,
		Short: shortAppDesc,
		Long:  longAppDesc,
		RunE:  run,
	}
)

func run(cmd *cobra.Command, args []string) error {
	// Setup logger & log file
	logFilePath, err := config.CreateLogFile()
	if err != nil {
		return err
	}
	logFile, err := os.OpenFile(logFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, config.LogsFileMod)
	if err != nil {
		return err
	}
	defer logFile.Close()
	configureLogger(logFile)
	log.Debug().Msg("Setup logger")
	return nil
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
