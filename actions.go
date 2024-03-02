package surgo

import (
	"fmt"
	"github.com/surrealdb/surrealdb.go"
	"reflect"
	"strings"
)

type DBWrap[T any] struct {
	db    IDB
	model string
}

// Model takes a pointer to a record and a database connection. It is
// used to provide type safety in queries. The name of the given record
// is used as the database table name (using reflect).
func Model[T any](db IDB) DBWrap[T] {
	return DBWrap[T]{
		db:    db,
		model: reflect.TypeOf(new(T)).Elem().Name(),
	}
}

// TODO: support for time.Time

func (dbw *DBWrap[T]) Select(options ...SelectOptsFunc) (*T, error) {
	var opts SelectOpts
	for _, option := range options {
		option(&opts)
	}

	query := fmt.Sprintf("SELECT %s%s FROM %s%s%s %s%s%s%s%s%s%s%s",
		fields(opts.fields),
		omit(opts.omit),
		only(opts.only),
		dbw.model,
		id(opts.id),
		where(opts.where),
		group(opts.groups),
		order(opts.order),
		limit(opts.limit),
		start(opts.start),
		fetch(opts.fetchFields),
		timeout(opts.timeout),
		parallel(opts.parallel),
	)

	query = strings.TrimSpace(query)
	res, err := dbw.db.Query(query + ";")
	data, err := surrealdb.SmartUnmarshal[T](res, err)
	if err != nil {
		return nil, err
	}
	return &data, nil
}
