package surgo

import (
	"fmt"
	"github.com/surrealdb/surrealdb.go"
)

func (dbm *DBModel[T]) Select(obj *T, options ...OptsFunc) error {
	var opts Opts
	for _, option := range options {
		option(&opts)
	}

	query := fmt.Sprintf("SELECT %s%s FROM %s%s%s %s%s%s%s%s%s%s%s",
		fields(opts.fields),
		omit(opts.omit),
		only(opts.only),
		dbm.model,
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

	res, err := dbm.db.Query(query)
	data, err := surrealdb.SmartUnmarshal[T](res, err)
	scan(&obj, data)
	return err
}

// TODO: support for ID field
// TODO: support for slices of records

func (dbm *DBModel[T]) Create(record *T, options ...OptsFunc) error {
	var opts Opts
	for _, option := range options {
		option(&opts)
	}

	query := fmt.Sprintf("CREATE %s%s%s%s%s%s%s",
		only(opts.only),
		dbm.model,
		id(opts.id),
		content(record),
		returns(opts.returns),
		timeout(opts.timeout),
		parallel(opts.parallel),
	)

	res, err := dbm.db.Query(query)
	data, err := surrealdb.SmartUnmarshal[T](res, err)
	if err != nil {
		return err
	}

	scan(&record, data)
	return nil
}

// TODO: support for ID field (scan)

func (dbm *DBModel[T]) Delete(ID string, options ...OptsFunc) (*T, error) {
	var opts Opts
	for _, option := range options {
		option(&opts)
	}

	query := fmt.Sprintf("DELETE %s%s%s %s%s%s%s",
		only(opts.only),
		dbm.model,
		id(ID),
		where(opts.where),
		returns(opts.returns),
		timeout(opts.timeout),
		parallel(opts.parallel),
	)

	res, err := dbm.db.Query(query)
	data, err := surrealdb.SmartUnmarshal[T](res, err)
	if err != nil {
		return nil, err
	}
	return &data, nil
}

// TODO support for set and merge

func (dbm *DBModel[T]) Update(record *T, options ...OptsFunc) error {
	var opts Opts
	for _, option := range options {
		option(&opts)
	}

	query := fmt.Sprintf("UPDATE %s%s%s%s%s%s%s",
		only(opts.only),
		dbm.model,
		id(opts.id),
		content(record),
		returns(opts.returns),
		timeout(opts.timeout),
		parallel(opts.parallel),
	)

	_, err := dbm.db.Query(query)
	return err
}

func (dbr *DBRelation[From, To, Edge]) Create(edge *Edge, fromID, toID OptsFunc, options ...OptsFunc) error {
	var fromOpts Opts
	var toOpts Opts
	fromID(&fromOpts)
	toID(&toOpts)

	var opts Opts
	for _, option := range options {
		option(&opts)
	}

	contentStr := ""
	if edge != nil {
		contentStr = content(edge)
	}

	query := fmt.Sprintf("RELATE %s%s%s->%s->%s%s%s %s%s%s",
		only(opts.only),
		dbr.from,
		id(fromOpts.id),
		dbr.edge,
		dbr.to,
		id(toOpts.id),
		contentStr,
		returns(opts.returns),
		timeout(opts.timeout),
		parallel(opts.parallel),
	)

	_, err := dbr.db.Query(query)
	return err
}

func (dbr *DBRelation[From, To, Edge]) Delete(fromID, toID OptsFunc, options ...OptsFunc) (*Edge, error) {
	var fromOpts Opts
	var toOpts Opts
	fromID(&fromOpts)
	toID(&toOpts)

	var opts Opts
	for _, option := range options {
		option(&opts)
	}

	whereStr := fmt.Sprintf("out=%s%s%s", dbr.to, id(toOpts.id), func() string {
		if opts.where == "" {
			return ""
		}
		return " AND " + opts.where
	}())

	query := fmt.Sprintf("DELETE %s%s%s->%s %s%s%s%s",
		only(opts.only),
		dbr.from,
		id(fromOpts.id),
		dbr.edge,
		where(whereStr),
		returns(opts.returns),
		timeout(opts.timeout),
		parallel(opts.parallel),
	)

	res, err := dbr.db.Query(query)
	data, err := surrealdb.SmartUnmarshal[Edge](res, err)
	return &data, err
}
