package errs

import (
	"errors"
	"fmt"
)

// SurgoError is a wrapper for errors returned by the Surgo package
type SurgoError struct {
	Err error
}

var (
	ErrNoConnection       = &SurgoError{fmt.Errorf("could not connect to SurrealDB")}
	ErrInvalidCredentials = &SurgoError{fmt.Errorf("invalid credentials")}
	ErrNoResult           = &SurgoError{fmt.Errorf("no result found")}
	ErrOutOfBounds        = &SurgoError{fmt.Errorf("index out of bounds")}
	ErrDatabase           = &SurgoError{fmt.Errorf("database error")}
	ErrUnmarshal          = &SurgoError{fmt.Errorf("unmarshal error")}
	ErrMarshal            = &SurgoError{fmt.Errorf("marshal error")}
)

func (e *SurgoError) With(err error) error {
	return errors.Join(e, err)
}

func (e *SurgoError) Withf(format string, args ...any) error {
	return e.With(fmt.Errorf(format, args...))
}

func (e *SurgoError) Error() string {
	return e.Err.Error()
}

func (e *SurgoError) Unwrap() error {
	return e.Err
}
