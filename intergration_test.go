package surgo

import (
	"github.com/stretchr/testify/assert"
	"math/rand"
	"testing"
)

type testObject struct {
	TimeSince int `db:"time_since"`
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
DEFINE TABLE test;
DEFINE FIELD time_since ON TABLE test TYPE int;
DEFINE FIELD name       ON TABLE test TYPE string;
INSERT INTO test (time_since, name) VALUES ($1, $2);
INSERT INTO test (time_since, name) VALUES ($3, $4);
SELECT * FROM test;
`, ts, name, ts2, name2)

	t.Run("Scan only", func(t *testing.T) {
		only := testObject{}
		err := db.Scan(&only, "SELECT * FROM ONLY test WHERE time_since = $1", ts)
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
