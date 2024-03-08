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

var model = Model[SampleModel](&DBMock{})
var relation = Relation[SampleModel, SampleModel, SampleRelation](&DBMock{})

func TestModel(t *testing.T) {
	t.Run("With a model", func(t *testing.T) {
		assert.Equal(t, "sample_model", model.model)
	})
}

func TestDBModel_Create(t *testing.T) {
	sampleModel := &SampleModel{ID: "123", Name: "foo", Age: 20, Bad: false}
	t.Run("With a model", func(t *testing.T) {
		err := model.Create(sampleModel)
		assert.NoError(t, err)
		assert.Equal(t,
			`CREATE sample_model CONTENT {id:"123",name:"foo",age:20,bad:false};`,
			model.db.(*DBMock).query,
		)
	})
	t.Run("With an id", func(t *testing.T) {
		err := model.Create(sampleModel, ID("123"))
		assert.NoError(t, err)
		assert.Equal(t,
			`CREATE sample_model:123 CONTENT {id:"123",name:"foo",age:20,bad:false};`,
			model.db.(*DBMock).query,
		)
	})
	t.Run("With a return", func(t *testing.T) {
		err := model.Create(sampleModel, Return(ReturnBefore))
		assert.NoError(t, err)
		assert.Equal(t,
			`CREATE sample_model CONTENT {id:"123",name:"foo",age:20,bad:false} RETURN BEFORE;`,
			model.db.(*DBMock).query,
		)
	})
	t.Run("With a timeout", func(t *testing.T) {
		err := model.Create(sampleModel, Timeout(10*time.Second))
		assert.NoError(t, err)
		assert.Equal(t,
			`CREATE sample_model CONTENT {id:"123",name:"foo",age:20,bad:false} TIMEOUT 10000ms;`,
			model.db.(*DBMock).query,
		)
	})
	t.Run("With a parallel", func(t *testing.T) {
		err := model.Create(sampleModel, Parallel())
		assert.NoError(t, err)
		assert.Equal(t,
			`CREATE sample_model CONTENT {id:"123",name:"foo",age:20,bad:false} PARALLEL;`,
			model.db.(*DBMock).query,
		)
	})
}

func TestDBModel_Select(t *testing.T) {
	t.Run("With a minimal select query", func(t *testing.T) {
		test := SampleModel{}
		err := model.Select(&test)
		assert.NoError(t, err)
		assert.Equal(t, "SELECT * FROM sample_model;", model.db.(*DBMock).query)
	})
	t.Run("With an id", func(t *testing.T) {
		test := SampleModel{}
		err := model.Select(&test, ID("123"))
		assert.NoError(t, err)
		assert.Equal(t, "SELECT * FROM sample_model:123;", model.db.(*DBMock).query)
	})
	t.Run("With a ranged id", func(t *testing.T) {
		test := SampleModel{}
		err := model.Select(&test, ID([2]any{SampleID{123}, SampleID{456}}))
		assert.NoError(t, err)
		assert.Equal(t, "SELECT * FROM sample_model:[123]..[456];", model.db.(*DBMock).query)
	})
	t.Run("With fields", func(t *testing.T) {
		test := SampleModel{}
		err := model.Select(&test, Fields("id", "name"))
		assert.NoError(t, err)
		assert.Equal(t, "SELECT id, name FROM sample_model;", model.db.(*DBMock).query)
	})
	t.Run("With omit", func(t *testing.T) {
		test := SampleModel{}
		err := model.Select(&test, Omit("id", "name"))
		assert.NoError(t, err)
		assert.Equal(t, "SELECT * OMIT id, name FROM sample_model;", model.db.(*DBMock).query)
	})
	t.Run("With only", func(t *testing.T) {
		test := SampleModel{}
		err := model.Select(&test, Only())
		assert.NoError(t, err)
		assert.Equal(t, "SELECT * FROM ONLY sample_model;", model.db.(*DBMock).query)
	})
	t.Run("With where", func(t *testing.T) {
		test := SampleModel{}
		err := model.Select(&test, Where(`name = "foo"`))
		assert.NoError(t, err)
		assert.Equal(t, `SELECT * FROM sample_model WHERE name = "foo";`, model.db.(*DBMock).query)
	})
	t.Run("With group by", func(t *testing.T) {
		test := SampleModel{}
		err := model.Select(&test, GroupBy("name", "age"))
		assert.NoError(t, err)
		assert.Equal(t, "SELECT * FROM sample_model GROUP BY name, age;", model.db.(*DBMock).query)
	})
	t.Run("With order by", func(t *testing.T) {
		test := SampleModel{}
		err := model.Select(&test, OrderBy(Asc("name"), Desc("age")))
		assert.NoError(t, err)
		assert.Equal(t, "SELECT * FROM sample_model ORDER BY name ASC, age DESC;", model.db.(*DBMock).query)
	})
	t.Run("With limit", func(t *testing.T) {
		test := SampleModel{}
		err := model.Select(&test, Limit(10))
		assert.NoError(t, err)
		assert.Equal(t, "SELECT * FROM sample_model LIMIT 10;", model.db.(*DBMock).query)
	})
	t.Run("With start", func(t *testing.T) {
		test := SampleModel{}
		err := model.Select(&test, Start(10))
		assert.NoError(t, err)
		assert.Equal(t, "SELECT * FROM sample_model START 10;", model.db.(*DBMock).query)
	})
	t.Run("With fetch", func(t *testing.T) {
		test := SampleModel{}
		err := model.Select(&test, Fetch("group.sub"))
		assert.NoError(t, err)
		assert.Equal(t, "SELECT * FROM sample_model FETCH group.sub;", model.db.(*DBMock).query)
	})
	t.Run("With timeout", func(t *testing.T) {
		test := SampleModel{}
		err := model.Select(&test, Timeout(10*time.Second))
		assert.NoError(t, err)
		assert.Equal(t, "SELECT * FROM sample_model TIMEOUT 10000ms;", model.db.(*DBMock).query)
	})
	t.Run("With parallel", func(t *testing.T) {
		test := SampleModel{}
		err := model.Select(&test, Parallel())
		assert.NoError(t, err)
		assert.Equal(t, "SELECT * FROM sample_model PARALLEL;", model.db.(*DBMock).query)
	})
}

