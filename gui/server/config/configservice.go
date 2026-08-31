package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/pelletier/go-toml/v2"
	"github.com/wailsapp/wails/v3/pkg/application"
)

const configFileName = "config.toml"

type ConfigService struct{}

// OpenConfigFile finds the default file or asks the user to select one, then
// returns its path and raw contents. Parsing and editing are deliberately not
// part of this first loading step.
func (s *ConfigService) OpenConfigFile() (*ConfigFile, error) {
	path, err := defaultConfigPath()
	if err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			return nil, err
		}
		return s.SelectConfigFile()
	}
	return s.readConfigFile(path)
}

// SelectConfigFile always opens the native picker, even when a default file
// exists. It is used when the user explicitly wants another configuration.
func (s *ConfigService) SelectConfigFile() (*ConfigFile, error) {
	path, err := application.Get().Dialog.OpenFile().
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
	return s.readConfigFile(path)
}

func (s *ConfigService) readConfigFile(path string) (*ConfigFile, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read configuration file %q: %w", path, err)
	}
	config, err := parseConfiguration(content)
	if err != nil {
		return nil, fmt.Errorf("parse configuration file %q: %w", path, err)
	}
	return &ConfigFile{
		Path:             path,
		Content:          string(content),
		Config:           config,
		ValidationErrors: validateConfiguration(config),
	}, nil
}

func parseConfiguration(content []byte) (*Configuration, error) {
	var config Configuration
	if err := toml.Unmarshal(content, &config); err != nil {
		return nil, fmt.Errorf("invalid TOML or configuration type: %w", err)
	}
	if config.Games.Enabled == nil {
		config.Games.Enabled = []string{}
	}
	return &config, nil
}

func validateConfiguration(config *Configuration) []ValidationError {
	var validationErrors []ValidationError

	supportedAuthentication := map[string]bool{
		"required": true,
		"cached":   true,
		"adaptive": true,
		"disabled": true,
	}
	if !supportedAuthentication[config.Authentication] {
		validationErrors = append(validationErrors, ValidationError{
			Field:   "Authentication",
			Message: "must be one of required, cached, adaptive or disabled",
		})
	}

	supportedGames := map[string]bool{
		"age1":   true,
		"age2":   true,
		"age3":   true,
		"age4":   true,
		"athens": true,
	}
	if len(config.Games.Enabled) == 0 {
		validationErrors = append(validationErrors, ValidationError{
			Field:   "Games.Enabled",
			Message: "must contain at least one game",
		})
	}
	for index, game := range config.Games.Enabled {
		if !supportedGames[game] {
			validationErrors = append(validationErrors, ValidationError{
				Field:   fmt.Sprintf("Games.Enabled.%d", index),
				Message: fmt.Sprintf("unsupported game %q", game),
			})
		}
	}

	if config.Announcement.Port < 1 || config.Announcement.Port > 65535 {
		validationErrors = append(validationErrors, ValidationError{
			Field:   "Announcement.Port",
			Message: "must be an integer between 1 and 65535",
		})
	}

	return validationErrors
}

func defaultConfigPath() (string, error) {
	executable, err := os.Executable()
	if err != nil {
		return "", fmt.Errorf("find executable: %w", err)
	}

	candidates := []string{
		filepath.Join(filepath.Dir(executable), "resources", "config", configFileName),
		filepath.Join(filepath.Dir(executable), configFileName),
		filepath.Join("resources", "config", configFileName),
		configFileName,
	}
	for _, candidate := range candidates {
		if _, err := os.Stat(candidate); err == nil {
			return candidate, nil
		} else if !os.IsNotExist(err) {
			return "", fmt.Errorf("check configuration file %q: %w", candidate, err)
		}
	}
	return "", os.ErrNotExist
}
