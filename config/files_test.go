package config_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/adrg/xdg"
	"github.com/jneo8/jujuspell/config" // Adjust this import path based on your actual module path

	"github.com/stretchr/testify/assert"
)

func TestCreateLogFile(t *testing.T) {

	// Reset environment after test execution.
	environ := os.Environ()
	defer func() {
		os.Clearenv()
		for _, env := range environ {
			parts := strings.SplitN(env, "=", 2)
			if len(parts) != 2 {
				continue
			}
			os.Setenv(parts[0], parts[1])
		}

		xdg.Reload()
	}()

	// Temporarily set XDG_STATE_HOME
	tempDir := t.TempDir()
	os.Setenv("XDG_STATE_HOME", tempDir)
	xdg.Reload()

	// Call the function under test
	logFilePath, err := config.CreateLogFile()

	// Assertions
	expectedPath := filepath.Join(tempDir, config.AppName, "jujuspell.log")
	assert.NoError(t, err)
	assert.Equal(t, expectedPath, logFilePath)
}

func TestEnsureFullPath(t *testing.T) {
	tempDir := t.TempDir()
	targetPath := filepath.Join(tempDir, "newdir")

	err := config.EnsureFullPath(targetPath, config.DefaultDirMod)
	assert.NoError(t, err)

	info, err := os.Stat(targetPath)
	assert.NoError(t, err)
	assert.True(t, info.IsDir())
}
