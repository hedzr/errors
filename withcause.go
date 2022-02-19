package errors

import "fmt"

type causer interface {
	Cause() error
}

// Cause1 returns the underlying cause of the error, if possible.
// Cause1 unwraps just one level of the inner wrapped error.
//
// An error value has a cause if it implements the following
// interface:
//
//     type causer interface {
//            Cause() error
//     }
//
// If the error does not implement Cause, the original error will
// be returned. If the error is nil, nil will be returned without further
// investigation.
func Cause1(err error) error {
	if e, ok := err.(causer); ok {
		return e.Cause()
	}
	return err
}

// Cause returns the underlying cause of the error recursively,
// if possible.
// An error value has a cause if it implements the following
// interface:
//
//     type causer interface {
//            Cause() error
//     }
//
// If the error does not implement Cause, the original error will
// be returned. If the error is nil, nil will be returned without further
// investigation.
func Cause(err error) error {
	for err != nil {
		if cause, ok := err.(causer); ok {
			err = cause.Cause()
		} else {
			break
		}
	}
	return err
}

type withCause struct {
	causer error
	msg    string
}

// WithCause is synonym of Wrap
func WithCause(cause error, message string, args ...interface{}) error {
	if len(args) > 0 {
		message = fmt.Sprintf(message, args...)
	}
	return &withCause{cause, message}
}

func (w *withCause) Error() string {
	if w.causer != nil {
		return w.msg + ": " + w.causer.Error()
	}
	return w.msg
}

// Attach appends errs
func (w *withCause) Attach(errs ...error) {
	for _, err := range errs {
		if err != nil {
			w.causer = err
		}
	}
}

// Cause returns the underlying cause of the error recursively,
// if possible.
func (w *withCause) Cause() error {
	return w.causer
}

// Unwrap returns the result of calling the Unwrap method on err, if
// `err`'s type contains an Unwrap method returning error.
// Otherwise, Unwrap returns nil.
func (w *withCause) Unwrap() error {
	return w.causer
}

// As finds the first error in `err`'s chain that matches target, and if so, sets
// target to that error value and returns true.
func (w *withCause) As(target interface{}) bool {
	return As(w.causer, target)
}

// Is reports whether any error in `err`'s chain matches target.
func (w *withCause) Is(target error) bool {
	return w.causer == target || Is(w.causer, target)
}

// TypeIs reports whether any error in `err`'s chain matches target.
func (w *withCause) TypeIs(target error) bool {
	return w.causer == target || TypeIs(w.causer, target)
}
