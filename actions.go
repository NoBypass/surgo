package surgo

import (
	"fmt"
	"github.com/surrealdb/surrealdb.go"
	"reflect"
	"strings"
)

// TODO: support for time.Time

func (dbw *DBWrap[T]) Select(obj *T, options ...OptsFunc) error {
	var opts Opts
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
	scan(&obj, data)
	return err
}

// TODO: support for ID field
// TODO: support for slices of records

func (dbw *DBWrap[T]) Create(content *T, options ...OptsFunc) error {
	var opts Opts
	for _, option := range options {
		option(&opts)
	}

	contentStr := " CONTENT {"
	v := reflect.ValueOf(content).Elem()
	t := v.Type()

	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		fieldType := t.Field(i)
		surrealTag := fieldType.Tag.Get("surreal")
		if surrealTag == "" {
			surrealTag = fieldType.Name
		}
		switch field.Kind() {
		case reflect.String:
			contentStr += fmt.Sprintf(`%s:"%v",`, surrealTag, field.Interface())
		default:
			contentStr += fmt.Sprintf("%s:%v,", surrealTag, field.Interface())
		}
	}

	contentStr = contentStr[:len(contentStr)-1] + "} "

	query := fmt.Sprintf("CREATE %s%s%s%s%s%s%s",
		only(opts.only),
		dbw.model,
		id(opts.id),
		contentStr,
		returns(opts.returns),
		timeout(opts.timeout),
		parallel(opts.parallel),
	)

	query = strings.TrimSpace(query)
	res, err := dbw.db.Query(query + ";")
	data, err := surrealdb.SmartUnmarshal[T](res, err)
	if err != nil {
		return err
	}

	scan(&content, data)
	return nil
}

func (dbw *DBWrap[T]) Delete(ID string, options ...OptsFunc) (*T, error) {
	var opts Opts
	for _, option := range options {
		option(&opts)
	}

	query := fmt.Sprintf("DELETE %s%s%s %s%s%s%s",
		only(opts.only),
		dbw.model,
		id(ID),
		where(opts.where),
		returns(opts.returns),
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

// TODO support for set and merge

func (dbw *DBWrap[T]) Update(content *T, options ...OptsFunc) error {
	var opts Opts
	for _, option := range options {
		option(&opts)
	}

	contentStr := " CONTENT {"
	v := reflect.ValueOf(content).Elem()
	t := v.Type()

	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		fieldType := t.Field(i)
		surrealTag := fieldType.Tag.Get("surreal")
		if surrealTag == "" {
			surrealTag = fieldType.Name
		}
		switch field.Kind() {
		case reflect.String:
			contentStr += fmt.Sprintf(`%s:"%v",`, surrealTag, field.Interface())
		default:
			contentStr += fmt.Sprintf("%s:%v,", surrealTag, field.Interface())
		}
	}

	contentStr = contentStr[:len(contentStr)-1] + "} "

	query := fmt.Sprintf("UPDATE %s%s%s%s%s%s%s",
		only(opts.only),
		dbw.model,
		id(opts.id),
		contentStr,
		returns(opts.returns),
		timeout(opts.timeout),
		parallel(opts.parallel),
	)

	query = strings.TrimSpace(query)
	_, err := dbw.db.Query(query + ";")
	if err != nil {
		return err
	}
	return nil
}
