package code

import (
	"fmt"
	"net/http"

	"github.com/pkg/errors"
)

type privateError struct {
	code Code
	err  error
}

func (e privateError) Error() string {
	return fmt.Sprintf("Code: %s, Message: %+v", e.code, e.err)
}

func (e privateError) Unwrap() error {
	return e.err
}

// Errorf formats error with code
func Errorf(c Code, format string, args ...interface{}) error {
	if c == OK {
		return nil
	}
	return privateError{
		code: c,
		err:  errors.WithStack(errors.New(fmt.Sprintf(format, args...))),
	}
}

// Error format error with code
func Error(c Code, format string) error {
	if c == OK {
		return nil
	}
	return privateError{
		code: c,
		err:  errors.WithStack(errors.New(format)),
	}
}

// GetCode gets code
func GetCode(err error) Code {
	if err == nil {
		return OK
	}
	if e, ok := err.(privateError); ok {
		return e.code
	}
	return Unknown
}

// GetError gets err
func GetError(err error) error {
	if err == nil {
		return nil
	}
	if e, ok := err.(privateError); ok {
		return e.err
	}
	return Error(Unknown, "failed to get error")
}

// GetStatusCode gets status code
func GetStatusCode(err error) int {
	switch GetCode(err) {
	case OK:
		return http.StatusOK

	case InvalidQuery:
		return http.StatusBadRequest

	default:
		return http.StatusInternalServerError
	}
}
