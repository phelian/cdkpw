package main

import "strings"

type Profile struct {
	Match   string `yaml:"match"`
	Profile string `yaml:"profile"`
}

type ProfileConfig struct {
	Profiles []Profile `yaml:"profiles"`
}

func findProfile(stackArg string, config *ProfileConfig) (string, bool) {
	for _, entry := range config.Profiles {
		if strings.Contains(stackArg, entry.Match) {
			return entry.Profile, true
		}
	}
	return "", false
}
