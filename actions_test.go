package surgo

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

type SampleModel struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Age  int    `json:"age"`
	Sex  bool   `json:"sex"`
	// CreatedAt time.Time `json:"created_at"`
}

type SampleID struct {
	Pos int `json:"pos"`
	// At time.Time `json:"at"`
}

type DBMock struct {
	query string
}

func (db *DBMock) Query(query string) (interface{}, error) {
	db.query = query
	return nil, nil
}

var testDB = Model[SampleModel](&DBMock{})

func TestModel(t *testing.T) {
	t.Run("With a model", func(t *testing.T) {
		assert.Equal(t, "SampleModel", testDB.model)
	})
}

func TestSelect(t *testing.T) {
	t.Run("With a minimal select query", func(t *testing.T) {
		_, err := testDB.Select()
		assert.NoError(t, err)
		assert.Equal(t, "SELECT * FROM SampleModel;", testDB.db.(*DBMock).query)
	})
	t.Run("With an id", func(t *testing.T) {
		_, err := testDB.Select(ID("123"))
		assert.NoError(t, err)
		assert.Equal(t, "SELECT * FROM SampleModel:123;", testDB.db.(*DBMock).query)
	})
	t.Run("With a ranged id", func(t *testing.T) {
		_, err := testDB.Select(ID([2]any{SampleID{123}, SampleID{456}}))
		assert.NoError(t, err)
		assert.Equal(t, "SELECT * FROM SampleModel:[123]..[456];", testDB.db.(*DBMock).query)
	})
	t.Run("With fields", func(t *testing.T) {
		_, err := testDB.Select(Fields("id", "name"))
		assert.NoError(t, err)
		assert.Equal(t, "SELECT id, name FROM SampleModel;", testDB.db.(*DBMock).query)
	})
	t.Run("With omit", func(t *testing.T) {
		_, err := testDB.Select(Omit("id", "name"))
		assert.NoError(t, err)
		assert.Equal(t, "SELECT * OMIT id, name FROM SampleModel;", testDB.db.(*DBMock).query)
	})
	t.Run("With only", func(t *testing.T) {
		_, err := testDB.Select(Only())
		assert.NoError(t, err)
		assert.Equal(t, "SELECT * FROM ONLY SampleModel;", testDB.db.(*DBMock).query)
	})
	t.Run("With where", func(t *testing.T) {
		_, err := testDB.Select(Where(`name = "foo"`))
		assert.NoError(t, err)
		assert.Equal(t, `SELECT * FROM SampleModel WHERE name = "foo";`, testDB.db.(*DBMock).query)
	})
	t.Run("With group by", func(t *testing.T) {
		_, err := testDB.Select(GroupBy("name", "age"))
		assert.NoError(t, err)
		assert.Equal(t, "SELECT * FROM SampleModel GROUP BY name, age;", testDB.db.(*DBMock).query)
	})
	t.Run("With order by", func(t *testing.T) {
		_, err := testDB.Select(OrderBy(Asc("name"), Desc("age")))
		assert.NoError(t, err)
		assert.Equal(t, "SELECT * FROM SampleModel ORDER BY name ASC, age DESC;", testDB.db.(*DBMock).query)
	})
	t.Run("With limit", func(t *testing.T) {
		_, err := testDB.Select(Limit(10))
		assert.NoError(t, err)
		assert.Equal(t, "SELECT * FROM SampleModel LIMIT 10;", testDB.db.(*DBMock).query)
	})
	t.Run("With start", func(t *testing.T) {
		_, err := testDB.Select(Start(10))
		assert.NoError(t, err)
		assert.Equal(t, "SELECT * FROM SampleModel START 10;", testDB.db.(*DBMock).query)
	})
	t.Run("With fetch", func(t *testing.T) {
		_, err := testDB.Select(Fetch("group.sub"))
		assert.NoError(t, err)
		assert.Equal(t, "SELECT * FROM SampleModel FETCH group.sub;", testDB.db.(*DBMock).query)
	})
	t.Run("With timeout", func(t *testing.T) {
		_, err := testDB.Select(Timeout(10 * time.Second))
		assert.NoError(t, err)
		assert.Equal(t, "SELECT * FROM SampleModel TIMEOUT 10000ms;", testDB.db.(*DBMock).query)
	})
	t.Run("With parallel", func(t *testing.T) {
		_, err := testDB.Select(Parallel())
		assert.NoError(t, err)
		assert.Equal(t, "SELECT * FROM SampleModel PARALLEL;", testDB.db.(*DBMock).query)
	})
}
