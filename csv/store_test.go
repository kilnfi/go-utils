package csv

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type testStruct struct {
	rec []string
}

func (s testStruct) MarshalCSV() ([]string, error) {
	return s.rec, nil
}

func (s *testStruct) UnmarshalCSV(records []string) error {
	s.rec = records
	return nil
}

func TestStore(t *testing.T) {
	dir := t.TempDir()
	store := NewStore(filepath.Join(dir, "test.csv"))

	t.Run("WriteAllStructs", func(t *testing.T) {
		require.NoError(t, store.WriteAllStructs([]interface{}{
			&testStruct{rec: []string{"a", "b", "c"}},
			&testStruct{rec: []string{"d", "e", "f"}},
			&testStruct{rec: []string{"g", "h", "i"}},
		}))
	})

	var recStructs []*testStruct
	t.Run("ReadAllStructs", func(t *testing.T) {
		require.NoError(t, store.ReadAllStructs(&recStructs))
		assert.Equal(
			t,
			[]*testStruct{
				{rec: []string{"a", "b", "c"}},
				{rec: []string{"d", "e", "f"}},
				{rec: []string{"g", "h", "i"}},
			},
			recStructs,
		)
	})
}
