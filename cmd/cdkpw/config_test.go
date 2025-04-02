package main

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/suite"
)

type configSuite struct {
	suite.Suite
	originalEnv string
	tempDir     string
}

func (s *configSuite) SetupTest() {
	s.originalEnv = os.Getenv("CDKPW_CONFIG")
	s.tempDir = s.T().TempDir()
}

func (s *configSuite) TearDownTest() {
	os.Setenv("CDKPW_CONFIG", s.originalEnv)
}

func (s *configSuite) TestGetConfigFile() {
	t := s.T()

	t.Run("uses CDKPW_CONFIG when set", func(t *testing.T) {
		expected := "/foo/config.yml"
		os.Setenv("CDKPW_CONFIG", expected)
		path, err := getConfigFile()
		s.Require().NoError(err)
		s.Equal(expected, path)
	})

	t.Run("falls back to $HOME/.cdk/.cdkpw.yml", func(t *testing.T) {
		os.Unsetenv("CDKPW_CONFIG")

		home, _ := os.UserHomeDir()
		expected := filepath.Join(home, ".cdk", defaultConfigFile)

		path, err := getConfigFile()
		s.Require().NoError(err)
		s.Equal(expected, path)
	})
}

func (s *configSuite) TestLoadValidConfig() {
	yamlContent := `
profiles:
  - match: Prod
    profile: prod_admin
  - match: Dev
    profile: dev_admin
`
	configPath := filepath.Join(s.tempDir, "config.yml")
	s.Require().NoError(os.WriteFile(configPath, []byte(yamlContent), 0600))

	os.Setenv("CDKPW_CONFIG", configPath)
	defer os.Unsetenv("CDKPW_CONFIG")

	config, err := loadConfig()

	s.Require().NoError(err)
	s.Require().Len(config.Profiles, 2)

	s.Equal("Prod", config.Profiles[0].Match)
	s.Equal("prod_admin", config.Profiles[0].Profile)
	s.Equal("Dev", config.Profiles[1].Match)
	s.Equal("dev_admin", config.Profiles[1].Profile)
}

func (s *configSuite) TestLoadInvalidConfig() {
	yamlContent := `
profiles:
  - match: Prod
profile: prod_admin
  			- match: Dev
    profile: dev_admin
`
	configPath := filepath.Join(s.tempDir, "config.yml")
	s.Require().NoError(os.WriteFile(configPath, []byte(yamlContent), 0600))

	os.Setenv("CDKPW_CONFIG", configPath)
	defer os.Unsetenv("CDKPW_CONFIG")

	_, err := loadConfig()
	s.Require().Error(err)
}

func (s *configSuite) TestGetConfigFile_HomeError() {
	original := getUserHomeDir
	defer func() { getUserHomeDir = original }()

	getUserHomeDir = func() (string, error) {
		return "", fmt.Errorf("mocked home dir error")
	}

	os.Unsetenv("CDKPW_CONFIG")

	_, err := getConfigFile()
	s.Error(err)
	s.Contains(err.Error(), "mocked home dir error")

	_, err = loadConfig()
	s.Error(err)
	s.Contains(err.Error(), "mocked home dir error")
}

func (s *configSuite) TestGetConfigFile_ReadError() {
	os.Setenv("CDKPW_CONFIG", "/does/not/exist")
	defer os.Unsetenv("CDKPW_CONFIG")

	_, err := loadConfig()
	s.Error(err)
	s.Contains(err.Error(), "could not read config file at")
}

func TestConfigSuite(t *testing.T) {
	suite.Run(t, new(configSuite))
}