func TestDBModel_Update(t *testing.T) {
	sampleModel := &SampleModel{ID: "123", Name: "foo", Age: 20, Bad: false}
	t.Run("With a model", func(t *testing.T) {
		err := model.Update(sampleModel)
		assert.NoError(t, err)
		assert.Equal(t,
			`UPDATE sample_model CONTENT {id:"123",name:"foo",age:20,bad:false};`,
			model.db.(*DBMock).query,
		)
	})
	t.Run("With an id", func(t *testing.T) {
		err := model.Update(sampleModel, ID("123"))
		assert.NoError(t, err)
		assert.Equal(t,
			`UPDATE sample_model:123 CONTENT {id:"123",name:"foo",age:20,bad:false};`,
			model.db.(*DBMock).query,
		)
	})
	t.Run("With a return", func(t *testing.T) {
		err := model.Update(sampleModel, Return(ReturnBefore))
		assert.NoError(t, err)
		assert.Equal(t,
			`UPDATE sample_model CONTENT {id:"123",name:"foo",age:20,bad:false} RETURN BEFORE;`,
			model.db.(*DBMock).query,
		)
	})
	t.Run("With a timeout", func(t *testing.T) {
		err := model.Update(sampleModel, Timeout(10*time.Second))
		assert.NoError(t, err)
		assert.Equal(t,
			`UPDATE sample_model CONTENT {id:"123",name:"foo",age:20,bad:false} TIMEOUT 10000ms;`,
			model.db.(*DBMock).query,
		)
	})
	t.Run("With a parallel", func(t *testing.T) {
		err := model.Update(sampleModel, Parallel())
		assert.NoError(t, err)
		assert.Equal(t,
			`UPDATE sample_model CONTENT {id:"123",name:"foo",age:20,bad:false} PARALLEL;`,
			model.db.(*DBMock).query,
		)
	})
}

func TestDBModel_Delete(t *testing.T) {
	t.Run("With a minimal delete query", func(t *testing.T) {
		_, err := model.Delete("123")
		assert.NoError(t, err)
		assert.Equal(t, "DELETE sample_model:123;", model.db.(*DBMock).query)
	})
	t.Run("With an id", func(t *testing.T) {
		_, err := model.Delete("123", ID("123"))
		assert.NoError(t, err)
		assert.Equal(t, "DELETE sample_model:123;", model.db.(*DBMock).query)
	})
	t.Run("With a return", func(t *testing.T) {
		_, err := model.Delete("123", Return(ReturnBefore))
		assert.NoError(t, err)
		assert.Equal(t, "DELETE sample_model:123 RETURN BEFORE;", model.db.(*DBMock).query)
	})
	t.Run("With a timeout", func(t *testing.T) {
		_, err := model.Delete("123", Timeout(10*time.Second))
		assert.NoError(t, err)
		assert.Equal(t, "DELETE sample_model:123 TIMEOUT 10000ms;", model.db.(*DBMock).query)
	})
	t.Run("With a parallel", func(t *testing.T) {
		_, err := model.Delete("123", Parallel())
		assert.NoError(t, err)
		assert.Equal(t, "DELETE sample_model:123 PARALLEL;", model.db.(*DBMock).query)
	})
	t.Run("With a where", func(t *testing.T) {
		_, err := model.Delete("123", Where(`name = "foo"`))
		assert.NoError(t, err)
		assert.Equal(t, `DELETE sample_model:123 WHERE name = "foo";`, model.db.(*DBMock).query)
	})
	t.Run("With an only", func(t *testing.T) {
		_, err := model.Delete("123", Only())
		assert.NoError(t, err)
		assert.Equal(t, "DELETE ONLY sample_model:123;", model.db.(*DBMock).query)
	})
}

