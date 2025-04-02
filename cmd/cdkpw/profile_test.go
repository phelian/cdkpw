package main

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type profileSuite struct {
	suite.Suite
	tempDir string
}

func (s *profileSuite) TestFindProfile_Matching() {
	config := &ProfileConfig{
		Profiles: []Profile{
			{Match: "Prod", Profile: "prod_admin"},
			{Match: "Dev", Profile: "dev_admin"},
			{Match: "Secure", Profile: "secure_admin"},
			{Match: "Api", Profile: "api_admin"},
		},
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
			actual, ok := findProfile(tt.stackArg, config)
			s.Equal(tt.want, actual)
			s.Equal(tt.found, ok)
		})
	}
}

func TestProfileSuite(t *testing.T) {
	suite.Run(t, new(profileSuite))
}
