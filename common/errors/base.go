package errors

import (
	std_errors "errors"
	"fmt"
)

// New returns a new error with a message and a stack.
func New(msg string) error {
	err := std_errors.New(msg)
	err = withStack(err, 2)
	return err
}

// Errorf returns a new error with a formatted message and a stack.
func Errorf(format string, args ...interface{}) error {
	err := fmt.Errorf(format, args...)
	err = withStack(err, 2)
	return err
}
