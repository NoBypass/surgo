package surgo

import (
	"time"
)

type Result struct {
	Data     any
	Error    error
	Duration time.Duration
}

// Unmarshal scans the Result into the given pointer struct or into a map.
func (r *Result) Unmarshal(dest any) error {
	if r.Error != nil {
		return r.Error
	}

	return scan(r.Data, dest)
}

// Query executes the query and returns the results.
func (db *DB) Query(query string, vars map[string]any) ([]Result, error) {
	vars = parseVars(vars)
	res, err := db.Conn.Query(query, vars)
	if err != nil {
		return nil, err
	}

	return db.respToResult(res)
}

// MustQuery executes the query and panics if an error occurs at any point.
func (db *DB) MustQuery(query string, vars map[string]any) {
	res, err := db.Query(query, vars)
	if err != nil {
		panic(err)
	} else if len(res) == 0 {
		panic("no results")
	} else {
		for _, r := range res {
			if r.Error != nil {
				panic(r.Error)
			}
		}
	}
}
