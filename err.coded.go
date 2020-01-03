// Copyright © 2019 Hedzr Yeh.

package errors

import (
	"bytes"
	"fmt"
	"strconv"
)

// CodedErr adds a error code
type CodedErr struct {
	code Code
	ExtErr
}

// NoCannedError detects mqttError object is not an error or not an canned-error (inners is empty)
func (e *CodedErr) NoCannedError() bool {
	return e.Number() == OK || e.HasAttachedErrors()
}

// HasAttachedErrors tests if any errors attached (nor nested) to `e` or not
func (e *CodedErr) HasAttachedErrors() bool {
	return len(e.GetAttachedErrors()) != 0
}

// Code put another code into CodedErr
func (e *CodedErr) Code(code Code) *CodedErr {
	e.code = code
	return e
}

// Equal compares with code
func (e *CodedErr) Equal(code Code) bool {
	return e.code == Code(code)
}

// EqualRecursive compares with code
func (e *CodedErr) EqualRecursive(code Code) bool {
	if e.Equal(code) {
		return true
	}

	b := false
	Walk(e, func(err error) (stop bool) {
		// log.Printf("  ___E : %+v", err)
		if c, ok := err.(interface{ Equal(code Code) bool }); ok {
			if c.Equal(Code(code)) {
				b = true
				return true
			}
		}
		return false
	})
	return b
}

// Number returns the code number
func (e *CodedErr) Number() Code {
	return e.code
}

// IsBoth tests if all codes presented
func (e *CodedErr) IsBoth(code ...Code) bool {
	for _, c := range code {
		if !e.EqualRecursive(c) {
			return false
		}
	}
	return true
}

// IsAny tests if any codes presented
func (e *CodedErr) IsAny(code ...Code) bool {
	for _, c := range code {
		if e.EqualRecursive(c) {
			return true
		}
	}
	return false
}

// Error for stringer interface
func (e *CodedErr) Error() string {
	var buf bytes.Buffer
	var s = strconv.Itoa(int(e.code))
	buf.WriteString(LeftPad(s, '0', 6))
	buf.WriteRune('|')
	// buf.WriteString(strconv.Itoa(int(e.code)))
	// buf.WriteRune('|')
	// buf.WriteString(fmt.Sprintf("%06d|", e.code))
	s = e.code.String()
	if len(s) > 0 {
		buf.WriteString(s)
		buf.WriteString("|")
	}
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
		e.msg = fmt.Sprintf(e.tmpl, args...)
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

// AttachIts attaches the nested errors into CodedErr
func (e *CodedErr) AttachIts(errors ...error) {
	_ = e.add(errors...)
}

// NestIts attaches the nested errors into CodedErr
func (e *CodedErr) NestIts(errors ...error) {
	_ = e.nest(errors...)
}
