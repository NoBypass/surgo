package surgo

import (
	"context"
	"github.com/NoBypass/surgo/v2/errs"
	"github.com/NoBypass/surgo/v2/marshal"
	"log"
	"time"
)

/* ---------- Tracer & Logger ---------- */

const (
	TraceQuery TraceType = iota
	TraceVars
	TraceResponse
	TraceEnd
)

// TraceType is used to specify the type of trace.
type TraceType int

// Logger can be implemented for logging and tracing purposes.
type Logger interface {
	// Error logs errors which occur while reading the SurrealDB
	// response as these errors cannot be returned normally.
	Error(err error)
	// Trace will only be called if a context was used to send the query.
	// If you want to do any tracing or detailed logging, you can use this method.
	Trace(ctx context.Context, t TraceType, data any)
}

type defaultLogger struct{}

func (l *defaultLogger) Error(err error) {
	log.Println(err)
}

func (l *defaultLogger) Trace(ctx context.Context, t TraceType, data any) {
}

type silentLogger struct {
	defaultLogger
}

func (l *silentLogger) Error(err error) {
}

/* ---------- Options ---------- */

type Option func(*DB)

// WithDefaultTimeout sets the default timeout for queries
func WithDefaultTimeout(timeout time.Duration) Option {
	return func(db *DB) {
		db.timeout = timeout
	}
}

// WithLogger sets the error logger for the DB
func WithLogger(l Logger) Option {
	if l == nil {
		return WithDisableLogging()
	}
	return func(db *DB) {
		db.logger = l
	}
}

// WithDisableLogging disables logging for the DB
func WithDisableLogging() Option {
	return func(db *DB) {
		db.logger = &silentLogger{}
	}
}

// WithFallbackTag sets the fallback tag for the Marshaler
func WithFallbackTag(tag string) Option {
	return func(db *DB) {
		db.Marshaler = marshal.Marshaler(tag)
	}
}

/* ---------- Misc ---------- */

func safeContext(ctx context.Context) context.Context {
	if ctx == nil {
		return context.Background()
	}
	return ctx
}

func resultsToQuery(res []any) ([]Query, error) {
	var results []Query
	for _, r := range res {
		q := r.(map[string]any)
		r, ok := q["result"]
		if !ok {
			return nil, errs.ErrMarshal.Withf("invalid response, missing result")
		}

		s, statusOk := q["status"]
		if !statusOk {
			return nil, errs.ErrMarshal.Withf("invalid response, missing status")
		}

		var ee error
		if str, ok := r.(string); s == "ERR" && ok {
			ee = errs.ErrDatabase.Withf(str)
			r = nil
		} else if str == "ERR" {
			return nil, errs.ErrMarshal.Withf("invalid response, unexpected error format")
		} else if r == nil {
			ee = errs.ErrNoResult
		}

		results = append(results, Query{
			Result: r,
			Error:  ee,
		})
	}
	return results, nil
}
