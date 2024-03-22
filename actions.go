package surgo

import (
	"fmt"
	"reflect"
	"strings"
	"time"
)

type Result struct {
	Data     any
	Error    error
	Duration time.Duration
	Query    Query
}

type Query struct {
	Query  string
	Params map[string]any
}

// Scan executes the query and scans the result into the given pointer struct or into a map.
// If multiple results are expected, a pointer to a slice of structs or maps can be passed.
// NOTE: Only the last result is scanned into the given object.
func (db *DB) Scan(scan any, query string, args ...any) error {
	v := reflect.ValueOf(scan)
	if v.Kind() != reflect.Ptr {
		return fmt.Errorf("scan must be a pointer")
	}

	res, err := db.Exec(query, args...)
	if err != nil {
		return err
	}

	last := res[len(res)-1]
	if last.Error != nil {
		return last.Error
	}

	return scanData(scan, last)
}

// Exec executes the query and returns the result. Parameters are supported as
// a map or simply multiple arguments. For Example:
//
//	db.Exec("SELECT * FROM table WHERE id = $id", map[string]any{"id": 1})
//
// or
//
//	db.Exec("SELECT * FROM table WHERE id = $1", 1)
func (db *DB) Exec(query string, args ...any) ([]Result, error) {
	params, err := parseParams(args)
	if err != nil {
		return nil, err
	}
	id, ok := params["$"]
	if ok {
		delete(params, "$")
		var s string
		switch id.(type) {
		case ID:
			s = id.(ID).string()
		case Range:
			s = id.(Range).string()
		}
		query = strings.Replace(query, ":$", fmt.Sprintf(":%s", s), 1)
	}
	return db.query(query, params)
}

// MustExec executes the query and panics if an error occurs at any point.
func (db *DB) MustExec(query string, args ...any) {
	res, err := db.Exec(query, args...)
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

type ID [2]any

type Range [2]ID
