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
	params, err := db.parseParams(args)
	if err != nil {
		return nil, err
	}
	return db.query(query, params)
}

// MustExec executes the query and panics if an error occurs at any point.
// Parameters are supported as a map or simply multiple arguments. For Example:
//
//	db.MustExec("SELECT * FROM table WHERE id = $id", map[string]any{"id": 1})
//
// or
//
//	db.MustExec("SELECT * FROM table WHERE id = $1", 1)
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

func (db *DB) query(query string, params map[string]any) ([]Result, error) {
	if !strings.HasSuffix(query, ";") {
		query = query + ";"
	}
	resp, err := db.db.Query(query, params)
	if err != nil {
		return nil, err
	}

	respSlice := resp.([]any)
	resSlice := make([]Result, len(respSlice))
	for i, s := range respSlice {
		m := s.(map[string]any)
		d, err := time.ParseDuration(m["time"].(string))
		if err != nil {
			return nil, err
		}

		resSlice[i] = Result{
			Data: m["result"],
			Error: func() error {
				e, ok := m["error"]
				if m["result"] == nil || !ok {
					return nil
				}
				return fmt.Errorf(e.(string))
			}(),
			Duration: d,
		}
	}
	return resSlice, nil
}
