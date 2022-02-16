// Copyright © 2020 Hedzr Yeh.

package errors

import (
	"fmt"
	"io"
)

// New returns an error with the supplied message.
// New also records the Stack trace at the point it was called.
func New(message string, args ...interface{}) *WithStackInfo {
	if len(args) > 0 {
		message = fmt.Sprintf(message, args...)
	}
	err := &withCause{
		causer: nil,
		msg:    message,
	}
	return &WithStackInfo{
		err,
		callers(),
	}
}

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

//
// ----------
//

// WithCauses holds a group of errors object.
type WithCauses struct {
	causers []error
	msg     string
	*Stack
}

func (w *WithCauses) Error() error {
	if len(w.causers) == 0 {
		return nil
	}
	return w.wrap(w.causers...)
}

func (w *WithCauses) wrap(errs ...error) error {
	return &causes{
		Causers: errs,
		Stack:   w.Stack,
	}
}

// Attach appends errs
func (w *WithCauses) Attach(errs ...error) {
	for _, ex := range errs {
		if ex != nil {
			w.causers = append(w.causers, ex)
		}
	}
	w.Stack = callers()
}

// Cause returns the underlying cause of the error, if possible.
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
func (w *WithCauses) Cause() error {
	if len(w.causers) == 0 {
		return nil
	}
	return w.causers[0]
}

// SetCause sets the underlying error manually if necessary.
func (w *WithCauses) SetCause(cause error) error {
	if cause == nil {
		return nil
	}
	if len(w.causers) == 0 {
		w.causers = append(w.causers, cause)
	} else {
		w.causers[0] = cause
	}
	return w.Cause()
}

// Causes returns the underlying cause of the errors.
func (w *WithCauses) Causes() []error {
	if len(w.causers) == 0 {
		return nil
	}
	return w.causers
}

// Unwrap returns the result of calling the Unwrap method on err, if
// `err`'s type contains an Unwrap method returning error.
// Otherwise, Unwrap returns nil.
func (w *WithCauses) Unwrap() error {
	return w.Cause()
}

// IsEmpty tests has attached errors
func (w *WithCauses) IsEmpty() bool {
	return len(w.causers) == 0
}

//
// ----------
//

// Is reports whether any error in `err`'s chain matches target.
func (w *WithCauses) Is(target error) bool {
	return IsSlice(w.causers, target)
	//if target == nil {
	//	//for _, e := range w.causers {
	//	//	if e == target {
	//	//		return true
	//	//	}
	//	//}
	//	return false
	//}
	//
	//isComparable := reflect.TypeOf(target).Comparable()
	//for {
	//	if isComparable {
	//		for _, e := range w.causers {
	//			if e == target {
	//				return true
	//			}
	//		}
	//		// return false
	//	}
	//
	//	for _, e := range w.causers {
	//		if x, ok := e.(interface{ Is(error) bool }); ok && x.Is(target) {
	//			return true
	//		}
	//		//if err := Unwrap(e); err == nil {
	//		//	return false
	//		//}
	//	}
	//	return false
	//}
}

func (w *WithCauses) TypeIs(target error) bool {
	return TypeIsSlice(w.causers, target)
}

// As finds the first error in `err`'s chain that matches target, and if so, sets
// target to that error value and returns true.
func (w *WithCauses) As(target interface{}) bool {
	return AsSlice(w.causers, target)
	//if target == nil {
	//	panic("errors: target cannot be nil")
	//}
	//val := reflect.ValueOf(target)
	//typ := val.Type()
	//if typ.Kind() != reflect.Ptr || val.IsNil() {
	//	panic("errors: target must be a non-nil pointer")
	//}
	//if e := typ.Elem(); e.Kind() != reflect.Interface && !e.Implements(errorType) {
	//	panic("errors: *target must be interface or implement error")
	//}
	//targetType := typ.Elem()
	//for _, err := range w.causers {
	//	for err != nil {
	//		if reflect.TypeOf(err).AssignableTo(targetType) {
	//			val.Elem().Set(reflect.ValueOf(err))
	//			return true
	//		}
	//		if x, ok := err.(interface{ As(interface{}) bool }); ok && x.As(target) {
	//			return true
	//		}
	//		err = Unwrap(err)
	//	}
	//}
	//return false
}

// WithStackInfo is exported now
type WithStackInfo struct {
	error
	*Stack
}

// WithStack annotates err with a Stack trace at the point WithStack was called.
// If err is nil, WithStack returns nil.
func WithStack(cause error) error {
	if cause == nil {
		return nil
	}
	return &WithStackInfo{cause, callers()}
}

