// Copyright © 2019 Hedzr Yeh.

package errors

import (
	"bytes"
	"fmt"
)

// ExtErr is a nestable error object
type ExtErr struct {
	inner *ExtErr
	errs  []error
	msg   string
	tmpl  string
}

// GetTemplateString returns e.tmpl member
func (e *ExtErr) GetTemplateString() string {
	return e.tmpl
}

// GetMsgString returns e.msg member
func (e *ExtErr) GetMsgString() string {
	return e.msg
}

// GetNestedError returns e.inner member (nested errors)
func (e *ExtErr) GetNestedError() *ExtErr {
	return e.inner
}

// GetAttachedErrors returns e.errs member (attached errors)
func (e *ExtErr) GetAttachedErrors() []error {
	return e.errs
}

// NoCannedError detects mqttError object is not an error or not an canned-error (inners is empty)
func (e *ExtErr) NoCannedError() bool {
	return e.HasAttachedErrors()
}

// HasAttachedErrors tests if any errors attached (nor nested) to `e` or not
func (e *ExtErr) HasAttachedErrors() bool {
	return len(e.errs) == 0
}

// Walk will walk all inner/attached and nested error objects inside e
func (e *ExtErr) Walk(fn func(err error) (stop bool)) {
	for _, ee := range e.errs {
		if fn(ee) {
			return
		}
		if ex, ok := ee.(Walkable); ok {
			ex.Walk(fn)
		}
	}
	if e.inner != nil {
		if !fn(e.inner) {
			e.inner.Walk(fn)
		}
	}
}

// Range can walk the inner/attached errors inside e
func (e *ExtErr) Range(fn func(err error) (stop bool)) {
	for _, ee := range e.errs {
		if fn(ee) {
			return
		}
	}
}

// Unwrap returns the result of calling the Unwrap method on err, if err's
// type contains an Unwrap method returning error.
// Otherwise, Unwrap returns nil.
func (e *ExtErr) Unwrap() error {
	if e.inner != nil {
		return e.inner
	}

	for _, ee := range e.errs {
		// if x, ok := ee.(interface{ Unwrap() error }); ok {
		// 	return x.Unwrap()
		// }
		return ee
	}
	return nil
}

// Is reports whether any error in err's chain matches target.
func (e *ExtErr) Is(err error) bool {
	if e == err {
		return true
	}

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

// AttachIts attaches the nested errors into ExtErr
func (e *ExtErr) AttachIts(errors ...error) {
	_ = e.add(errors...)
}

// NestIts attaches the nested errors into ExtErr
func (e *ExtErr) NestIts(errors ...error) {
	_ = e.nest(errors...)
}

func (e *ExtErr) nest(errs ...error) *ExtErr {
	z := e
	for {
		if z.inner != nil {
			z = z.inner
		} else if len(z.errs) == 0 {
			z.errs = errs
			return e // z
		} else if errs[0] != nil {
			z.inner = &ExtErr{errs: errs}
			return e // z
		}
	}
}

func (e *ExtErr) add(errs ...error) *ExtErr {
	if len(errs) > 0 && errs[0] != nil {
		e.errs = append(e.errs, errs...)
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
