package surgo

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestInternals(t *testing.T) {
	t.Run("Parse struct to map", func(t *testing.T) {
		t.Parallel()
		obj := testObject{TimeSince: 1, Name: "John"}
		m := structToMap(obj)
		assert.Equal(t, map[string]any{"id": 1, "Name": "John"}, m)
	})
	t.Run("Parse slice to map", func(t *testing.T) {
		t.Parallel()
		slice := []any{1, "John"}
		m := sliceToMap(slice)
		assert.Equal(t, map[string]any{"1": 1, "2": "John"}, m)
	})
	t.Run("Parse array to map", func(t *testing.T) {
		t.Parallel()
		array := [2]any{1, "John"}
		m := sliceToMap(array)
		assert.Equal(t, map[string]any{"1": 1, "2": "John"}, m)
	})
}
