package marshal

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestMarshaler_Unmarshal(t *testing.T) {
	m := Marshaler("json")

	t.Run("simple bool unmarshal", func(t *testing.T) {
		var b bool
		err := m.Unmarshal(true, &b)
		assert.NoError(t, err)
		assert.True(t, b)
	})
	t.Run("simple string unmarshal", func(t *testing.T) {
		var s string
		err := m.Unmarshal("test", &s)
		assert.NoError(t, err)
		assert.Equal(t, "test", s)
	})
	t.Run("simple time unmarshal", func(t *testing.T) {
		var tm time.Time
		err := m.Unmarshal("2020-01-01T00:00:00Z", &tm)
		assert.NoError(t, err)
		assert.Equal(t, time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC), tm)
	})
	t.Run("simple duration unmarshal", func(t *testing.T) {
		var d time.Duration
		err := m.Unmarshal("1s", &d)
		assert.NoError(t, err)
		assert.Equal(t, time.Second, d)
	})
	t.Run("all numbers unmarshal", func(t *testing.T) {
		numberTypes := []any{
			42, int8(42), int16(42), int32(42), int64(42),
			uint(42), uint8(42), uint16(42), uint32(42), uint64(42),
			float32(42.42), 42.42,
		}

		for _, src := range numberTypes {
			for _, dest := range numberTypes {
				srcNum, destNum := src, dest
				err := m.Unmarshal(srcNum, &destNum)
				assert.NoError(t, err)
				assert.Equal(t, srcNum, destNum)
			}
		}
	})
	t.Run("simple slice unmarshal", func(t *testing.T) {
		var s []string
		err := m.Unmarshal([]string{"test", "test2"}, &s)
		assert.NoError(t, err)
		assert.Equal(t, []string{"test", "test2"}, s)
	})
	t.Run("byte unmarshal", func(t *testing.T) {
		var b byte
		err := m.Unmarshal(byte('a'), &b)
		assert.NoError(t, err)
		assert.Equal(t, byte('a'), b)
	})
	t.Run("byte slice unmarshal", func(t *testing.T) {
		var b []byte
		err := m.Unmarshal([]byte("test"), &b)
		assert.NoError(t, err)
		assert.Equal(t, []byte("test"), b)
	})
	t.Run("assign to any", func(t *testing.T) {
		var a any
		err := m.Unmarshal("test", &a)
		assert.NoError(t, err)
		assert.Equal(t, "test", a)
	})
	t.Run("map to map", func(t *testing.T) {
		var m1 map[string]any
		err := m.Unmarshal(map[string]any{"test": "test", "num": 42}, &m1)
		assert.NoError(t, err)
		assert.Equal(t, map[string]any{"test": "test", "num": 42}, m1)
	})
	t.Run("map to struct", func(t *testing.T) {
		type testStruct struct {
			Test string
			Num  int
		}
		var s testStruct

		err := m.Unmarshal(map[string]any{"Test": "test", "Num": 42}, &s)
		assert.NoError(t, err)
		assert.Equal(t, testStruct{"test", 42}, s)
	})
	t.Run("map to struct with tags", func(t *testing.T) {
		type testStruct struct {
			Test string `db:"test"`
			Num  int    `db:"num"`
		}
		var s testStruct

		err := m.Unmarshal(map[string]any{"test": "test", "num": 42}, &s)
		assert.NoError(t, err)
		assert.Equal(t, testStruct{"test", 42}, s)
	})
	t.Run("map to struct with fallback tags", func(t *testing.T) {
		type testStruct struct {
			Test string `json:"testjson" db:"test"`
			Num  int    `json:"numjson"`
		}
		var s testStruct

		err := m.Unmarshal(map[string]any{"testjson": "incorrect", "test": "test", "Num": 11, "numjson": 42}, &s)
		assert.NoError(t, err)
		assert.Equal(t, testStruct{"test", 42}, s)
	})
	t.Run("map to slice of structs", func(t *testing.T) {
		type testStruct struct {
			Test string
			Num  int
		}
		var s []testStruct

		err := m.Unmarshal([]map[string]any{{"Test": "test", "Num": 42}, {"Test": "test2", "Num": 43}}, &s)
		assert.NoError(t, err)
		assert.Equal(t, []testStruct{{"test", 42}, {"test2", 43}}, s)
	})
	t.Run("nested map to struct", func(t *testing.T) {
		type (
			nestedStruct struct {
				Value string
			}
			testStruct struct {
				Nested nestedStruct
			}
		)

		var s testStruct
		err := m.Unmarshal(map[string]any{"Nested": map[string]any{"Value": "test"}}, &s)
		assert.NoError(t, err)
		assert.Equal(t, testStruct{nestedStruct{"test"}}, s)
	})
	t.Run("nested map to struct pointer", func(t *testing.T) {
		type (
			nestedStruct struct {
				Value string
			}
			testStruct struct {
				Nested *nestedStruct
			}
		)

		var s testStruct
		err := m.Unmarshal(map[string]any{"Nested": map[string]any{"Value": "test"}}, &s)
		assert.NoError(t, err)
		assert.Equal(t, testStruct{&nestedStruct{"test"}}, s)
	})
	t.Run("slice of any to slice of defined", func(t *testing.T) {
		var s []string
		err := m.Unmarshal([]any{"a", "b"}, &s)
		assert.NoError(t, err)
		assert.Equal(t, []string{"a", "b"}, s)
	})
	t.Run("map as any to struct ", func(t *testing.T) {
		type testStruct struct {
			Test string
			Num  int
		}
		var s testStruct
		var in any = map[string]any{"Test": "test", "Num": 42}
		err := m.Unmarshal(in, &s)
		assert.NoError(t, err)
		assert.Equal(t, testStruct{"test", 42}, s)
	})
	t.Run("type alias unmarshal", func(t *testing.T) {
		type testAlias string
		var a testAlias
		err := m.Unmarshal("test", &a)
		assert.NoError(t, err)
		assert.Equal(t, testAlias("test"), a)
	})
	t.Run("src is null", func(t *testing.T) {
		var s string
		err := m.Unmarshal(nil, &s)
		assert.NoError(t, err)
		assert.Equal(t, "", s)
	})
	t.Run("null in src object", func(t *testing.T) {
		type testStruct struct {
			Test *struct {
				A string
			} `db:"test"`
		}
		var s testStruct
		err := m.Unmarshal(map[string]any{"test": nil}, &s)
		assert.NoError(t, err)
		assert.Equal(t, testStruct{
			Test: nil,
		}, s)
	})
	t.Run("anonymous struct", func(t *testing.T) {
		type Anonymous struct {
			Test string `db:"test"`
		}
		type testStruct struct {
			Anonymous
		}
		var s testStruct
		err := m.Unmarshal(map[string]any{"test": "test"}, &s)
		assert.NoError(t, err)
		assert.Equal(t, "test", s.Test)
	})
}
