package surgo

import (
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
	"time"
)

type SampleModel struct {
	ID   string `surreal:"id"`
	Name string `surreal:"name"`
	Age  int    `surreal:"age"`
	Bad  bool   `surreal:"bad"`
	// CreatedAt time.Time `json:"created_at"`
}

type SampleID struct {
	Pos int
	// At time.Time `json:"at"`
}

type DBMock struct {
	query string
}

func (db *DBMock) Query(query string) (interface{}, error) {
	db.query = strings.TrimSpace(query) + ";"
	return nil, nil
}

var testDB = Model[SampleModel](&DBMock{})

func TestModel(t *testing.T) {
	t.Run("With a model", func(t *testing.T) {
		assert.Equal(t, "SampleModel", testDB.model)
	})
}

func TestCreate(t *testing.T) {
	sampleModel := &SampleModel{ID: "123", Name: "foo", Age: 20, Bad: false}
	t.Run("With a model", func(t *testing.T) {
		err := testDB.Create(sampleModel)
		assert.NoError(t, err)
		assert.Equal(t,
			`CREATE SampleModel CONTENT {id:"123",name:"foo",age:20,bad:false};`,
			testDB.db.(*DBMock).query,
		)
	})
	t.Run("With an id", func(t *testing.T) {
		err := testDB.Create(sampleModel, ID("123"))
		assert.NoError(t, err)
		assert.Equal(t,
			`CREATE SampleModel:123 CONTENT {id:"123",name:"foo",age:20,bad:false};`,
			testDB.db.(*DBMock).query,
		)
	})
	t.Run("With a return", func(t *testing.T) {
		err := testDB.Create(sampleModel, Return(ReturnBefore))
		assert.NoError(t, err)
		assert.Equal(t,
			`CREATE SampleModel CONTENT {id:"123",name:"foo",age:20,bad:false} RETURN BEFORE;`,
			testDB.db.(*DBMock).query,
		)
	})
	t.Run("With a timeout", func(t *testing.T) {
		err := testDB.Create(sampleModel, Timeout(10*time.Second))
		assert.NoError(t, err)
		assert.Equal(t,
			`CREATE SampleModel CONTENT {id:"123",name:"foo",age:20,bad:false} TIMEOUT 10000ms;`,
			testDB.db.(*DBMock).query,
		)
	})
	t.Run("With a parallel", func(t *testing.T) {
		err := testDB.Create(sampleModel, Parallel())
		assert.NoError(t, err)
		assert.Equal(t,
			`CREATE SampleModel CONTENT {id:"123",name:"foo",age:20,bad:false} PARALLEL;`,
			testDB.db.(*DBMock).query,
		)
	})
}

func TestSelect(t *testing.T) {
	t.Run("With a minimal select query", func(t *testing.T) {
		test := SampleModel{}
		err := testDB.Select(&test)
		assert.NoError(t, err)
		assert.Equal(t, "SELECT * FROM SampleModel;", testDB.db.(*DBMock).query)
	})
	t.Run("With an id", func(t *testing.T) {
		test := SampleModel{}
		err := testDB.Select(&test, ID("123"))
		assert.NoError(t, err)
		assert.Equal(t, "SELECT * FROM SampleModel:123;", testDB.db.(*DBMock).query)
	})
	t.Run("With a ranged id", func(t *testing.T) {
		test := SampleModel{}
		err := testDB.Select(&test, ID([2]any{SampleID{123}, SampleID{456}}))
		assert.NoError(t, err)
		assert.Equal(t, "SELECT * FROM SampleModel:[123]..[456];", testDB.db.(*DBMock).query)
	})
	t.Run("With fields", func(t *testing.T) {
		test := SampleModel{}
		err := testDB.Select(&test, Fields("id", "name"))
		assert.NoError(t, err)
		assert.Equal(t, "SELECT id, name FROM SampleModel;", testDB.db.(*DBMock).query)
	})
	t.Run("With omit", func(t *testing.T) {
		test := SampleModel{}
		err := testDB.Select(&test, Omit("id", "name"))
		assert.NoError(t, err)
		assert.Equal(t, "SELECT * OMIT id, name FROM SampleModel;", testDB.db.(*DBMock).query)
	})
	t.Run("With only", func(t *testing.T) {
		test := SampleModel{}
		err := testDB.Select(&test, Only())
		assert.NoError(t, err)
		assert.Equal(t, "SELECT * FROM ONLY SampleModel;", testDB.db.(*DBMock).query)
	})
	t.Run("With where", func(t *testing.T) {
		test := SampleModel{}
		err := testDB.Select(&test, Where(`name = "foo"`))
		assert.NoError(t, err)
		assert.Equal(t, `SELECT * FROM SampleModel WHERE name = "foo";`, testDB.db.(*DBMock).query)
	})
	t.Run("With group by", func(t *testing.T) {
		test := SampleModel{}
		err := testDB.Select(&test, GroupBy("name", "age"))
		assert.NoError(t, err)
		assert.Equal(t, "SELECT * FROM SampleModel GROUP BY name, age;", testDB.db.(*DBMock).query)
	})
	t.Run("With order by", func(t *testing.T) {
		test := SampleModel{}
		err := testDB.Select(&test, OrderBy(Asc("name"), Desc("age")))
		assert.NoError(t, err)
		assert.Equal(t, "SELECT * FROM SampleModel ORDER BY name ASC, age DESC;", testDB.db.(*DBMock).query)
	})
	t.Run("With limit", func(t *testing.T) {
		test := SampleModel{}
		err := testDB.Select(&test, Limit(10))
		assert.NoError(t, err)
		assert.Equal(t, "SELECT * FROM SampleModel LIMIT 10;", testDB.db.(*DBMock).query)
	})
	t.Run("With start", func(t *testing.T) {
		test := SampleModel{}
		err := testDB.Select(&test, Start(10))
		assert.NoError(t, err)
		assert.Equal(t, "SELECT * FROM SampleModel START 10;", testDB.db.(*DBMock).query)
	})
	t.Run("With fetch", func(t *testing.T) {
		test := SampleModel{}
		err := testDB.Select(&test, Fetch("group.sub"))
		assert.NoError(t, err)
		assert.Equal(t, "SELECT * FROM SampleModel FETCH group.sub;", testDB.db.(*DBMock).query)
	})
	t.Run("With timeout", func(t *testing.T) {
		test := SampleModel{}
		err := testDB.Select(&test, Timeout(10*time.Second))
		assert.NoError(t, err)
		assert.Equal(t, "SELECT * FROM SampleModel TIMEOUT 10000ms;", testDB.db.(*DBMock).query)
	})
	t.Run("With parallel", func(t *testing.T) {
		test := SampleModel{}
		err := testDB.Select(&test, Parallel())
		assert.NoError(t, err)
		assert.Equal(t, "SELECT * FROM SampleModel PARALLEL;", testDB.db.(*DBMock).query)
	})
}

