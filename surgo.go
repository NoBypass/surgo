package surgo

import (
	"errors"
	"github.com/NoBypass/surgo/v2/surrealdb"
	"github.com/NoBypass/surgo/v2/surrealdb/pkg/conn/gorilla"
	"os"
)

var fallbackTag = ""

type DB struct {
	Conn DBConn
}

type DBConn interface {
	Close() error
	Query(string, any) (any, error)
}

type Credentials struct {
	Namespace, Database, Scope, Username, Password string
}

var (
	ErrNoConnection       = errors.New("could not connect to SurrealDB")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrNoResult           = errors.New("no result found")
	ErrQuery              = errors.New("query failed")
)

// Connect connects to a SurrealDB instance and returns a DB object.
func Connect(url string, creds *Credentials) (*DB, error) {
	ws := gorilla.Create()
	db, err := surrealdb.New(url, ws)
	if err != nil {
		return nil, ErrNoConnection
	}

	_, err = db.Signin(&surrealdb.Auth{
		Namespace: creds.Namespace,
		Database:  creds.Database,
		Scope:     creds.Scope,
		Username:  creds.Username,
		Password:  creds.Password,
	})
	if err != nil {
		return nil, ErrInvalidCredentials
	}

	fallbackTag = os.Getenv("SURGO_FALLBACK_TAG")

	return &DB{db}, nil
}

// MustConnect connects to a SurrealDB instance and returns a DB object.
// If an error occurs, it panics.
func MustConnect(url string, creds *Credentials) *DB {
	db, err := Connect(url, creds)
	if err != nil {
		panic(err)
	}
	return db
}

// Close closes the connection to the SurrealDB instance.
func (db *DB) Close() error {
	return db.Conn.Close()
}

// Scan executes the query and scans the result into the given pointer struct or into a map.
// If multiple results are expected, a pointer to a slice of structs or maps can be passed.
// NOTE: Only the last result (the last query if multiple are present) is scanned into the
// given object. If any of the queries fail, the error is returned.
func (db *DB) Scan(dest any, query string, vars map[string]any) error {
	results, err := db.Query(query, vars)
	if err != nil {
		return err
	}

	for _, res := range results {
		if res.Error != nil {
			return res.Error
		}
	}

	return results[len(results)-1].Unmarshal(dest)
}
