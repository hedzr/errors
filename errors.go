// Copyright © 2020 Hedzr Yeh.

package errors

import (
	"fmt"
	"github.com/pkg/errors"
	"io"
	"reflect"
)

func New(message string, args ...interface{}) error {
	if len(args) > 0 {
		message = fmt.Sprintf(message, args...)
	}
	return errors.New(message)
}

type causer interface {
	Cause() error
}

// Cause1 returns the underlying cause of the error, if possible.
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

func WithCause(cause error, msg string) error {
	return &withCause{cause, msg}
}

func (w *withCause) Error() string {
	return w.msg + ": " + w.causer.Error()
}

func (w *withCause) Cause() error {
	return w.causer
}

func (w *withCause) Unwrap() error {
	return w.causer
}

func (w *withCause) Is(target error) bool {
	if target == nil {
		return w.causer == target
	}

	isComparable := reflect.TypeOf(target).Comparable()
	for {
		if isComparable && w.causer == target {
			return true
		}
		if x, ok := w.causer.(interface{ Is(error) bool }); ok && x.Is(target) {
			return true
		}
		// TODO: consider supporing target.Is(err). This would allow
		// user-definable predicates, but also may allow for coping with sloppy
		// APIs, thereby making it easier to get away with them.
		if err := Unwrap(w.causer); err == nil {
			return false
		}
	}
}

type withCauses struct {
	causers []error
	msg     string
	*stack
}

func (w *withCauses) Error() error {
	if len(w.causers) == 0 {
		return nil
	}
	return w.wrap(w.causers...)
}

func (w *withCauses) wrap(errs ...error) error {
	return &causes{
		Causers: errs,
		stack:   w.stack,
	}
}

func (w *withCauses) Cause() error {
	if len(w.causers) == 0 {
		return nil
	}
	return w.causers[0]
}

func (w *withCauses) Causes() []error {
	if len(w.causers) == 0 {
		return nil
	}
	return w.causers
}

func (w *withCauses) Unwrap() error {
	return w.Cause()
}

func (w *withCauses) Attach(errs ...error) {
	w.causers = append(w.causers, errs...)
	w.stack = callers()
}

func (w *withCauses) IsEmpty() bool {
	return len(w.causers) == 0
}

func (w *withCauses) Is(target error) bool {
	if target == nil {
		for _, e := range w.causers {
			if e == target {
				return true
			}
		}
		return false
	}

	isComparable := reflect.TypeOf(target).Comparable()
	for {
		if isComparable {
			for _, e := range w.causers {
				if e == target {
					return true
				}
			}
			return false
		}

		for _, e := range w.causers {
			if x, ok := e.(interface{ Is(error) bool }); ok && x.Is(target) {
				return true
			}
			if err := Unwrap(e); err == nil {
				return false
			}
		}
		return false
	}
}

// Wrap returns an error annotating err with a stack trace
// at the point Wrap is called, and the supplied message.
// If err is nil, Wrap returns nil.
func Wrap(err error, message string, args ...interface{}) error {
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
	return &withStack{
		err,
		callers(),
	}
}

type withStack struct {
	error
	*stack
}

// WithStack annotates err with a stack trace at the point WithStack was called.
// If err is nil, WithStack returns nil.
func WithStack(cause error) error {
	if cause == nil {
		return nil
	}
	return &withStack{cause, callers()}
}

func (w *withStack) Cause() error {
	return w.error
}

func (w *withStack) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		if s.Flag('+') {
			fmt.Fprintf(s, "%+v", w.Cause())
			w.stack.Format(s, verb)
			return
		}
		fallthrough
	case 's':
		io.WriteString(s, w.Error())
	case 'q':
		fmt.Fprintf(s, "%q", w.Error())
	}
}

func (w *withStack) Is(target error) bool {
	if x, ok := w.error.(interface{ Is(error) bool }); ok && x.Is(target) {
		return true
	}
	return false
}

func (w *withStack) Unwrap() error {
	if x, ok := w.error.(interface{ Unwrap() error }); ok {
		return x.Unwrap()
	}
	return nil
}

func (w *withStack) Attach(errs ...error) {
	if x, ok := w.error.(interface{ Attach(errs ...error) }); ok {
		x.Attach(errs...)
	}
}

func (w *withStack) IsEmpty() bool {
	if x, ok := w.error.(interface{ IsEmpty() bool }); ok {
		return x.IsEmpty()
	}
	return false
}

// func As(err error, target interface{}) bool
// func Is(err, target error) bool
// func New(text string) error
// func Unwrap(err error) error

// Unwrap returns the result of calling the Unwrap method on err, if err's
// type contains an Unwrap method returning error.
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

// Is reports whether any error in err's chain matches target.
//
// The chain consists of err itself followed by the sequence of errors obtained by
// repeatedly calling Unwrap.
//
// An error is considered to match a target if it is equal to that target or if
// it implements a method Is(error) bool such that Is(target) returns true.
func Is(err, target error) bool {
	if target == nil {
		return err == target
	}

	isComparable := reflect.TypeOf(target).Comparable()
	for {
		if isComparable && err == target {
			return true
		}
		if x, ok := err.(interface{ Is(error) bool }); ok && x.Is(target) {
			return true
		}
		// TODO: consider supporing target.Is(err). This would allow
		// user-definable predicates, but also may allow for coping with sloppy
		// APIs, thereby making it easier to get away with them.
		if err = Unwrap(err); err == nil {
			return false
		}
	}
}

// As finds the first error in err's chain that matches target, and if so, sets
// target to that error value and returns true.
//
// The chain consists of err itself followed by the sequence of errors obtained by
// repeatedly calling Unwrap.
//
// An error matches target if the error's concrete value is assignable to the value
// pointed to by target, or if the error has a method As(interface{}) bool such that
// As(target) returns true. In the latter case, the As method is responsible for
// setting target.
//
// As will panic if target is not a non-nil pointer to either a type that implements
// error, or to any interface type. As returns false if err is nil.
func As(err error, target interface{}) bool {
	if target == nil {
		panic("errors: target cannot be nil")
	}
	val := reflect.ValueOf(target)
	typ := val.Type()
	if typ.Kind() != reflect.Ptr || val.IsNil() {
		panic("errors: target must be a non-nil pointer")
	}
	if e := typ.Elem(); e.Kind() != reflect.Interface && !e.Implements(errorType) {
		panic("errors: *target must be interface or implement error")
	}
	targetType := typ.Elem()
	for err != nil {
		if reflect.TypeOf(err).AssignableTo(targetType) {
			val.Elem().Set(reflect.ValueOf(err))
			return true
		}
		if x, ok := err.(interface{ As(interface{}) bool }); ok && x.As(target) {
			return true
		}
		err = Unwrap(err)
	}
	return false
}

var errorType = reflect.TypeOf((*error)(nil)).Elem()