package main

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type argsSuite struct {
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

func TestArgs(t *testing.T) {
	suite.Run(t, new(argsSuite))
}
