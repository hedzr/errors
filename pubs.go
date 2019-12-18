// Copyright © 2019 Hedzr Yeh.

package errors

import "fmt"

// New ExtErr error object with message and allows attach more nested errors
func New(msg string, args ...interface{}) *ExtErr {
	e := &ExtErr{msg: fmt.Sprintf(msg, args...)}
	return e
}

// NewTemplate ExtErr error object with string template and allows attach more nested errors
func NewTemplate(tmpl string) *ExtErr {
	e := &ExtErr{tmpl: tmpl}
	return e
}

// NewWithError ExtErr error object with nested errors
func NewWithError(errs ...error) *ExtErr {
	return New("unknown error").Attach(errs...)
}

// NewCodedError error object with nested errors
func NewCodedError(code Code, errs ...error) *CodedErr {
	e := &CodedErr{code: code}
	return e.Attach(errs...)
}

// // NewWithCodeMsg ExtErr error object with nested errors
// func NewWithCodeMsg(code Code, msg string, errors ...error) *CodedErr {
// 	return &CodedErr{Code: code, ExtErr: *New(msg, errors...)}
// }
//
// func add(msg string, errs ...error) *ExtErr {
// 	if len(errs) == 0 {
// 		return &ExtErr{msg: msg}
// 	} else if len(errs) == 1 {
// 		err := errs[0]
// 		if e, ok := err.(*ExtErr); ok {
// 			return &ExtErr{msg: msg, innerEE: e}
// 		}
// 		return &ExtErr{msg: msg, innerErr: err}
// 	}
//
// 	return add("", errs[1:]...)
// }

// Walkable interface
type Walkable interface {
	Walk(fn func(err error) (stop bool))
}

// Ranged interface
type Ranged interface {
	Range(fn func(err error) (stop bool))
}

// CanWalk tests if err is walkable
func CanWalk(err error) (ok bool) {
	_, ok = err.(Walkable)
	return
}

// CanRange tests if err is range-able
func CanRange(err error) (ok bool) {
	_, ok = err.(Ranged)
	return
}

// CanUnwrap tests if err is unwrap-able
func CanUnwrap(err error) (ok bool) {
	_, ok = err.(interface{ Unwrap() error })
	return
}

// CanIs tests if err is is-able
func CanIs(err error) (ok bool) {
	_, ok = err.(interface{ Is(error) bool })
	return
}

// CanAs tests if err is as-able
func CanAs(err error) (ok bool) {
	_, ok = err.(interface{ As(interface{}) bool })
	return
}

// Walk will walk all inner and nested error objects inside err
func Walk(err error, fn func(err error) (stop bool)) {
	if !fn(err) {
		if ee, ok := err.(Walkable); ok {
			ee.Walk(fn)
		}
	}
}

// Range can walk the inner/attached errors inside err
func Range(err error, fn func(err error) (stop bool)) {
	if !fn(err) {
		if ee, ok := err.(Ranged); ok {
			ee.Range(fn)
		}
	}
}

// Equal tests if code number presented recursively
func Equal(err error, code Code) bool {
	if x, ok := err.(interface{ EqualRecursive(code Code) bool }); ok {
		return x.EqualRecursive(code)
	}
	return false
}

// IsAny tests if any codes presented
func IsAny(err error, code ...Code) bool {
	if x, ok := err.(interface{ IsAny(codes ...Code) bool }); ok {
		return x.IsAny(code...)
	}
	return false
}

// IsBoth tests if all codes presented
func IsBoth(err error, code ...Code) bool {
	if x, ok := err.(interface{ IsBoth(codes ...Code) bool }); ok {
		return x.IsBoth(code...)
	}
	return false
}

// Attach attaches the nested errors into CodedErr
func Attach(err error, errs ...error) {
	if x, ok := err.(interface{ AttachIts(errors ...error) }); ok {
		x.AttachIts(errs...)
	}
}

// Nest attaches the nested errors into CodedErr
func Nest(err error, errs ...error) {
	if x, ok := err.(interface{ NestIts(errors ...error) }); ok {
		x.NestIts(errs...)
	}
}
