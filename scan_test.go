package surgo

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

type testStruct struct {
	Success   bool
	Text      string
	CreatedAt int `db:"created_at"`
	Age       float64
}

type testStructWithTime struct {
	Time time.Time
}

type nestedTestStruct struct {
	Title string
	Test  testStruct
}

type nestedTestStructPtr struct {
	Title string
	Test  *testStruct
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
	t.Run("scan single map to nested struct with pointer", func(t *testing.T) {
		m := map[string]any{
			"title": "hello",
			"test": map[string]any{
				"text": "world",
			},
		}

		var s nestedTestStructPtr
		err := scan(m, &s)
		assert.NoError(t, err)
		assert.Equal(t, nestedTestStructPtr{
			Title: "hello",
			Test: &testStruct{
				Text: "world",
			},
		}, s)
	})
	t.Run("scan multiple maps to nested with pointer", func(t *testing.T) {
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

		var s []nestedTestStructPtr
		err := scan(m, &s)
		assert.NoError(t, err)
		assert.Equal(t, []nestedTestStructPtr{
			{
				Title: "hello",
			},
			{
				Test: &testStruct{
					Text: "world",
				},
			},
		}, s)
	})
	t.Run("scan single map to nested struct with tag and pointer", func(t *testing.T) {
		m := map[string]any{
			"title": "hello",
			"test": map[string]any{
				"created_at": 123,
			},
		}

		var s nestedTestStructPtr
		err := scan(m, &s)
		assert.NoError(t, err)
		assert.Equal(t, nestedTestStructPtr{
			Title: "hello",
			Test: &testStruct{
				CreatedAt: 123,
			},
		}, s)
	})
	t.Run("scan nested nil struct to pointer", func(t *testing.T) {
		m := map[string]any{
			"title": "hello",
			"test":  nil,
		}

		var s nestedTestStructPtr
		err := scan(m, &s)
		assert.NoError(t, err)
		assert.Equal(t, nestedTestStructPtr{
			Title: "hello",
			Test:  nil,
		}, s)
	})
	t.Run("scan to time.Time field", func(t *testing.T) {
		m := map[string]any{
			"time": "2011-01-01T00:00:00Z",
		}

		var s testStructWithTime
		err := scan(m, &s)
		assert.NoError(t, err)
		assert.Equal(t, testStructWithTime{
			Time: time.Date(2011, 1, 1, 0, 0, 0, 0, time.UTC),
		}, s)
	})
}
