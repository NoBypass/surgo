package surgo

import (
	"context"
	"github.com/NoBypass/surgo/v2/errs"
	"github.com/NoBypass/surgo/v2/marshal"
	"github.com/NoBypass/surgo/v2/rpc"
	"time"
)

type DB struct {
	Conn      *rpc.WebsocketConn
	Marshaler marshal.Marshaler
	timeout   time.Duration
	logger    Logger

	// ctx is only populated if WithContext is used.
	ctx context.Context
}

// Credentials contains the necessary information to connect to a SurrealDB instance.
// Namespace, Database and Scope are optional, if not provided the signin will happen
// on the Root, Namespace or Database level respectively.
type Credentials struct {
	Namespace string `json:"NS,omitempty"`
	Database  string `json:"DB,omitempty"`
	Scope     string `json:"SC,omitempty"`
	Username  string `json:"user,omitempty"`
	Password  string `json:"pass,omitempty"`
}

// Connect connects to a SurrealDB instance and returns a DB object.
func Connect(url string, creds *Credentials, opts ...Option) (*DB, error) {
	db := &DB{
		Marshaler: marshal.Marshaler(""),
		timeout:   10 * time.Second,
		logger:    &defaultLogger{},
	}

	for _, opt := range opts {
		opt(db)
	}

	c, err := rpc.NewWebsocketConn(url, db.logger)
	if err != nil {
		return nil, errs.ErrNoConnection.With(err)
	}
	db.Conn = c

	ctx, cancel := context.WithTimeout(context.Background(), db.timeout)
	defer cancel()
	_, err = c.Send(ctx, "signin", []any{creds})
	if err != nil {
		return nil, errs.ErrInvalidCredentials.With(err)
	}

	return db, nil
}

// Close closes the connection to the SurrealDB instance.
func (db *DB) Close() error {
	return db.Conn.Close()
}

func (db DB) WithContext(ctx context.Context) *DB {
	return &DB{
		Conn:      db.Conn,
		Marshaler: db.Marshaler,
		timeout:   db.timeout,
		logger:    db.logger,
		ctx:       ctx,
	}
}
