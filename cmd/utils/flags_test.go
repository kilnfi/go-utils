//go:build !integration
// +build !integration

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

		expectedFlagDesc string
	}{
		{
			baseDesc:         "This is a great flag",
			expectedFlagDesc: "This is a great flag",
		},
		{
			baseDesc:         "This is a great flag",
			envVar:           "GREAT",
			expectedFlagDesc: `This is a great flag [env: GREAT]`,
		},
		{
			baseDesc:         "This is a great flag",
			expectedFlagDesc: `This is a great flag`,
		},
		{
			baseDesc:         "This is a great flag",
			envVar:           "GREAT",
			expectedFlagDesc: `This is a great flag [env: GREAT]`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.testDesc, func(t *testing.T) {
			assert.Equal(t, tt.expectedFlagDesc, FlagDesc(tt.baseDesc, tt.envVar))
		})
	}
}
