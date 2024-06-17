package surgo

import (
	"github.com/stretchr/testify/assert"
	"reflect"
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

type testStructWithID struct {
	ID ID
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
	t.Run("scan incompatible to pointer", func(t *testing.T) {
		m := map[string]any{
			"test": "other_test:123",
		}

		var s *nestedTestStructPtr
		err := scan(m, &s)
		assert.NoError(t, err)
		assert.Nil(t, s.Test)
	})
	t.Run("scan incompatible", func(t *testing.T) {
		m := map[string]any{
			"title": 123,
		}

		var s *nestedTestStructPtr
		err := scan(m, &s)
		assert.NoError(t, err)
		assert.Equal(t, "", s.Title)
	})
	t.Run("scan string to time.Time", func(t *testing.T) {
		m := "2011-01-01T00:00:00Z"

		var s time.Time
		err := scan(m, &s)
		assert.NoError(t, err)
		assert.Equal(t, time.Date(2011, 1, 1, 0, 0, 0, 0, time.UTC), s)
	})
	t.Run("scan string to time.Duration", func(t *testing.T) {
		m := "1h"

		var s time.Duration
		err := scan(m, &s)
		assert.NoError(t, err)
		assert.Equal(t, time.Hour, s)
	})
	t.Run("scan int id string to ID", func(t *testing.T) {
		m := map[string]any{
			"id": "test:123",
		}

		var s testStructWithID
		err := scan(m, &s)
		s.ID.Convert(reflect.Int)
		assert.NoError(t, err)
		assert.Equal(t, testStructWithID{
			ID: ID{123},
		}, s)
	})
	t.Run("scan string id string to ID", func(t *testing.T) {
		m := map[string]any{
			"id": "test:`primary`",
		}

		var s testStructWithID
		err := scan(m, &s)
		assert.NoError(t, err)
		assert.Equal(t, testStructWithID{
			ID: ID{"primary"},
		}, s)
	})
	t.Run("scan to map with string value", func(t *testing.T) {
		m := map[string]any{
			"title":  "test",
			"title2": "test2",
		}

		s := make(map[string]string)
		err := scan(m, &s)
		assert.NoError(t, err)
		assert.Equal(t, map[string]string{
			"title":  "test",
			"title2": "test2",
		}, s)
	})
	t.Run("scan to map with mixed values", func(t *testing.T) {
		m := map[string]any{
			"title":       "test",
			"description": 123,
			"exists":      true,
		}

		s := make(map[string]any)
		err := scan(m, &s)
		assert.NoError(t, err)
		assert.Equal(t, map[string]any{
			"title":       "test",
			"description": 123,
			"exists":      true,
		}, s)
	})
	t.Run("scan to map with struct values", func(t *testing.T) {
		m := map[string]any{
			"1": map[string]any{
				"success":    true,
				"text":       "hello",
				"age":        1.23,
				"created_at": 123,
			},
			"2": map[string]any{
				"success":    false,
				"text":       "bye",
				"age":        0,
				"created_at": 0,
			},
		}

		s := make(map[string]testStruct)
		err := scan(m, &s)
		assert.NoError(t, err)
		assert.Equal(t, map[string]testStruct{
			"1": {
				Success:   true,
				Text:      "hello",
				Age:       1.23,
				CreatedAt: 123,
			},
			"2": {
				Success:   false,
				Text:      "bye",
				Age:       0,
				CreatedAt: 0,
			},
		}, s)
	})
}