func TestUpdate(t *testing.T) {
	sampleModel := &SampleModel{ID: "123", Name: "foo", Age: 20, Bad: false}
	t.Run("With a model", func(t *testing.T) {
		err := testDB.Update(sampleModel)
		assert.NoError(t, err)
		assert.Equal(t,
			`UPDATE SampleModel CONTENT {id:"123",name:"foo",age:20,bad:false};`,
			testDB.db.(*DBMock).query,
		)
	})
	t.Run("With an id", func(t *testing.T) {
		err := testDB.Update(sampleModel, ID("123"))
		assert.NoError(t, err)
		assert.Equal(t,
			`UPDATE SampleModel:123 CONTENT {id:"123",name:"foo",age:20,bad:false};`,
			testDB.db.(*DBMock).query,
		)
	})
	t.Run("With a return", func(t *testing.T) {
		err := testDB.Update(sampleModel, Return(ReturnBefore))
		assert.NoError(t, err)
		assert.Equal(t,
			`UPDATE SampleModel CONTENT {id:"123",name:"foo",age:20,bad:false} RETURN BEFORE;`,
			testDB.db.(*DBMock).query,
		)
	})
	t.Run("With a timeout", func(t *testing.T) {
		err := testDB.Update(sampleModel, Timeout(10*time.Second))
		assert.NoError(t, err)
		assert.Equal(t,
			`UPDATE SampleModel CONTENT {id:"123",name:"foo",age:20,bad:false} TIMEOUT 10000ms;`,
			testDB.db.(*DBMock).query,
		)
	})
	t.Run("With a parallel", func(t *testing.T) {
		err := testDB.Update(sampleModel, Parallel())
		assert.NoError(t, err)
		assert.Equal(t,
			`UPDATE SampleModel CONTENT {id:"123",name:"foo",age:20,bad:false} PARALLEL;`,
			testDB.db.(*DBMock).query,
		)
	})
}

func TestDelete(t *testing.T) {
	t.Run("With a minimal delete query", func(t *testing.T) {
		_, err := testDB.Delete("123")
		assert.NoError(t, err)
		assert.Equal(t, "DELETE SampleModel:123;", testDB.db.(*DBMock).query)
	})
	t.Run("With an id", func(t *testing.T) {
		_, err := testDB.Delete("123", ID("123"))
		assert.NoError(t, err)
		assert.Equal(t, "DELETE SampleModel:123;", testDB.db.(*DBMock).query)
	})
	t.Run("With a return", func(t *testing.T) {
		_, err := testDB.Delete("123", Return(ReturnBefore))
		assert.NoError(t, err)
		assert.Equal(t, "DELETE SampleModel:123 RETURN BEFORE;", testDB.db.(*DBMock).query)
	})
	t.Run("With a timeout", func(t *testing.T) {
		_, err := testDB.Delete("123", Timeout(10*time.Second))
		assert.NoError(t, err)
		assert.Equal(t, "DELETE SampleModel:123 TIMEOUT 10000ms;", testDB.db.(*DBMock).query)
	})
	t.Run("With a parallel", func(t *testing.T) {
		_, err := testDB.Delete("123", Parallel())
		assert.NoError(t, err)
		assert.Equal(t, "DELETE SampleModel:123 PARALLEL;", testDB.db.(*DBMock).query)
	})
	t.Run("With a where", func(t *testing.T) {
		_, err := testDB.Delete("123", Where(`name = "foo"`))
		assert.NoError(t, err)
		assert.Equal(t, `DELETE SampleModel:123 WHERE name = "foo";`, testDB.db.(*DBMock).query)
	})
	t.Run("With an only", func(t *testing.T) {
		_, err := testDB.Delete("123", Only())
		assert.NoError(t, err)
		assert.Equal(t, "DELETE ONLY SampleModel:123;", testDB.db.(*DBMock).query)
	})
}
