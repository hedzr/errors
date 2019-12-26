// Copyright © 2019 Hedzr Yeh.

package errors

import (
	"errors"
	"runtime"
	"strings"
)

// Walkable interface
type Walkable interface {
	Walk(fn func(err error) (stop bool))
}

// Ranged interface
type Ranged interface {
	Range(fn func(err error) (stop bool))
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

// TextContains test if a text fragment is included by err
func TextContains(err error, text string) bool {
	return strings.Index(err.Error(), text) >= 0
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

// DumpStacksAsString returns stack tracing information like debug.PrintStack()
func DumpStacksAsString(allRoutines bool) string {
	buf := make([]byte, 16384)
	buf = buf[:runtime.Stack(buf, allRoutines)]
	// fmt.Printf("=== BEGIN goroutine stack dump ===\n%s\n=== END goroutine stack dump ===\n", buf)
	return string(buf)
}

// HasInnerErrors detects if nested or attached errors present
func HasInnerErrors(err error) (yes bool) {
	if Unwrap(err) != nil {
		return true
	}
	return false
}

// HasAttachedErrors detects if attached errors present
func HasAttachedErrors(err error) (yes bool) {
	if ex, ok := err.(interface{ HasAttachedErrors() bool }); ok {
		return ex.HasAttachedErrors()
	}
	return false
}

// HasWrappedError detects if nested or wrapped errors present
//
// nested error: ExtErr.inner
// wrapped error: fmt.Errorf("... %w ...", err)
func HasWrappedError(err error) (yes bool) {
	if ex, ok := err.(interface{ GetNestedError() *ExtErr }); ok {
		return ex.GetNestedError() != nil
	} else if errors.Unwrap(err) != nil {
		return true
	}
	return false
}
