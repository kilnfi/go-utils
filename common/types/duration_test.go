//go:build !integration
// +build !integration

package common

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUnmarshal(t *testing.T) {
	tests := []struct {
		desc string

		// JSON body of the request
		bytes            []byte
		expectedDuration Duration
	}{
		{
			desc:             "int",
			bytes:            []byte(`20`),
			expectedDuration: Duration{time.Duration(20)},
		},
		{
			desc:             "string",
			bytes:            []byte(`"15s30ns"`),
			expectedDuration: Duration{time.Duration(15000000030)},
		},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			dur := new(Duration)
			err := json.Unmarshal(tt.bytes, dur)
			require.NoError(t, err, "Unmarshal should not error")
			assert.Equal(t, tt.expectedDuration, *dur, "Duration should be correct")
		})
	}
}

func TestMarshal(t *testing.T) {
	tests := []struct {
		desc string

		// JSON body of the request
		v            interface{}
		expectedByte []byte
		expectedErr  error
	}{
		{
			desc:         "duration",
			v:            Duration{time.Duration(20)},
			expectedByte: []byte(`"20ns"`),
		},
		{
			desc: "structPtr",
			v: struct {
				Duration *Duration `json:"key"`
			}{
				Duration: &Duration{time.Duration(20)},
			},
			expectedByte: []byte(`{"key":"20ns"}`),
		},
		{
			desc: "structValue",
			v: struct {
				Duration Duration `json:"key"`
			}{
				Duration: Duration{time.Duration(20)},
			},
			expectedByte: []byte(`{"key":"20ns"}`),
		},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			b, err := json.Marshal(tt.v)
			if tt.expectedErr != nil {
				require.Error(t, err, "Marshal should error")
				assert.Equal(t, tt.expectedErr.Error(), err.Error(), "Error message should be correct")
			} else {
				require.NoError(t, err, "Marshal should not error")
				require.Equal(t, tt.expectedByte, b)
			}
		})
	}
}
