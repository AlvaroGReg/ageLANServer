package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/wailsapp/wails/v3/pkg/application"
)

const configFileName = "config.game.toml"

// ConfigFile is the small bridge between Go and the React editor.
type ConfigFile struct {
	Path    string `json:"path"`
	Content string `json:"content"`
}

type ConfigService struct{}

// OpenConfigFile finds the default file or asks the user to select one, then
// returns its path and contents. Parsing and editing remain in the frontend by now.
func (s *ConfigService) OpenConfigFile() (*ConfigFile, error) {
	path := defaultConfigPath()
	if _, err := os.Stat(path); os.IsNotExist(err) {
		path, err = application.Get().Dialog.OpenFile().
			SetTitle("Select configuration file").
			CanChooseFiles(true).
			CanChooseDirectories(false).
			AddFilter("TOML files", "*.toml").
			PromptForSingleSelection()
		if err != nil {
			return nil, fmt.Errorf("open configuration picker: %w", err)
		}
		if path == "" {
			return nil, errors.New("no configuration file selected")
		}
	} else if err != nil {
		return nil, fmt.Errorf("check default configuration file: %w", err)
	}

	content, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read configuration file %q: %w", path, err)
	}
	return &ConfigFile{Path: path, Content: string(content)}, nil
}

func defaultConfigPath() string {
	executable, err := os.Executable()
	if err != nil {
		return configFileName
	}
	return filepath.Join(filepath.Dir(executable), configFileName)
}
