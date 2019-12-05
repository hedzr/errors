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
	innerEE  *ExtErr
	innerErr error
	msg      string
	tmpl     string
}

// Unwrap returns the result of calling the Unwrap method on err, if err's
// type contains an Unwrap method returning error.
// Otherwise, Unwrap returns nil.
func (e *ExtErr) Unwrap() error {
	if e.innerErr != nil {
		return e.innerErr
	}
	if e.innerEE != nil {
		return e.innerEE
	}
	return nil
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

// Attach attaches the nested errors into ExtErr
func (e *ExtErr) Attach(errors ...error) *ExtErr {
	return e.add(errors...)
}

// Nest attaches the nested errors into ExtErr
func (e *ExtErr) Nest(errors ...error) *ExtErr {
	return e.add(errors...)
}

func (e *ExtErr) add(errs ...error) *ExtErr {
	switch len(errs) {
	case 0:
	case 1:
		err := errs[0]
		if e, ok := err.(*ExtErr); ok {
			return &ExtErr{innerEE: e}
		}
		e.innerErr = err
	default:
		return e.add(errs[1:]...)
	}
	return e
}

func (e *ExtErr) Error() string {
	var buf bytes.Buffer
	if len(e.msg) == 0 {
		buf.WriteString("error")
	} else {
		buf.WriteString(e.msg)
	}
	if e.innerErr != nil {
		// buf.WriteString("[")
		buf.WriteString(", ")
		buf.WriteString(e.innerErr.Error())
		// buf.WriteString("]")
	}
	if e.innerEE != nil {
		buf.WriteString("[")
		buf.WriteString(e.innerEE.Error())
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
	_ = e.add(errors...)
	return e
}

func (e *CodedErr) Error() string {
	var buf bytes.Buffer
	buf.WriteString(fmt.Sprintf("%06d|%s|", e.code, e.code.String()))
	buf.WriteString(e.ExtErr.Error())
	return buf.String()
}
