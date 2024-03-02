package surgo

import (
	"context"
	"fmt"
	"github.com/surrealdb/surrealdb.go"
)

type Option [2]string

type DB struct {
	db  *surrealdb.DB
	ctx context.Context
}

type IDB interface {
	Query(string) (interface{}, error)
}

func New(ctx context.Context, url string, options ...Option) (*DB, error) {
	db, err := surrealdb.New(fmt.Sprintf("wss://%s/rpc", url))
	if err != nil {
		db, err = surrealdb.New(fmt.Sprintf("ws://%s/rpc", url))
		if err != nil {
			return nil, err
		}
	}

	confMap := make(map[string]string)
	for _, option := range options {
		confMap[option[0]] = option[1]
	}

	if _, err = db.Signin(confMap); err != nil {
		return nil, err
	}

	ns, ok := confMap["NS"]
	if !ok {
		ns = ""
	}
	dbName, ok := confMap["DB"]
	if !ok {
		dbName = ""
	}
	_, err = db.Use(ns, dbName)
	if err != nil {
		return nil, err
	}

	return &DB{
		db:  db,
		ctx: ctx,
	}, nil
}

func User(username string) Option {
	return Option{"user", username}
}

func Pass(password string) Option {
	return Option{"pass", password}
}

func Namespace(namespace string) Option {
	return Option{"NS", namespace}
}

func Database(database string) Option {
	return Option{"DB", database}
}
