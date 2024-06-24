package error

import "fmt"

type Type int

const (
	NotFound Type = iota
	Forbidden
)

func (t Type) String() string {
	switch t {
	case NotFound:
		return "Not Found"
	case Forbidden:
		return "Forbidden"
	default:
		return "Unknown"
	}
}

type Error struct {
	ErrorType Type
	message   string
	wrapped   error
}

func (e *Error) Error() string {
	if e.wrapped != nil {
		return fmt.Sprintf("ErrorType: %s, message: %s\n%s", e.ErrorType, e.message, e.wrapped.Error())
	}
	return fmt.Sprintf("ErrorType: %s, message: %s", e.ErrorType, e.message)
}

func NewError(errorType Type, message string, original error) *Error {
	return &Error{
		ErrorType: errorType,
		message:   message,
		wrapped:   original,
	}
}