// Cause returns the underlying cause of the error, if possible.
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
func (w *WithStackInfo) Cause() error {
	return w.error
}

// SetCause sets the underlying error manually if necessary.
func (w *WithStackInfo) SetCause(cause error) error {
	w.error = cause
	return w
}

// Format formats the stack of Frames according to the fmt.Formatter interface.
//
//    %s	lists source files for each Frame in the stack
//    %v	lists the source file and line number for each Frame in the stack
//
// Format accepts flags that alter the printing of some verbs, as follows:
//
//    %+v   Prints filename, function, and line number for each Frame in the stack.
func (w *WithStackInfo) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		if s.Flag('+') {
			_, _ = fmt.Fprintf(s, "%+v", w.Cause())
			w.Stack.Format(s, verb)
			return
		}
		fallthrough
	case 's':
		_, _ = io.WriteString(s, w.Error())
	case 'q':
		_, _ = fmt.Fprintf(s, "%q", w.Error())
	}
}

// Is reports whether any error in `err`'s chain matches target.
func (w *WithStackInfo) Is(target error) bool {
	if x, ok := w.error.(interface{ Is(error) bool }); ok && x.Is(target) {
		return true
	}
	return w.error == target
}

// TypeIs reports whether any error in `err`'s chain matches target.
func (w *WithStackInfo) TypeIs(target error) bool {
	if x, ok := w.error.(interface{ TypeIs(error) bool }); ok && x.TypeIs(target) {
		return true
	}
	return w.error == target
}

// As finds the first error in `err`'s chain that matches target, and if so, sets
// target to that error value and returns true.
func (w *WithStackInfo) As(target interface{}) bool {
	return As(w.error, target)
	//if target == nil {
	//	panic("errors: target cannot be nil")
	//}
	//val := reflect.ValueOf(target)
	//typ := val.Type()
	//if typ.Kind() != reflect.Ptr || val.IsNil() {
	//	panic("errors: target must be a non-nil pointer")
	//}
	//if e := typ.Elem(); e.Kind() != reflect.Interface && !e.Implements(errorType) {
	//	panic("errors: *target must be interface or implement error")
	//}
	//targetType := typ.Elem()
	//err := w.error
	//for err != nil {
	//	if reflect.TypeOf(err).AssignableTo(targetType) {
	//		val.Elem().Set(reflect.ValueOf(err))
	//		return true
	//	}
	//	if x, ok := err.(interface{ As(interface{}) bool }); ok && x.As(target) {
	//		return true
	//	}
	//	err = Unwrap(err)
	//}
	//return false
}

// Unwrap returns the result of calling the Unwrap method on err, if
// `err`'s type contains an Unwrap method returning error.
// Otherwise, Unwrap returns nil.
func (w *WithStackInfo) Unwrap() error {
	if w.error != nil {
		return w.error
	}
	//if x, ok := w.error.(interface{ Unwrap() error }); ok {
	//	return x.Unwrap()
	//}
	return nil
}

// Attach appends errs
// WithStackInfo.Attach() can only wrap and hold one child error object.
func (w *WithStackInfo) Attach(errs ...error) *WithStackInfo {
	if w.error == nil {
		if len(errs) > 1 {
			panic("*WithStackInfo.Attach() can only wrap one child error object.")
		}
		for _, e := range errs {
			if e != nil {
				w.error = e
			}
		}
		return w
	}

	if x, ok := w.error.(interface{ Attach(errs ...error) }); ok {
		x.Attach(errs...)
	}

	return w
}

// AttachGenerals appends errs if the general object is a error object
// WithStackInfo.AttachGenerals() can only wrap and hold one child error object.
func (w *WithStackInfo) AttachGenerals(errs ...interface{}) *WithStackInfo {
	if w.error == nil {
		if len(errs) > 1 {
			panic("*WithStackInfo.AttachGenerals() can only wrap one child error object.")
		}
		for _, e := range errs {
			if e1, ok := e.(error); ok {
				w.error = e1
			}
		}
		return w
	}

	if x, ok := w.error.(interface{ AttachGenerals(errs ...interface{}) }); ok {
		x.AttachGenerals(errs...)
	}

	return w
}

// IsEmpty tests has attached errors
func (w *WithStackInfo) IsEmpty() bool {
	if x, ok := w.error.(interface{ IsEmpty() bool }); ok {
		return x.IsEmpty()
	}
	return false
}
