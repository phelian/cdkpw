package main

import (
	"os/exec"
	"testing"

	"github.com/stretchr/testify/suite"
)

type argsSuite struct {
	suite.Suite
}

type commandSuite struct {
	suite.Suite
}

func (s *argsSuite) TestParseArgs() {
	tests := []struct {
		name     string
		input    []string
		expected CDKCommand
	}{
		{
			name:  "basic diff with profile and context",
			input: []string{"diff", "my-stack", "--profile", "my-profile", "--context", "key=value", "--exclusively"},
			expected: CDKCommand{
				Action:    "diff",
				StackName: "my-stack",
				Profile:   "my-profile",
				Context:   []string{"--context", "key=value"},
				Flags:     []string{"--exclusively"},
			},
		},
		{
			name:  "deploy without profile",
			input: []string{"deploy", "other-stack"},
			expected: CDKCommand{
				Action:    "deploy",
				StackName: "other-stack",
			},
		},
		{
			name:  "diff with -c context",
			input: []string{"diff", "-c", "debug=true", "stack-name"},
			expected: CDKCommand{
				Action:    "diff",
				StackName: "stack-name",
				Context:   []string{"-c", "debug=true"},
			},
		},
		{
			name:  "destroy with flags before stack",
			input: []string{"destroy", "--exclusively", "secure-stack"},
			expected: CDKCommand{
				Action:    "destroy",
				StackName: "secure-stack",
				Flags:     []string{"--exclusively"},
			},
		},
		{
			name:  "missing action",
			input: []string{},
			expected: CDKCommand{
				Action:    "",
				StackName: "",
			},
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			actual := parseArgs(tt.input)
			s.Equal(tt.expected.Action, actual.Action, "Action")
			s.Equal(tt.expected.StackName, actual.StackName, "StackName")
			s.Equal(tt.expected.Profile, actual.Profile, "Profile")
			s.Equal(tt.expected.Context, actual.Context, "Context")
			s.Equal(tt.expected.Flags, actual.Flags, "Flags")
		})
	}
}

func (s *commandSuite) TestSetProfile() {
	cmd := &CDKCommand{
		RawArgs: []string{"deploy", "MyStack"},
	}

	cmd.SetProfile("my-profile")

	s.Equal("my-profile", cmd.Profile)
	s.Equal([]string{"deploy", "MyStack", "--profile", "my-profile"}, cmd.RawArgs)

	// Should not override an existing profile
	cmd.SetProfile("another-profile")
	s.Equal("my-profile", cmd.Profile)
	s.Equal([]string{"deploy", "MyStack", "--profile", "my-profile"}, cmd.RawArgs)
}

func (s *commandSuite) TestIsProfiled() {
	s.True((&CDKCommand{Profile: "prod"}).IsProfiled())
	s.False((&CDKCommand{}).IsProfiled())
}

var mockExecutedArgs []string

func mockExecCommand(command string, args ...string) *exec.Cmd {
	mockExecutedArgs = append([]string{command}, args...)
	return exec.Command("echo", "mocked cdk call")
}

func (s *commandSuite) TestExecute() {
	original := execCommand
	defer func() { execCommand = original }()

	execCommand = mockExecCommand // inject mock

	cmd := &CDKCommand{
		RawArgs: []string{"deploy", "MyStack"},
	}

	cmd.Execute()

	s.Equal([]string{"cdk", "deploy", "MyStack"}, mockExecutedArgs)
}

func TestArgsAndCommand(t *testing.T) {
	suite.Run(t, new(argsSuite))
	suite.Run(t, new(commandSuite))
}
