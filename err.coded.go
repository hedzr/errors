// Copyright © 2019 Hedzr Yeh.

package errors

import (
	"bytes"
	"fmt"
)

// CodedErr adds a error code
type CodedErr struct {
	code Code
	ExtErr
}

// Code put another code into CodedErr
func (e *CodedErr) Code(code Code) *CodedErr {
	e.code = code
	return e
}

func (e *CodedErr) Error() string {
	var buf bytes.Buffer
	buf.WriteString(fmt.Sprintf("%06d|%s|", e.code, e.code.String()))
	buf.WriteString(e.ExtErr.Error())
	return buf.String()
}

// Template setup a string format template.
// Coder could compile the error object with formatting args later.
//
// Note that `ExtErr.Template()` had been overrided here
func (e *CodedErr) Template(tmpl string) *CodedErr {
	e.tmpl = tmpl
	return e
}

// Format compiles the final msg with string template and args
//
// Note that `ExtErr.Template()` had been overridden here
func (e *CodedErr) Format(args ...interface{}) *CodedErr {
	if len(args) == 0 {
		e.msg = e.tmpl
	} else {
		e.msg = fmt.Sprintf(e.tmpl, args)
	}
	return e
}

// Msg encodes a formattable msg with args into ExtErr
//
// Note that `ExtErr.Template()` had been overridden here
func (e *CodedErr) Msg(msg string, args ...interface{}) *CodedErr {
	if len(args) == 0 {
		e.msg = msg
	} else {
		e.msg = fmt.Sprintf(msg, args...)
	}
	return e
}

// Attach attaches the nested errors into CodedErr
//
// Note that `ExtErr.Template()` had been overridden here
func (e *CodedErr) Attach(errors ...error) *CodedErr {
	_ = e.add(errors...)
	return e
}

// Nest attaches the nested errors into CodedErr
//
// Note that `ExtErr.Template()` had been overridden here
func (e *CodedErr) Nest(errors ...error) *CodedErr {
	_ = e.nest(errors...)
	return e
}
