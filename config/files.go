package config

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/adrg/xdg"
)

const (
	LogsFile                  = "jujuspell.log"
	DefaultDirMod os.FileMode = 0744
	LogsFileMod   os.FileMode = 0644
)

func CreateLogFile() (string, error) {
	appLogDir, err := xdg.StateFile(AppName)
	if err != nil {
		return "", err
	}
	if err := EnsureFullPath(appLogDir, DefaultDirMod); err != nil {
		return "", err
	}
	AppLogFile := filepath.Join(appLogDir, LogsFile)
	return AppLogFile, nil
}

func EnsureFullPath(path string, mod os.FileMode) error {
	if _, err := os.Stat(path); errors.Is(err, fs.ErrNotExist) {
		if err = os.MkdirAll(path, mod); err != nil {
			return err
		}
	}
	return nil
}
