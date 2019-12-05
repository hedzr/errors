// Copyright © 2019 Hedzr Yeh.

package errors

import (
	"bytes"
	"fmt"
)

// New ExtErr error object with message and nested errors
func New(msg string, errors ...error) *ExtErr {
	return add(msg, errors...)
}

// NewWithError ExtErr error object with nested errors
func NewWithError(errors ...error) *ExtErr {
	return add("unknown error", errors...)
}

// NewWithCode ExtErr error object with nested errors
func NewWithCode(code Code, errors ...error) *CodedErr {
	return &CodedErr{Code: code, ExtErr: *NewWithError(errors...)}
}

// NewWithCodeMsg ExtErr error object with nested errors
func NewWithCodeMsg(code Code, msg string, errors ...error) *CodedErr {
	return &CodedErr{Code: code, ExtErr: *New(msg, errors...)}
}

func add(msg string, errs ...error) *ExtErr {
	if len(errs) == 0 {
		return &ExtErr{msg: msg}
	} else if len(errs) == 1 {
		err := errs[0]
		if e, ok := err.(*ExtErr); ok {
			return &ExtErr{msg: msg, innerEE: e}
		}
		return &ExtErr{msg: msg, innerErr: err}
	}

	return add("", errs[1:]...)
}

// CodedErr adds a error code
type CodedErr struct {
	Code Code
	ExtErr
}

// ExtErr is a nestable error object
type ExtErr struct {
	innerEE  *ExtErr
	innerErr error
	msg      string
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

func (e *CodedErr) Error() string {
	var buf bytes.Buffer
	buf.WriteString(fmt.Sprintf("%06d|%s|", e.Code, e.Code.String()))
	buf.WriteString(e.ExtErr.Error())
	return buf.String()
}
