package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

const defaultConfigFile = ".cdkpw.yml"

var getUserHomeDir = os.UserHomeDir

type Profile struct {
	Match   string `yaml:"match"`
	Profile string `yaml:"profile"`
}

type Verbose int

const (
	SILENT Verbose = iota // 0
	INFO                  // 1
	DEBUG                 // 2
)

type Config struct {
	Profiles    []Profile `yaml:"profiles"`
	CdkLocation string    `yaml:"cdkLocation"`
	Verbose     Verbose   `yaml:"verbose"`
}

func (c *Config) findProfile(stackArg string) (string, bool) {
	// Find all matching profiles
	var matches []Profile
	for _, entry := range c.Profiles {
		if strings.Contains(stackArg, entry.Match) {
			matches = append(matches, entry)
		}
	}
	
	// No matches found
	if len(matches) == 0 {
		return "", false
	}
	
	// Sort by match string length (longest first = most specific)
	// If multiple profiles match, prefer the most specific one
	bestMatch := matches[0]
	for _, match := range matches[1:] {
		if len(match.Match) > len(bestMatch.Match) {
			bestMatch = match
		}
	}
	
	if c.Verbose >= INFO {
		fmt.Printf("cdkpw: Using profile %s for stack %s\n", bestMatch.Profile, stackArg)
	}
	return bestMatch.Profile, true
}

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

func loadConfig() (*Config, error) {
	configPath, err := getConfigFile()
	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("could not read config file at %s: %w", configPath, err)
	}

	config := Config{}
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("invalid YAML in %s: %w", configPath, err)
	}

	if config.CdkLocation == "" {
		config.CdkLocation = "cdk"
	}
	return &config, nil
}
