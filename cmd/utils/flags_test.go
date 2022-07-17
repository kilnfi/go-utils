package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFlagDesc(t *testing.T) {
	tests := []struct {
		testDesc string

		baseDesc string
		envVar   string
		dfault   interface{}

		expectedFlagDesc string
	}{
		{
			baseDesc:         "This is a great flag",
			expectedFlagDesc: "This is a great flag",
		},
		{
			baseDesc: "This is a great flag",
			envVar:   "GREAT",
			expectedFlagDesc: `This is a great flag
  Environment variable: GREAT`,
		},
		{
			baseDesc: "This is a great flag",
			dfault:   []string{"a", "b"},
			expectedFlagDesc: `This is a great flag
  Default: ["a","b"]`,
		},
		{
			baseDesc: "This is a great flag",
			envVar:   "GREAT",
			dfault:   []string{"a", "b"},
			expectedFlagDesc: `This is a great flag
  Environment variable: GREAT
  Default: ["a","b"]`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.testDesc, func(t *testing.T) {
			assert.Equal(t, tt.expectedFlagDesc, flagDesc(tt.baseDesc, tt.envVar, tt.dfault))
		})
	}
}
