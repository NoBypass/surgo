package surgo

import (
	"context"
	"github.com/NoBypass/surgo/v2/errs"
	"iter"
)

const scanCtxKey = "surgo:is_scan"

type (
	Result struct {
		Error   error
		Queries []Query
	}
	Query struct {
		Result any
		Error  error
	}
)

// Query executes the query and returns the results. The error is
// only not nil if the whole call failed. If a query fails, the error
// is stored in the result struct.
func (db *DB) Query(query string, vars map[string]any) (result Result) {
	ctx, cancel := context.WithTimeout(safeContext(db.ctx), db.timeout)
	defer cancel()

	db.logger.Trace(ctx, TraceQuery, query)
	vars = db.Marshaler.Marshal(vars)
	db.logger.Trace(ctx, TraceVars, vars)
	if ctx.Value(scanCtxKey) == nil {
		defer func() {
			db.logger.Trace(ctx, TraceEnd, result)
		}()
	}

	res, err := db.Conn.Send(ctx, "query", []any{query, vars})
	if err != nil {
		return Result{Error: err}
	}

	db.logger.Trace(ctx, TraceResponse, res)

	queries, err := resultsToQuery(res.([]any))
	return Result{
		Error:   err,
		Queries: queries,
	}
}

// Scan executes the query and scans the result into the given pointer struct or into a map.
// If multiple results are expected, a pointer to a slice of structs or maps can be passed.
// NOTE: Only the last result (the last query if multiple are present) is scanned into the
// given object. If any of the queries fail, the error is returned.
func (db DB) Scan(dest any, query string, vars map[string]any) error {
	db.ctx = context.WithValue(safeContext(db.ctx), scanCtxKey, true)
	defer func() {
		db.logger.Trace(db.ctx, TraceEnd, nil)
	}()

	result := db.Query(query, vars)
	if result.Error != nil {
		return result.Error
	}

	queryResult, err := result.Last()
	if err != nil {
		return err
	}

	return db.Marshaler.Unmarshal(queryResult, dest)
}

// First returns the result of the first query. If the query failed, the error is returned.
func (r *Result) First() (any, error) {
	if r.Error != nil {
		return nil, r.Error
	} else if len(r.Queries) == 0 {
		return nil, errs.ErrNoResult
	}
	return r.Queries[0].Result, r.Queries[0].Error
}

// Last returns the result of the last query. If the query failed, the error is returned.
func (r *Result) Last() (any, error) {
	if r.Error != nil {
		return nil, r.Error
	} else if len(r.Queries) == 0 {
		return nil, errs.ErrNoResult
	}
	return r.Queries[len(r.Queries)-1].Result, r.Queries[len(r.Queries)-1].Error
}

// At returns the result of the query at the given index. If the query failed, the error is returned.
func (r *Result) At(i int) (any, error) {
	if r.Error != nil {
		return nil, r.Error
	} else if i < 0 || i >= len(r.Queries) {
		return nil, errs.ErrOutOfBounds
	}
	return r.Queries[i].Result, r.Queries[i].Error
}

// Iter returns an iterator over the queries. If the query failed, the error is returned.
func (r *Result) Iter() iter.Seq2[any, error] {
	return func(yield func(any, error) bool) {
		if r.Error != nil {
			yield(nil, r.Error)
			return
		}

		for _, q := range r.Queries {
			if !yield(q.Result, q.Error) {
				return
			}
		}
	}
}
