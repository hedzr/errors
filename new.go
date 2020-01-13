// Copyright © 2019 Hedzr Yeh.

package errors

import "fmt"

// New ExtErr error object with message and allows attach more nested errors.
// New also records the stack trace at the point it was called.
func New(format string, args ...interface{}) *ExtErr {
	if len(args) == 0 {
		return &ExtErr{msg: format, stack: callers()}
	}
	e := &ExtErr{msg: fmt.Sprintf(format, args...), stack: callers()}
	return e
}

// NewTemplate ExtErr error object with string template and allows attach more nested errors
func NewTemplate(tmpl string) *ExtErr {
	e := &ExtErr{tmpl: tmpl}
	return e
}

// NewWithError ExtErr error object with nested errors
// NewWithError also records the stack trace at the point it was called.
func NewWithError(errs ...error) *ExtErr {
	return New("unknown error").Attach(errs...)
}

// NewCodedError error object with nested errors
// NewCodedError also records the stack trace at the point it was called.
func NewCodedError(code Code, errs ...error) *CodedErr {
	e := &CodedErr{code: code, ExtErr: ExtErr{stack: callers()}}
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
