package surgo

import (
	"fmt"
	"github.com/surrealdb/surrealdb.go"
)

func (dbm *DBModel[T]) FindOne(obj *T, options ...OptsFunc) error {
	options = append(options, Only())
	data, err := surrealdb.SmartUnmarshal[T](dbm.selectConstructor(options...))
	scanStruct(&obj, data)
	return err
}

func (dbm *DBModel[T]) Find(obj *[]T, options ...OptsFunc) error {
	data, err := surrealdb.SmartUnmarshal[[]T](dbm.selectConstructor(options...))
	scanSlice(obj, data)
	return err
}

func (dbm *DBModel[T]) selectConstructor(options ...OptsFunc) (any, error) {
	var opts Opts
	for _, option := range options {
		option(&opts)
	}

	query := fmt.Sprintf("SELECT %s%s FROM %s%s%s %s%s%s%s%s%s%s%s",
		fields(opts.fields),
		omit(opts.omit),
		only(opts.only),
		model(dbm.model, opts.model),
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

	return dbm.db.Query(query)
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
		model(dbm.model, opts.model),
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

	scanStruct(&record, data)
	return nil
}

// TODO: support for ID field (scan)

func (dbm *DBModel[T]) Delete(options ...OptsFunc) (*T, error) {
	var opts Opts
	for _, option := range options {
		option(&opts)
	}

	query := fmt.Sprintf("DELETE %s%s%s %s%s%s%s",
		only(opts.only),
		model(dbm.model, opts.model),
		id(opts.id),
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
		model(dbm.model, opts.model),
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

// TODO: use scan

func (dbr *DBRelation[From, To, Edge]) Delete(fromID, toID OptsFunc, options ...OptsFunc) (*Edge, error) {
	var fromOpts Opts
	var toOpts Opts
	fromID(&fromOpts)
	toID(&toOpts)

	options = append(
		options,
		Where(fmt.Sprintf("out=%s%s", dbr.to, id(toOpts.id))),
		overrideModel(fmt.Sprintf("%s%s->%s", dbr.from, id(fromOpts.id), dbr.edge)),
	)
	return dbr.model.Delete(options...)
}
