package error

import "fmt"

type Type int

const (
	NotFound Type = iota
	NotAuthorized
)

type Error struct {
	ErrorType Type
	message   string
	wrapped   error
}

func (e *Error) Error() string {
	if e.wrapped != nil {
		return fmt.Sprintf("ErrorType: %d, message: %s\n%s", e.ErrorType, e.message, e.wrapped.Error())
	}
	return fmt.Sprintf("ErrorType: %d, message: %s", e.ErrorType, e.message)
}

func NewError(errorType Type, message string, original error) *Error {
	return &Error{
		ErrorType: errorType,
		message:   message,
		wrapped:   original,
	}
}
