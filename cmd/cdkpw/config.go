package main

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

const defaultConfigFile = ".cdkpw.yml"

var getUserHomeDir = os.UserHomeDir

// getConfigPath retrieves the path to the configuration file.
func getConfigFile() (string, error) {
	if customConfigPath := os.Getenv("CDKPW_CONFIG"); customConfigPath != "" {
		return customConfigPath, nil
	}

	home, err := getUserHomeDir()
	if err != nil {
		return "", fmt.Errorf("unable to determine config directory: %w", err)
	}
	return filepath.Join(home, ".cdk", defaultConfigFile), nil
}

func loadConfig() (*ProfileConfig, error) {
	configPath, err := getConfigFile()
	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("could not read config file at %s: %w", configPath, err)
	}

	var config ProfileConfig
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("invalid YAML in %s: %w", configPath, err)
	}
	return &config, nil
}
