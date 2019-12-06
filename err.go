// Copyright © 2019 Hedzr Yeh.

package errors

import (
	"bytes"
	"fmt"
)

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
func NewWithError(errors ...error) *ExtErr {
	return New("unknown error").Attach(errors...)
}

// NewCodedError error object with nested errors
func NewCodedError(code Code) *CodedErr {
	return &CodedErr{code: code}
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

// CodedErr adds a error code
type CodedErr struct {
	code Code
	ExtErr
}

// ExtErr is a nestable error object
type ExtErr struct {
	inner *ExtErr
	errs  []error
	msg   string
	tmpl  string
}

// Unwrap returns the result of calling the Unwrap method on err, if err's
// type contains an Unwrap method returning error.
// Otherwise, Unwrap returns nil.
func (e *ExtErr) Unwrap() error {
	if e.inner != nil {
		return e.inner
	}

	for _, ee := range e.errs {
		if x, ok := ee.(interface{ Unwrap() error }); ok {
			return x.Unwrap()
		}
	}
	return nil
}

// Is reports whether any error in err's chain matches target.
func (e *ExtErr) Is(err error) bool {
	if e.inner != nil {
		if e.inner == err {
			return true
		}
		if e.inner.Is(err) {
			return true
		}
	}

	for _, ee := range e.errs {
		if ee == err {
			return true
		}
		if i, ok := ee.(interface{ Is(error) bool }); ok && i.Is(err) {
			return true
		}
	}
	return false
}

// As finds the first error in err's chain that matches target, and if so, sets
// target to that error value and returns true.
func (e *ExtErr) As(target interface{}) bool {
	if e.inner != nil {
		if As(e.inner, target) {
			return true
		}
	}

	for _, ee := range e.errs {
		if i, ok := ee.(interface{ As(interface{}) bool }); ok && i.As(target) {
			return true
		}
	}
	return false
}

// Template setup a string format template.
// Coder could compile the error object with formatting args later.
func (e *ExtErr) Template(tmpl string) *ExtErr {
	e.tmpl = tmpl
	return e
}

// Format compiles the final msg with string template and args
func (e *ExtErr) Format(args ...interface{}) *ExtErr {
	if len(args) == 0 {
		e.msg = e.tmpl
	} else {
		e.msg = fmt.Sprintf(e.tmpl, args)
	}
	return e
}

// Msg encodes a formattable msg with args into ExtErr
func (e *ExtErr) Msg(msg string, args ...interface{}) *ExtErr {
	if len(args) == 0 {
		e.msg = msg
	} else {
		e.msg = fmt.Sprintf(msg, args...)
	}
	return e
}

// Attach attaches a group of errors into ExtErr
func (e *ExtErr) Attach(errors ...error) *ExtErr {
	return e.add(errors...)
}

// Nest attaches the nested errors into ExtErr
func (e *ExtErr) Nest(errors ...error) *ExtErr {
	return e.nest(errors...)
}

func (e *ExtErr) nest(errs ...error) *ExtErr {
	z := e
	for {
		if z.inner != nil {
			z = z.inner
		} else if len(z.errs) == 0 {
			z.errs = errs
			return e // z
		} else {
			z.inner = &ExtErr{errs: errs}
			return e // z
		}
	}
}

func (e *ExtErr) add(errs ...error) *ExtErr {
	e.errs = append(e.errs, errs...)
	return e
}

func (e *ExtErr) Error() string {
	var buf bytes.Buffer
	if len(e.msg) == 0 {
		buf.WriteString("error")
	} else {
		buf.WriteString(e.msg)
	}

	for _, ee := range e.errs {
		// buf.WriteString("[")
		buf.WriteString(", ")
		buf.WriteString(ee.Error())
		// buf.WriteString("]")
	}
	if e.inner != nil {
		buf.WriteString("[")
		buf.WriteString(e.inner.Error())
		buf.WriteString("]")
	}
	return buf.String()
}

// Code put another code into CodedErr
func (e *CodedErr) Code(code Code) *CodedErr {
	e.code = code
	return e
}

// Template setup a string format template.
// Coder could compile the error object with formatting args later.
func (e *CodedErr) Template(tmpl string) *CodedErr {
	e.tmpl = tmpl
	return e
}

// Format compiles the final msg with string template and args
func (e *CodedErr) Format(args ...interface{}) *CodedErr {
	if len(args) == 0 {
		e.msg = e.tmpl
	} else {
		e.msg = fmt.Sprintf(e.tmpl, args)
	}
	return e
}

// Msg encodes a formattable msg with args into ExtErr
func (e *CodedErr) Msg(msg string, args ...interface{}) *CodedErr {
	if len(args) == 0 {
		e.msg = msg
	} else {
		e.msg = fmt.Sprintf(msg, args...)
	}
	return e
}

// Attach attaches the nested errors into CodedErr
func (e *CodedErr) Attach(errors ...error) *CodedErr {
	_ = e.add(errors...)
	return e
}

// Nest attaches the nested errors into CodedErr
func (e *CodedErr) Nest(errors ...error) *CodedErr {
	_ = e.nest(errors...)
	return e
}

func (e *CodedErr) Error() string {
	var buf bytes.Buffer
	buf.WriteString(fmt.Sprintf("%06d|%s|", e.code, e.code.String()))
	buf.WriteString(e.ExtErr.Error())
	return buf.String()
}
