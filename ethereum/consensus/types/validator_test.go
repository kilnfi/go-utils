//go:build !integration
// +build !integration

package types

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestValidatorUnmarshalMarshalCSV(t *testing.T) {
	record := []string{
		"10",
		"active_ongoing",
		"34919114473",
		"0x8efba2238a00d678306c6258105b058e3c8b0c1f36e821de42da7319c4221b77aa74135dab1860235e19d6515575c381",
		"0x00ea42f2e2c8e339759f42e72b2f6801485abfdfbd416f0ffdd1d1b07b33a9c0",
		"32000000000",
		"false",
		"0",
		"0",
		"18446744073709551615",
		"18446744073709551615",
	}
	v := new(Validator)
	t.Run("UnmarshalCSV", func(t *testing.T) {
		require.NoError(t, v.UnmarshalCSV(record))
	})

	t.Run("MarshalCSV", func(t *testing.T) {
		record2, err := v.MarshalCSV()
		require.NoError(t, err)
		assert.Equal(t, record, record2)
	})

}
