package surgo

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

type testStruct struct {
	Success   bool
	Text      string
	CreatedAt int `db:"created_at"`
	Age       float64
}

type nestedTestStruct struct {
	Title string
	Test  testStruct
}

func Test_scan(t *testing.T) {
	t.Run("scan single map to struct", func(t *testing.T) {
		m := map[string]any{
			"success": true,
			"text":    "hello",
			"age":     1.23,
		}

		var s testStruct
		err := scan(m, &s)
		assert.NoError(t, err)
		assert.Equal(t, testStruct{
			Success: true,
			Text:    "hello",
			Age:     1.23,
		}, s)
	})
	t.Run("scan single map to struct with to map", func(t *testing.T) {
		m := map[string]any{
			"created_at": 123,
		}

		var s testStruct
		err := scan(m, &s)
		assert.NoError(t, err)
		assert.Equal(t, testStruct{
			CreatedAt: 123,
		}, s)
	})
	t.Run("scan multiple maps to struct", func(t *testing.T) {
		m := []map[string]any{
			{
				"text": "hello",
			},
			{
				"age": 1.23,
			},
		}

		var s []testStruct
		err := scan(m, &s)
		assert.NoError(t, err)
		assert.Equal(t, []testStruct{
			{
				Text: "hello",
			},
			{
				Age: 1.23,
			},
		}, s)
	})
	t.Run("scan multiple maps to struct with to map", func(t *testing.T) {
		m := []map[string]any{
			{
				"created_at": 123,
			},
		}

		var s []testStruct
		err := scan(m, &s)
		assert.NoError(t, err)
		assert.Equal(t, []testStruct{
			{
				CreatedAt: 123,
			},
		}, s)
	})
	t.Run("scan single value to variable", func(t *testing.T) {
		m := 123

		var v int
		err := scan(m, &v)
		assert.NoError(t, err)
		assert.Equal(t, 123, v)
	})
	t.Run("scan multiple values to variable", func(t *testing.T) {
		m := []any{
			"hello",
			1.23,
		}

		var v []any
		err := scan(m, &v)
		assert.NoError(t, err)
		assert.Equal(t, []any{"hello", 1.23}, v)
	})
	t.Run("scan single map to nested struct", func(t *testing.T) {
		m := map[string]any{
			"title": "hello",
			"test": map[string]any{
				"text": "world",
			},
		}

		var s nestedTestStruct
		err := scan(m, &s)
		assert.NoError(t, err)
		assert.Equal(t, nestedTestStruct{
			Title: "hello",
			Test: testStruct{
				Text: "world",
			},
		}, s)
	})
	t.Run("scan multiple maps to nested struct", func(t *testing.T) {
		m := []map[string]any{
			{
				"title": "hello",
			},
			{
				"test": map[string]any{
					"text": "world",
				},
			},
		}

		var s []nestedTestStruct
		err := scan(m, &s)
		assert.NoError(t, err)
		assert.Equal(t, []nestedTestStruct{
			{
				Title: "hello",
			},
			{
				Test: testStruct{
					Text: "world",
				},
			},
		}, s)
	})
	t.Run("scan single map to nested struct with tag", func(t *testing.T) {
		m := map[string]any{
			"title": "hello",
			"test": map[string]any{
				"created_at": 123,
			},
		}

		var s nestedTestStruct
		err := scan(m, &s)
		assert.NoError(t, err)
		assert.Equal(t, nestedTestStruct{
			Title: "hello",
			Test: testStruct{
				CreatedAt: 123,
			},
		}, s)
	})
}
