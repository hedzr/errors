package old

import "fmt"

// func As(err error, target interface{}) bool
// func Is(err, target error) bool
// func New(text string) error
// func Unwrap(err error) error

// Unwrap returns the result of calling the Unwrap method on err, if
// `err`'s type contains an Unwrap method returning error.
// Otherwise, Unwrap returns nil.
func Unwrap(err error) error {
	u, ok := err.(interface {
		Unwrap() error
	})
	if !ok {
		return nil
	}
	return u.Unwrap()
}

// Wrap returns an error annotating err with a Stack trace
// at the point Wrap is called, and the supplied message.
// If err is nil, Wrap returns nil.
func Wrap(err error, message string, args ...interface{}) *WithStackInfo {
	if err == nil {
		return nil
	}
	if len(args) > 0 {
		message = fmt.Sprintf(message, args...)
	}
	err = &withCause{
		causer: err,
		msg:    message,
	}
	return &WithStackInfo{
		err,
		callers(1),
	}
}
