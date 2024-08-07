package constants

import "errors"

var (
	InvalidResponse = errors.New("invalid SurrealDB response")
	ErrQuery        = errors.New("error occurred processing the SurrealDB query")
	ErrNoRow        = errors.New("error no row")
)
