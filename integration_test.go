package surgo

import (
	"context"
	"testing"
	"time"
)

type TestObj struct {
	Title    string
	Num      int
	Floating float64
	Boolean  bool
	Time     time.Time
	Duration time.Duration
}

type SecondTest struct {
	ID     int    `db:"id"`
	Thing  string `db:"test_tag"`
	Ignore string `db:"-"`
	Omit   string `db:"omit,omitempty"`
}

func TestIntegration(t *testing.T) {
	db, err := Connect("ws://localhost:8000/rpc", &Credentials{
		Username:  "admin",
		Password:  "admin",
		Namespace: "test",
		Database:  "test",
	})
	if err != nil {
		t.Skipf("skipping integration test %v", err)
	}

	defer db.Close()

	result := db.Query("CREATE test:test CONTENT $test", map[string]any{
		"test": TestObj{
			Title:    "test",
			Num:      1,
			Floating: 1.1,
			Boolean:  true,
			Time:     time.Now(),
			Duration: time.Second + time.Millisecond*500,
		},
	})
	if result.Error != nil {
		t.Errorf("unexpected error: %v", err)
	} else {
		r, err := result.First()
		t.Logf("result: %v | %v", r, err)
	}

	var test = new(TestObj)
	err = db.Scan(test, "SELECT * FROM ONLY test:test", map[string]any{})
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	} else {
		t.Logf("result: %v", test)
	}

	var test2 = new(TestObj)
	err = db.Scan(test2, "SELECT * FROM ONLY test:unavailable", map[string]any{})
	if err != nil {
		t.Logf("expected error: %v", err)
	} else {
		t.Errorf("result: %v", test2)
	}

	result2 := db.Query("INSERT INTO test $tests RETURN AFTER", map[string]any{
		"tests": []SecondTest{
			{
				ID:     1,
				Thing:  "test",
				Ignore: "ignore",
				Omit:   "omit",
			},
			{
				ID:     2,
				Thing:  "test2",
				Ignore: "ignore2",
				Omit:   "",
			},
		},
	})
	for res, err := range result2.Iter() {
		t.Logf("result: %v | %v", res, err)
	}

	db.logger = &testLogger{t: t}
	_ = db.Query("SELECT * FROM ONLY test:unavailable", map[string]any{})
	_ = db.Scan(test2, "SELECT * FROM ONLY test:unavailable", map[string]any{})

	db.logger = &testLogger2{t: t}
	ctx := context.WithValue(context.Background(), "test", "test")
	_ = db.WithContext(ctx).Query("SELECT * FROM ONLY test:unavailable", map[string]any{})
	_ = db.WithContext(ctx).Scan(test2, "SELECT * FROM ONLY test:unavailable", map[string]any{})
}

type testLogger struct {
	t *testing.T
}

func (l *testLogger) Error(err error) {
	l.t.Logf("error: %v", err)
}

func (l *testLogger) Trace(ctx context.Context, t TraceType, data any) {
	l.t.Logf("trace: %v | %v", t, data)
}

type testLogger2 struct {
	testLogger
	t *testing.T
}

func (l *testLogger2) Trace(ctx context.Context, t TraceType, data any) {
	l.t.Logf("trace: %v | %v | %v", t, ctx.Value("test"), data)
}
