package surgo

import "fmt"

type (
	ErrNoConnection       error
	ErrInvalidCredentials error
	ErrNoResult           error
	ErrQuery              error
)

func newErrNoConnection(underlying error) error {
	return ErrNoConnection(fmt.Errorf("could not connect to SurrealDB: %w", underlying))
}

func newErrInvalidCredentials(underlying error) error {
	return ErrInvalidCredentials(fmt.Errorf("invalid credentials: %w", underlying))
}

func newErrNoResult(underlying error) error {
	return ErrNoResult(fmt.Errorf("no result found: %w", underlying))
}

func newErrQuery(underlying error) error {
	return ErrQuery(fmt.Errorf("query failed: %w", underlying))
}
