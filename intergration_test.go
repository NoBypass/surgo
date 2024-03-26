package surgo

import (
	"github.com/stretchr/testify/assert"
	"math/rand"
	"testing"
)

type testObject struct {
	TimeSince int `DB:"time_since"`
	Name      string
}

func TestIntegration(t *testing.T) {
	db, err := Connect("127.0.0.1:8000",
		User("root"),
		Password("1234"),
		Database("test"),
		Namespace("test"),
	)
	if err != nil {
		t.Skip("SurrealDB not available")
	}
	defer db.Close()

	ts := rand.Intn(100)
	ts2 := rand.Intn(100)
	name := string(rune(rand.Intn(100)))
	name2 := string(rune(rand.Intn(100)))
	db.MustExec(`
REMOVE TABLE test;
REMOVE TABLE other;
DEFINE TABLE test;
DEFINE TABLE other;
DEFINE FIELD time_since ON TABLE test TYPE int;
DEFINE FIELD name       ON TABLE test TYPE string;
INSERT INTO test (time_since, name) VALUES ($1, $2);
INSERT INTO test (time_since, name) VALUES ($3, $4);
CREATE other:123;
CREATE other:456;
SELECT * FROM test;
`, ts, name, ts2, name2)

	t.Run("With id", func(t *testing.T) {
		var obj string
		err := db.Scan(&obj, "SELECT * FROM other:$", ID{123})
		assert.NoError(t, err)
	})
	t.Run("With ranged id", func(t *testing.T) {
		var obj []string
		err := db.Scan(&obj, "SELECT * FROM other:$", Range{ID{123}, ID{457}})
		assert.NoError(t, err)
	})
	t.Run("With string id", func(t *testing.T) {
		var obj string
		err := db.Scan(&obj, "SELECT * FROM other:$", ID{"123"})
		assert.NoError(t, err)
	})
	t.Run("With ranged string id", func(t *testing.T) {
		var obj []string
		err := db.Scan(&obj, "SELECT * FROM other:$", Range{ID{"123"}, ID{"457"}})
		assert.NoError(t, err)
	})
	t.Run("Scan only", func(t *testing.T) {
		only := testObject{}
		err := db.Scan(&only, "SELECT * FROM ONLY test WHERE time_since = $time_since", testObject{TimeSince: ts})
		assert.NoError(t, err)
		assert.Equal(t, testObject{TimeSince: ts, Name: name}, only)
	})
	t.Run("Scan only without only in query", func(t *testing.T) {
		failOnly := testObject{}
		err := db.Scan(&failOnly, "SELECT * FROM test")
		assert.NoError(t, err)
		assert.Equal(t, testObject{TimeSince: ts, Name: name}, failOnly)
	})
	t.Run("Scan only to slice", func(t *testing.T) {
		var only []testObject
		err := db.Scan(&only, "SELECT * FROM ONLY test WHERE time_since = $1", ts)
		assert.NoError(t, err)
		assert.Equal(t, []testObject{{TimeSince: ts, Name: name}}, only)
	})

	t.Run("Scan many", func(t *testing.T) {
		var many []testObject
		err := db.Scan(&many, "SELECT * FROM test")
		assert.NoError(t, err)
		assert.ElementsMatch(t, []testObject{{TimeSince: ts, Name: name}, {TimeSince: ts2, Name: name2}}, many)
	})

	t.Run("Scan to single variable", func(t *testing.T) {
		var count int
		err := db.Scan(&count, "SELECT count() FROM test GROUP ALL")
		assert.NoError(t, err)
		assert.Equal(t, 2, count)
	})
}
