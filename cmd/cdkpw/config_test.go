package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/suite"
)

type configSuite struct {
	suite.Suite
	originalEnv string
	tempDir     string
}

func (s *configSuite) TestConfigFindProfile_Matching() {
	config := &Config{
		Profiles: []Profile{
			{Match: "Prod", Profile: "prod_admin"},
			{Match: "Dev", Profile: "dev_admin"},
			{Match: "Secure", Profile: "secure_admin"},
			{Match: "Api", Profile: "api_admin"},
		},
		CdkLocation: "/usr/local/bin/cdk",
	}

	tests := []struct {
		name     string
		stackArg string
		want     string
		found    bool
	}{
		{
			name:     "matches Prod",
			stackArg: "ProdAppStack",
			want:     "prod_admin",
			found:    true,
		},
		{
			name:     "matches Dev",
			stackArg: "DevWorkerStack",
			want:     "dev_admin",
			found:    true,
		},
		{
			name:     "matches Secure",
			stackArg: "SecureZoneEKS",
			want:     "secure_admin",
			found:    true,
		},
		{
			name:     "matches Api",
			stackArg: "CustomerApiStack",
			want:     "api_admin",
			found:    true,
		},
		{
			name:     "no match",
			stackArg: "StagingStack",
			want:     "",
			found:    false,
		},
		{
			name:     "empty stack arg",
			stackArg: "",
			want:     "",
			found:    false,
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			actual, ok := config.findProfile(tt.stackArg)
			s.Equal(tt.want, actual)
			s.Equal(tt.found, ok)
		})
	}
	s.T().Run("cdk location", func(t *testing.T) {
		s.Equal(config.CdkLocation, "/usr/local/bin/cdk")
	})
}

func (s *configSuite) TestFindProfile_VerbosePrints() {
	t := s.T()

	cfg := Config{
		Profiles: []Profile{
			{Match: "Backup", Profile: "backup_admin"},
		},
		Verbose: INFO,
	}

	// Capture stdout
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	// Run the function
	stackArg := "BackupStack"
	profile, ok := cfg.findProfile(stackArg)

	// Stop capturing
	w.Close()
	var buf bytes.Buffer
	io.Copy(&buf, r)
	os.Stdout = old

	// Check results
	if !ok {
		t.Fatalf("Expected to find a matching profile")
	}
	if profile != "backup_admin" {
		t.Errorf("Expected profile 'backup_admin', got '%s'", profile)
	}

	output := buf.String()
	expected := fmt.Sprintf("cdkpw: Using profile %s for stack %s\n", profile, stackArg)
	if !strings.Contains(output, expected) {
		t.Errorf("Expected output to contain %q, but got %q", expected, output)
	}
}

func (s *configSuite) TestFindProfile_Silent() {
	t := s.T()

	cfg := Config{
		Profiles: []Profile{
			{Match: "Backup", Profile: "backup_admin"},
		},
		Verbose: SILENT,
	}

	// Capture stdout
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	// Run the function
	stackArg := "BackupStack"
	profile, ok := cfg.findProfile(stackArg)

	// Stop capturing
	w.Close()
	var buf bytes.Buffer
	io.Copy(&buf, r)
	os.Stdout = old

	// Check results
	if !ok {
		t.Fatalf("Expected to find a matching profile")
	}
	if profile != "backup_admin" {
		t.Errorf("Expected profile 'backup_admin', got '%s'", profile)
	}

	output := buf.String()
	if output != "" {
		t.Errorf("Expected no output, but got: %q", output)
	}
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
