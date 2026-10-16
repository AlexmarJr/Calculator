package models

// ErrorKind describes a business failure without depending on a transport.
type ErrorKind int

const (
	InvalidInput ErrorKind = iota + 1
	UndefinedResult
)

type Error struct {
	Kind    ErrorKind
	Message string
}

func (e *Error) Error() string {
	return e.Message
}

func Invalid(message string) error {
	return &Error{
		Kind:    InvalidInput,
		Message: message,
	}
}

func Undefined(message string) error {
	return &Error{
		Kind:    UndefinedResult,
		Message: message,
	}
}
