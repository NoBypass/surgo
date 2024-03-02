package surgo

import (
	"fmt"
	"reflect"
	"strings"
	"time"
)

type SelectOptsFunc func(*SelectOpts)
type SelectOpts struct {
	timeout     time.Duration
	fetchFields []string
	order       []string
	groups      []string
	fields      []string
	omit        []string
	where       string
	id          string
	parallel    bool
	only        bool
	limit       int
	start       int
}

// Fields is an option for `Select` to specify the fields to return.
// You can also use the `AS` keyword to use an alias for the fields.
func Fields(fields ...string) SelectOptsFunc {
	return func(o *SelectOpts) {
		o.fields = fields
	}
}

// Only is an option for `Select` to specifies to only select a single record.
func Only() SelectOptsFunc {
	return func(o *SelectOpts) {
		o.only = true
	}
}

// Omit is an option for `Select` to specify the fields to omit from the result.
func Omit(fields ...string) SelectOptsFunc {
	return func(o *SelectOpts) {
		o.omit = fields
	}
}

// Where is an option for `Select` to specify the condition to filter the records.
func Where(condition string) SelectOptsFunc {
	return func(o *SelectOpts) {
		o.where = condition
	}
}

// GroupBy is an option for `Select` to specify the fields to group the records by.
func GroupBy(fields ...string) SelectOptsFunc {
	return func(o *SelectOpts) {
		o.groups = fields
	}
}

type OrderOptsFunc func(*OrderOpts)
type OrderOpts struct {
	orders []string
}

// OrderBy is an option for `Select` to specify the fields to order the records by.
// It takes a list of `OrderOptsFunc` to specify the order of the fields. You can
// use the `Asc` and `Desc` functions to specify the order.
func OrderBy(fields ...OrderOptsFunc) SelectOptsFunc {
	var opts OrderOpts
	for _, option := range fields {
		option(&opts)
	}
	return func(o *SelectOpts) {
		o.order = opts.orders
	}
}

// Asc is an option for `OrderBy` to specify the ascending order of the field.
func Asc(field string) OrderOptsFunc {
	return func(o *OrderOpts) {
		o.orders = append(o.orders, fmt.Sprintf("%s ASC", field))
	}
}

// Desc is an option for `OrderBy` to specify the descending order of the field.
func Desc(field string) OrderOptsFunc {
	return func(o *OrderOpts) {
		o.orders = append(o.orders, fmt.Sprintf("%s DESC", field))
	}
}

// Start is an option for `Select` to specify the starting index of the records to return.
func Start(start int) SelectOptsFunc {
	return func(o *SelectOpts) {
		o.start = start
	}
}

// Limit is an option for `Select` to specify the maximum number of records to return.
func Limit(limit int) SelectOptsFunc {
	return func(o *SelectOpts) {
		o.limit = limit
	}
}

// Fetch is an option for `Select` to specify the fields to fetch from the result.
func Fetch(fields ...string) SelectOptsFunc {
	return func(o *SelectOpts) {
		o.fetchFields = fields
	}
}

// Timeout is an option for `Select` to specify the maximum time to wait for the query to complete.
func Timeout(d time.Duration) SelectOptsFunc {
	return func(o *SelectOpts) {
		o.timeout = d
	}
}

type SliceOrString[T any] interface {
	~[]T | string
}

// TODO: support for slices and singular values in ranged ID's

// ID is an option for `Select` to specify the ID of the record to select.
// The id can be a single string or if you want to use ranges, you can pass
// a slice with a length of 2, where the first element is the start and the
// second element is the end of the range. An example of a ranged id is:
//
//	`ID([2]string{Bar{123}, Bar{456}})` -> `SELECT * FROM Foo:[123]..[456];`
//
// (The `Bar` type is just an example, you can use any type you want. It is
// strongly recommended to use a struct for this purpose, so you can make
// use of type safety.)
func ID[T ~string | ~[2]any](id T) SelectOptsFunc {
	idStrs := make([]string, 0, 2)
	if reflect.TypeOf(id).Kind() == reflect.Array {
		for _, v := range reflect.ValueOf(id).Interface().([2]interface{}) {
			t := reflect.ValueOf(v)
			fields := t.NumField()
			if fields == 0 {
				idStrs = append(idStrs, "")
				continue
			}
			curr := "["
			for j := range fields {
				curr += fmt.Sprintf("%v", t.Field(j).Interface())
				if j != fields-1 {
					curr += ", "
				}
			}
			idStrs = append(idStrs, curr+"]")
		}
	} else {
		idStrs = append(idStrs, fmt.Sprintf("%s", id))
	}
	return func(o *SelectOpts) {
		o.id = strings.Join(idStrs, "..")
	}
}

// Parallel is an option for `Select` to specify to run the query in parallel.
func Parallel() SelectOptsFunc {
	return func(o *SelectOpts) {
		o.parallel = true
	}
}
