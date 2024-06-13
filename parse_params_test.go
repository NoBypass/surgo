package surgo

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
	"time"
)

type mockQueryAgent struct {
	mock.Mock
}

func (mqa *mockQueryAgent) Query(sql string, vars any) (any, error) {
	args := mqa.Called(sql, vars)
	return args.Get(0), args.Error(1)
}

func (mqa *mockQueryAgent) Close() {}

type arbitraryData struct {
	Num     int
	Text    string
	Boolean bool `db:"bool"`
}

var nilMap map[string]any
var emptyResponse = []any{
	map[string]any{
		"result": nil,
		"status": "OK",
		"time":   "0s",
	},
}

func Test_parseQuery(t *testing.T) {
	t.Run("Unmodified query", func(t *testing.T) {
		m := new(mockQueryAgent)
		db := &DB{m}
		query := "SELECT * FROM test;"

		m.On("Query", mock.Anything, mock.Anything).Return(emptyResponse, nil)
		_, err := db.Exec(query)

		assert.NoError(t, err)
		m.AssertCalled(t, "Query", query, nilMap)
	})
	t.Run("Query a slice parameter", func(t *testing.T) {
		m := new(mockQueryAgent)
		db := &DB{m}
		query := "SELECT * FROM test WHERE id = $1;"

		m.On("Query", mock.Anything, mock.Anything).Return(emptyResponse, nil)
		_, err := db.Exec(query, 1)

		assert.NoError(t, err)
		m.AssertCalled(t, "Query", query, map[string]any{"1": 1})
	})
	t.Run("Query wuth zero value parameter", func(t *testing.T) {
		m := new(mockQueryAgent)
		db := &DB{m}
		query := "SELECT * FROM test WHERE id = $1;"

		m.On("Query", mock.Anything, mock.Anything).Return(emptyResponse, nil)
		_, err := db.Exec(query, 0)

		assert.NoError(t, err)
		m.AssertCalled(t, "Query", query, map[string]any{"1": 0})
	})
	t.Run("Query with multiple slice parameters", func(t *testing.T) {
		m := new(mockQueryAgent)
		db := &DB{m}
		query := "SELECT * FROM test WHERE id = $1 AND name = $2;"

		m.On("Query", mock.Anything, mock.Anything).Return(emptyResponse, nil)
		_, err := db.Exec(query, 1, "test")

		assert.NoError(t, err)
		m.AssertCalled(t, "Query", query, map[string]any{"1": 1, "2": "test"})
	})
	t.Run("Query with map parameter", func(t *testing.T) {
		m := new(mockQueryAgent)
		db := &DB{m}
		query := "SELECT * FROM test WHERE id = $id;"

		m.On("Query", mock.Anything, mock.Anything).Return(emptyResponse, nil)
		_, err := db.Exec(query, map[string]any{"id": 1})

		assert.NoError(t, err)
		m.AssertCalled(t, "Query", query, map[string]any{"id": 1})
	})
	t.Run("Query with struct parameter", func(t *testing.T) {
		m := new(mockQueryAgent)
		db := &DB{m}
		query := "SELECT * FROM test WHERE id = $num;"

		m.On("Query", mock.Anything, mock.Anything).Return(emptyResponse, nil)
		_, err := db.Exec(query, arbitraryData{Num: 1})

		assert.NoError(t, err)
		m.AssertCalled(t, "Query", query, map[string]any{"num": 1, "bool": false, "text": ""})
	})
	t.Run("Query with multiple map parameters", func(t *testing.T) {
		m := new(mockQueryAgent)
		db := &DB{m}
		query := "SELECT * FROM test WHERE id = $num AND name = $text;"

		m.On("Query", mock.Anything, mock.Anything).Return(emptyResponse, nil)
		_, err := db.Exec(query, map[string]any{"num": 1}, map[string]any{"text": "test"})

		assert.NoError(t, err)
		m.AssertCalled(t, "Query", query, map[string]any{"num": 1, "text": "test"})
	})
	t.Run("Query with struct parameter and tag", func(t *testing.T) {
		m := new(mockQueryAgent)
		db := &DB{m}
		query := "SELECT * FROM test WHERE id = $bool;"

		m.On("Query", mock.Anything, mock.Anything).Return(emptyResponse, nil)
		_, err := db.Exec(query, arbitraryData{Boolean: true})

		assert.NoError(t, err)
		m.AssertCalled(t, "Query", query, map[string]any{"bool": true, "num": 0, "text": ""})
	})
	t.Run("Query with multiple struct parameters", func(t *testing.T) {
		m := new(mockQueryAgent)
		db := &DB{m}
		query := "SELECT * FROM test WHERE id = $num AND name = $text;"

		m.On("Query", mock.Anything, mock.Anything).Return(emptyResponse, nil)
		_, err := db.Exec(query, arbitraryData{Num: 1}, arbitraryData{Text: "test"})

		assert.NoError(t, err)
		m.AssertCalled(t, "Query", query, map[string]any{"num": 0, "text": "test", "bool": false})
	})
	t.Run("Query with a mix of map and struct parameters", func(t *testing.T) {
		m := new(mockQueryAgent)
		db := &DB{m}
		query := "SELECT * FROM test WHERE id = $num AND name = $text;"

		m.On("Query", mock.Anything, mock.Anything).Return(emptyResponse, nil)
		_, err := db.Exec(query, arbitraryData{Text: "test"}, map[string]any{"num": 1})

		assert.NoError(t, err)
		m.AssertCalled(t, "Query", query, map[string]any{"num": 1, "text": "test", "bool": false})
	})
	t.Run("Query with a mix of slice and struct parameters", func(t *testing.T) {
		m := new(mockQueryAgent)
		db := &DB{m}
		query := "SELECT * FROM test WHERE id = $1 AND name = $text;"

		m.On("Query", mock.Anything, mock.Anything).Return(emptyResponse, nil)
		_, err := db.Exec(query, 1, arbitraryData{Text: "test"})

		assert.NoError(t, err)
		m.AssertCalled(t, "Query", query, map[string]any{"1": 1, "text": "test", "bool": false, "num": 0})
	})
	t.Run("Query with a mix of slice and map parameters", func(t *testing.T) {
		m := new(mockQueryAgent)
		db := &DB{m}
		query := "SELECT * FROM test WHERE id = $1 AND name = $text;"

		m.On("Query", mock.Anything, mock.Anything).Return(emptyResponse, nil)
		_, err := db.Exec(query, 1, map[string]any{"text": "test"})

		assert.NoError(t, err)
		m.AssertCalled(t, "Query", query, map[string]any{"1": 1, "text": "test"})
	})
	t.Run("Query with an id parameter", func(t *testing.T) {
		m := new(mockQueryAgent)
		db := &DB{m}
		query := "SELECT * FROM test:$"

		m.On("Query", mock.Anything, mock.Anything).Return(emptyResponse, nil)
		_, err := db.Exec(query, ID{1})

		assert.NoError(t, err)
		m.AssertCalled(t, "Query", "SELECT * FROM test:1;", map[string]any{})
	})
	t.Run("Query with a string id parameter", func(t *testing.T) {
		m := new(mockQueryAgent)
		db := &DB{m}
		query := "SELECT * FROM test:$"

		m.On("Query", mock.Anything, mock.Anything).Return(emptyResponse, nil)
		_, err := db.Exec(query, ID{"1"})

		assert.NoError(t, err)
		m.AssertCalled(t, "Query", "SELECT * FROM test:`1`;", map[string]any{})
	})
	t.Run("Query with multiple id parameters", func(t *testing.T) {
		m := new(mockQueryAgent)
		db := &DB{m}
		query := "SELECT * FROM test:$, foo:$"

		m.On("Query", mock.Anything, mock.Anything).Return(emptyResponse, nil)
		_, err := db.Exec(query, ID{1}, ID{2})

		assert.NoError(t, err)
		m.AssertCalled(t, "Query", "SELECT * FROM test:1, foo:2;", map[string]any{})
	})
	t.Run("Query with an array id", func(t *testing.T) {
		m := new(mockQueryAgent)
		db := &DB{m}
		query := "SELECT * FROM test:$"

		m.On("Query", mock.Anything, mock.Anything).Return(emptyResponse, nil)
		_, err := db.Exec(query, ID{1, "2"})

		assert.NoError(t, err)
		m.AssertCalled(t, "Query", "SELECT * FROM test:[1, '2'];", map[string]any{})
	})
	t.Run("Query with a range id", func(t *testing.T) {
		m := new(mockQueryAgent)
		db := &DB{m}
		query := "SELECT * FROM test:$"

		m.On("Query", mock.Anything, mock.Anything).Return(emptyResponse, nil)
		_, err := db.Exec(query, Range{ID{1, "2"}, ID{3, "4"}})

		assert.NoError(t, err)
		m.AssertCalled(t, "Query", "SELECT * FROM test:[1, '2']..[3, '4'];", map[string]any{})
	})
	t.Run("Query with an id and slice parameter", func(t *testing.T) {
		m := new(mockQueryAgent)
		db := &DB{m}
		query := "SELECT * FROM test:$ WHERE id = $1"

		m.On("Query", mock.Anything, mock.Anything).Return(emptyResponse, nil)
		_, err := db.Exec(query, ID{"1"}, 1)

		assert.NoError(t, err)
		m.AssertCalled(t, "Query", "SELECT * FROM test:`1` WHERE id = $1;", map[string]any{"1": 1})
	})
	t.Run("Query with an id, map and slice parameter", func(t *testing.T) {
		m := new(mockQueryAgent)
		db := &DB{m}
		query := "SELECT * FROM test:$ WHERE name = $name AND id = $1"

		m.On("Query", mock.Anything, mock.Anything).Return(emptyResponse, nil)
		_, err := db.Exec(query, ID{"1"}, map[string]any{"name": "test"}, 1)

		assert.NoError(t, err)
		m.AssertCalled(t, "Query", "SELECT * FROM test:`1` WHERE name = $name AND id = $1;", map[string]any{
			"name": "test",
			"1":    1,
		})
	})
	t.Run("Query with a datetime id", func(t *testing.T) {
		m := new(mockQueryAgent)
		db := &DB{m}
		query := "SELECT * FROM test:$"

		m.On("Query", mock.Anything, mock.Anything).Return(emptyResponse, nil)
		now := time.Now()
		_, err := db.Exec(query, ID{now})

		assert.NoError(t, err)
		m.AssertCalled(t, "Query", fmt.Sprintf("SELECT * FROM test:`%s`;", now.Format(time.RFC3339)), map[string]any{})
	})
	t.Run("Query with a ranged date id", func(t *testing.T) {
		m := new(mockQueryAgent)
		db := &DB{m}
		query := "SELECT * FROM test:$"

		m.On("Query", mock.Anything, mock.Anything).Return(emptyResponse, nil)
		now := time.Now()
		_, err := db.Exec(query, Range{
			ID{"test", now},
			ID{"test", now.AddDate(0, 0, 1)},
		})

		assert.NoError(t, err)
		m.AssertCalled(t, "Query",
			fmt.Sprintf("SELECT * FROM test:['test', <datetime>'%s']..['test', <datetime>'%s'];",
				now.Format(time.RFC3339),
				now.AddDate(0, 0, 1).Format(time.RFC3339)),
			map[string]any{},
		)
	})
	t.Run("Query with a time and duration", func(t *testing.T) {
		m := new(mockQueryAgent)
		db := &DB{m}
		query := "CREATE test SET time = $1 + $2"

		m.On("Query", mock.Anything, mock.Anything).Return(emptyResponse, nil)
		now := time.Now()
		_, err := db.Exec(query, now, time.Hour)

		assert.NoError(t, err)
		m.AssertCalled(t, "Query", "CREATE test SET time = $1 + $2;", map[string]any{
			"1": now.Format(time.RFC3339),
			"2": "1h",
		})
	})
}