func TestDBRelation_Create(t *testing.T) {
	t.Run("With a minimal create query", func(t *testing.T) {
		err := relation.Create(nil, ID("123"), ID("456"))
		assert.NoError(t, err)
		assert.Equal(t, "RELATE sample_model:123->sample_relation->sample_model:456;", relation.db.(*DBMock).query)
	})
	t.Run("With a return", func(t *testing.T) {
		err := relation.Create(nil, ID("123"), ID("456"), Return(ReturnBefore))
		assert.NoError(t, err)
		assert.Equal(t, "RELATE sample_model:123->sample_relation->sample_model:456 RETURN BEFORE;", relation.db.(*DBMock).query)
	})
	t.Run("With a timeout", func(t *testing.T) {
		err := relation.Create(nil, ID("123"), ID("456"), Timeout(10*time.Second))
		assert.NoError(t, err)
		assert.Equal(t, "RELATE sample_model:123->sample_relation->sample_model:456 TIMEOUT 10000ms;", relation.db.(*DBMock).query)
	})
	t.Run("With a parallel", func(t *testing.T) {
		err := relation.Create(nil, ID("123"), ID("456"), Parallel())
		assert.NoError(t, err)
		assert.Equal(t, "RELATE sample_model:123->sample_relation->sample_model:456 PARALLEL;", relation.db.(*DBMock).query)
	})
	t.Run("With an only", func(t *testing.T) {
		err := relation.Create(nil, ID("123"), ID("456"), Only())
		assert.NoError(t, err)
		assert.Equal(t, "RELATE ONLY sample_model:123->sample_relation->sample_model:456;", relation.db.(*DBMock).query)
	})
	t.Run("With content", func(t *testing.T) {
		stamp := time.Now()
		err := relation.Create(&SampleRelation{At: stamp}, ID("123"), ID("456"))
		assert.NoError(t, err)
		assert.Equal(t, fmt.Sprintf(`RELATE sample_model:123->sample_relation->sample_model:456 CONTENT {at:"%s"};`, time.Now().Format(time.RFC3339)),
			relation.db.(*DBMock).query)
	})
}

func TestDBRelation_Delete(t *testing.T) {
	t.Run("With a minimal delete query", func(t *testing.T) {
		_, err := relation.Delete(ID("123"), ID("456"))
		assert.NoError(t, err)
		assert.Equal(t, "DELETE sample_model:123->sample_relation WHERE out=sample_model:456;", relation.db.(*DBMock).query)
	})
	t.Run("With a return", func(t *testing.T) {
		_, err := relation.Delete(ID("123"), ID("456"), Return(ReturnBefore))
		assert.NoError(t, err)
		assert.Equal(t, "DELETE sample_model:123->sample_relation WHERE out=sample_model:456 RETURN BEFORE;", relation.db.(*DBMock).query)
	})
	t.Run("With a timeout", func(t *testing.T) {
		_, err := relation.Delete(ID("123"), ID("456"), Timeout(10*time.Second))
		assert.NoError(t, err)
		assert.Equal(t, "DELETE sample_model:123->sample_relation WHERE out=sample_model:456 TIMEOUT 10000ms;", relation.db.(*DBMock).query)
	})
	t.Run("With a parallel", func(t *testing.T) {
		_, err := relation.Delete(ID("123"), ID("456"), Parallel())
		assert.NoError(t, err)
		assert.Equal(t, "DELETE sample_model:123->sample_relation WHERE out=sample_model:456 PARALLEL;", relation.db.(*DBMock).query)
	})
	t.Run("With an only", func(t *testing.T) {
		_, err := relation.Delete(ID("123"), ID("456"), Only())
		assert.NoError(t, err)
		assert.Equal(t, "DELETE ONLY sample_model:123->sample_relation WHERE out=sample_model:456;", relation.db.(*DBMock).query)
	})
	t.Run("With a where", func(t *testing.T) {
		_, err := relation.Delete(ID("123"), ID("456"), Where(`at = 0`))
		assert.NoError(t, err)
		assert.Equal(t, "DELETE sample_model:123->sample_relation WHERE out=sample_model:456 AND at = 0;", relation.db.(*DBMock).query)
	})
}
