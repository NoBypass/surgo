package surgo

import (
	"fmt"
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
}

type SampleRelation struct {
	At time.Time `surreal:"at"`
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

var testModel = Model[SampleModel](&DBMock{})
var testRelation = Relation[SampleModel, SampleModel, SampleRelation](&DBMock{})

func TestModel(t *testing.T) {
	t.Run("With a model", func(t *testing.T) {
		assert.Equal(t, "sample_model", testModel.model)
	})
}

func TestDBModel_Create(t *testing.T) {
	sampleModel := &SampleModel{ID: "123", Name: "foo", Age: 20, Bad: false}
	t.Run("With a model", func(t *testing.T) {
		err := testModel.Create(sampleModel)
		assert.NoError(t, err)
		assert.Equal(t,
			`CREATE sample_model CONTENT {name:"foo",age:20,bad:false};`,
			testModel.db.(*DBMock).query,
		)
	})
	t.Run("With an id", func(t *testing.T) {
		err := testModel.Create(sampleModel, ID("123"))
		assert.NoError(t, err)
		assert.Equal(t,
			`CREATE sample_model:123 CONTENT {name:"foo",age:20,bad:false};`,
			testModel.db.(*DBMock).query,
		)
	})
	t.Run("With a return", func(t *testing.T) {
		err := testModel.Create(sampleModel, Return(ReturnBefore))
		assert.NoError(t, err)
		assert.Equal(t,
			`CREATE sample_model CONTENT {name:"foo",age:20,bad:false} RETURN BEFORE;`,
			testModel.db.(*DBMock).query,
		)
	})
	t.Run("With a timeout", func(t *testing.T) {
		err := testModel.Create(sampleModel, Timeout(10*time.Second))
		assert.NoError(t, err)
		assert.Equal(t,
			`CREATE sample_model CONTENT {name:"foo",age:20,bad:false} TIMEOUT 10000ms;`,
			testModel.db.(*DBMock).query,
		)
	})
	t.Run("With a parallel", func(t *testing.T) {
		err := testModel.Create(sampleModel, Parallel())
		assert.NoError(t, err)
		assert.Equal(t,
			`CREATE sample_model CONTENT {name:"foo",age:20,bad:false} PARALLEL;`,
			testModel.db.(*DBMock).query,
		)
	})
}

func TestDBModel_Select(t *testing.T) {
	t.Run("With a minimal select query", func(t *testing.T) {
		test := SampleModel{}
		err := testModel.Select(&test)
		assert.NoError(t, err)
		assert.Equal(t, "SELECT * FROM sample_model;", testModel.db.(*DBMock).query)
	})
	t.Run("With an id", func(t *testing.T) {
		test := SampleModel{}
		err := testModel.Select(&test, ID("123"))
		assert.NoError(t, err)
		assert.Equal(t, "SELECT * FROM sample_model:123;", testModel.db.(*DBMock).query)
	})
	t.Run("With a ranged id", func(t *testing.T) {
		test := SampleModel{}
		err := testModel.Select(&test, ID([2]any{SampleID{123}, SampleID{456}}))
		assert.NoError(t, err)
		assert.Equal(t, "SELECT * FROM sample_model:[123]..[456];", testModel.db.(*DBMock).query)
	})
	t.Run("With fields", func(t *testing.T) {
		test := SampleModel{}
		err := testModel.Select(&test, Fields("id", "name"))
		assert.NoError(t, err)
		assert.Equal(t, "SELECT id, name FROM sample_model;", testModel.db.(*DBMock).query)
	})
	t.Run("With omit", func(t *testing.T) {
		test := SampleModel{}
		err := testModel.Select(&test, Omit("id", "name"))
		assert.NoError(t, err)
		assert.Equal(t, "SELECT * OMIT id, name FROM sample_model;", testModel.db.(*DBMock).query)
	})
	t.Run("With only", func(t *testing.T) {
		test := SampleModel{}
		err := testModel.Select(&test, Only())
		assert.NoError(t, err)
		assert.Equal(t, "SELECT * FROM ONLY sample_model;", testModel.db.(*DBMock).query)
	})
	t.Run("With where", func(t *testing.T) {
		test := SampleModel{}
		err := testModel.Select(&test, Where(`name = "foo"`))
		assert.NoError(t, err)
		assert.Equal(t, `SELECT * FROM sample_model WHERE name = "foo";`, testModel.db.(*DBMock).query)
	})
	t.Run("With group by", func(t *testing.T) {
		test := SampleModel{}
		err := testModel.Select(&test, GroupBy("name", "age"))
		assert.NoError(t, err)
		assert.Equal(t, "SELECT * FROM sample_model GROUP BY name, age;", testModel.db.(*DBMock).query)
	})
	t.Run("With order by", func(t *testing.T) {
		test := SampleModel{}
		err := testModel.Select(&test, OrderBy(Asc("name"), Desc("age")))
		assert.NoError(t, err)
		assert.Equal(t, "SELECT * FROM sample_model ORDER BY name ASC, age DESC;", testModel.db.(*DBMock).query)
	})
	t.Run("With limit", func(t *testing.T) {
		test := SampleModel{}
		err := testModel.Select(&test, Limit(10))
		assert.NoError(t, err)
		assert.Equal(t, "SELECT * FROM sample_model LIMIT 10;", testModel.db.(*DBMock).query)
	})
	t.Run("With start", func(t *testing.T) {
		test := SampleModel{}
		err := testModel.Select(&test, Start(10))
		assert.NoError(t, err)
		assert.Equal(t, "SELECT * FROM sample_model START 10;", testModel.db.(*DBMock).query)
	})
	t.Run("With fetch", func(t *testing.T) {
		test := SampleModel{}
		err := testModel.Select(&test, Fetch("group.sub"))
		assert.NoError(t, err)
		assert.Equal(t, "SELECT * FROM sample_model FETCH group.sub;", testModel.db.(*DBMock).query)
	})
	t.Run("With timeout", func(t *testing.T) {
		test := SampleModel{}
		err := testModel.Select(&test, Timeout(10*time.Second))
		assert.NoError(t, err)
		assert.Equal(t, "SELECT * FROM sample_model TIMEOUT 10000ms;", testModel.db.(*DBMock).query)
	})
	t.Run("With parallel", func(t *testing.T) {
		test := SampleModel{}
		err := testModel.Select(&test, Parallel())
		assert.NoError(t, err)
		assert.Equal(t, "SELECT * FROM sample_model PARALLEL;", testModel.db.(*DBMock).query)
	})
}

func TestDBModel_Update(t *testing.T) {
	sampleModel := &SampleModel{ID: "123", Name: "foo", Age: 20, Bad: false}
	t.Run("With a model", func(t *testing.T) {
		err := testModel.Update(sampleModel)
		assert.NoError(t, err)
		assert.Equal(t,
			`UPDATE sample_model CONTENT {name:"foo",age:20,bad:false};`,
			testModel.db.(*DBMock).query,
		)
	})
	t.Run("With an id", func(t *testing.T) {
		err := testModel.Update(sampleModel, ID("123"))
		assert.NoError(t, err)
		assert.Equal(t,
			`UPDATE sample_model:123 CONTENT {name:"foo",age:20,bad:false};`,
			testModel.db.(*DBMock).query,
		)
	})
	t.Run("With a return", func(t *testing.T) {
		err := testModel.Update(sampleModel, Return(ReturnBefore))
		assert.NoError(t, err)
		assert.Equal(t,
			`UPDATE sample_model CONTENT {name:"foo",age:20,bad:false} RETURN BEFORE;`,
			testModel.db.(*DBMock).query,
		)
	})
	t.Run("With a timeout", func(t *testing.T) {
		err := testModel.Update(sampleModel, Timeout(10*time.Second))
		assert.NoError(t, err)
		assert.Equal(t,
			`UPDATE sample_model CONTENT {name:"foo",age:20,bad:false} TIMEOUT 10000ms;`,
			testModel.db.(*DBMock).query,
		)
	})
	t.Run("With a parallel", func(t *testing.T) {
		err := testModel.Update(sampleModel, Parallel())
		assert.NoError(t, err)
		assert.Equal(t,
			`UPDATE sample_model CONTENT {name:"foo",age:20,bad:false} PARALLEL;`,
			testModel.db.(*DBMock).query,
		)
	})
}

func TestDBModel_Delete(t *testing.T) {
	t.Run("With a minimal delete query", func(t *testing.T) {
		_, err := testModel.Delete()
		assert.NoError(t, err)
		assert.Equal(t, "DELETE sample_model;", testModel.db.(*DBMock).query)
	})
	t.Run("With an id", func(t *testing.T) {
		_, err := testModel.Delete(ID("123"))
		assert.NoError(t, err)
		assert.Equal(t, "DELETE sample_model:123;", testModel.db.(*DBMock).query)
	})
	t.Run("With a return", func(t *testing.T) {
		_, err := testModel.Delete(Return(ReturnBefore))
		assert.NoError(t, err)
		assert.Equal(t, "DELETE sample_model RETURN BEFORE;", testModel.db.(*DBMock).query)
	})
	t.Run("With a timeout", func(t *testing.T) {
		_, err := testModel.Delete(Timeout(10 * time.Second))
		assert.NoError(t, err)
		assert.Equal(t, "DELETE sample_model TIMEOUT 10000ms;", testModel.db.(*DBMock).query)
	})
	t.Run("With a parallel", func(t *testing.T) {
		_, err := testModel.Delete(Parallel())
		assert.NoError(t, err)
		assert.Equal(t, "DELETE sample_model PARALLEL;", testModel.db.(*DBMock).query)
	})
	t.Run("With a where", func(t *testing.T) {
		_, err := testModel.Delete(Where(`name = "foo"`))
		assert.NoError(t, err)
		assert.Equal(t, `DELETE sample_model WHERE name = "foo";`, testModel.db.(*DBMock).query)
	})
	t.Run("With an only", func(t *testing.T) {
		_, err := testModel.Delete(Only())
		assert.NoError(t, err)
		assert.Equal(t, "DELETE ONLY sample_model;", testModel.db.(*DBMock).query)
	})
}

func TestDBRelation_Create(t *testing.T) {
	t.Run("With a minimal create query", func(t *testing.T) {
		err := testRelation.Create(nil, ID("123"), ID("456"))
		assert.NoError(t, err)
		assert.Equal(t, "RELATE sample_model:123->sample_relation->sample_model:456;", testRelation.db.(*DBMock).query)
	})
	t.Run("With a return", func(t *testing.T) {
		err := testRelation.Create(nil, ID("123"), ID("456"), Return(ReturnBefore))
		assert.NoError(t, err)
		assert.Equal(t, "RELATE sample_model:123->sample_relation->sample_model:456 RETURN BEFORE;", testRelation.db.(*DBMock).query)
	})
	t.Run("With a timeout", func(t *testing.T) {
		err := testRelation.Create(nil, ID("123"), ID("456"), Timeout(10*time.Second))
		assert.NoError(t, err)
		assert.Equal(t, "RELATE sample_model:123->sample_relation->sample_model:456 TIMEOUT 10000ms;", testRelation.db.(*DBMock).query)
	})
	t.Run("With a parallel", func(t *testing.T) {
		err := testRelation.Create(nil, ID("123"), ID("456"), Parallel())
		assert.NoError(t, err)
		assert.Equal(t, "RELATE sample_model:123->sample_relation->sample_model:456 PARALLEL;", testRelation.db.(*DBMock).query)
	})
	t.Run("With an only", func(t *testing.T) {
		err := testRelation.Create(nil, ID("123"), ID("456"), Only())
		assert.NoError(t, err)
		assert.Equal(t, "RELATE ONLY sample_model:123->sample_relation->sample_model:456;", testRelation.db.(*DBMock).query)
	})
	t.Run("With content", func(t *testing.T) {
		stamp := time.Now()
		err := testRelation.Create(&SampleRelation{At: stamp}, ID("123"), ID("456"))
		assert.NoError(t, err)
		assert.Equal(t, fmt.Sprintf(`RELATE sample_model:123->sample_relation->sample_model:456 CONTENT {at:"%s"};`, time.Now().Format(time.RFC3339)),
			testRelation.db.(*DBMock).query)
	})
}

func TestDBRelation_Delete(t *testing.T) {
	t.Run("With a minimal delete query", func(t *testing.T) {
		_, err := testRelation.Delete(ID("123"), ID("456"))
		assert.NoError(t, err)
		assert.Equal(t, "DELETE sample_model:123->sample_relation WHERE out=sample_model:456;", testRelation.db.(*DBMock).query)
	})
	t.Run("With a return", func(t *testing.T) {
		_, err := testRelation.Delete(ID("123"), ID("456"), Return(ReturnBefore))
		assert.NoError(t, err)
		assert.Equal(t, "DELETE sample_model:123->sample_relation WHERE out=sample_model:456 RETURN BEFORE;", testRelation.db.(*DBMock).query)
	})
	t.Run("With a timeout", func(t *testing.T) {
		_, err := testRelation.Delete(ID("123"), ID("456"), Timeout(10*time.Second))
		assert.NoError(t, err)
		assert.Equal(t, "DELETE sample_model:123->sample_relation WHERE out=sample_model:456 TIMEOUT 10000ms;", testRelation.db.(*DBMock).query)
	})
	t.Run("With a parallel", func(t *testing.T) {
		_, err := testRelation.Delete(ID("123"), ID("456"), Parallel())
		assert.NoError(t, err)
		assert.Equal(t, "DELETE sample_model:123->sample_relation WHERE out=sample_model:456 PARALLEL;", testRelation.db.(*DBMock).query)
	})
	t.Run("With an only", func(t *testing.T) {
		_, err := testRelation.Delete(ID("123"), ID("456"), Only())
		assert.NoError(t, err)
		assert.Equal(t, "DELETE ONLY sample_model:123->sample_relation WHERE out=sample_model:456;", testRelation.db.(*DBMock).query)
	})
	t.Run("With a where", func(t *testing.T) {
		_, err := testRelation.Delete(ID("123"), ID("456"), Where(`at = 0`))
		assert.NoError(t, err)
		assert.Equal(t, "DELETE sample_model:123->sample_relation WHERE at = 0 AND out=sample_model:456;", testRelation.db.(*DBMock).query)
	})
}
