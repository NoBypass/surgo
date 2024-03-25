package surgo

import (
	"fmt"
	"github.com/surrealdb/surrealdb.go"
)

type Option struct {
	key string
	val any
}

type DB struct {
	db QueryAgent
}

type QueryAgent interface {
	Query(sql string, vars any) (any, error)
	Close()
}

// Connect connects to a SurrealDB instance and returns a DB object.
func Connect(url string, options ...Option) (*DB, error) {
	db, err := surrealdb.New(fmt.Sprintf("wss://%s/rpc", url))
	if err != nil {
		db, err = surrealdb.New(fmt.Sprintf("ws://%s/rpc", url))
		if err != nil {
			return nil, err
		}
	}

	confMap := make(map[string]any)
	for _, option := range options {
		confMap[option.key] = option.val
	}

	if _, ok := confMap["agent"]; !ok {
		confMap["agent"] = DefaultAgent
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
	_, err = db.Use(ns.(string), dbName.(string))
	if err != nil {
		return nil, err
	}

	return &DB{
		db: confMap["agent"].(AgentFunc)(db),
	}, nil
}

// MustConnect connects to a SurrealDB instance and returns a DB object.
// If an error occurs, it panics.
func MustConnect(url string, options ...Option) *DB {
	db, err := Connect(url, options...)
	if err != nil {
		panic(err)
	}
	return db
}

func (db *DB) Close() {
	db.db.Close()
}

func User(username string) Option {
	return Option{"user", username}
}

func Password(password string) Option {
	return Option{"pass", password}
}

func Namespace(namespace string) Option {
	return Option{"NS", namespace}
}

func Database(database string) Option {
	return Option{"DB", database}
}

var DefaultAgent = func(db *surrealdb.DB) QueryAgent {
	return db
}

type AgentFunc func(db *surrealdb.DB) QueryAgent

func CustomAgent(agentFunc AgentFunc) Option {
	return Option{"agent", agentFunc}
}
