package errors

import (
	"fmt"
	"reflect"
	"strings"
)

// As finds the first error in `err`'s chain that matches target,
// and if so, sets target to that error value and returns true.
//
// The chain consists of err itself followed by the sequence of errors
// obtained by repeatedly calling Unwrap.
//
// An error matches target if the error's concrete value is assignable
// to the value pointed to by target, or if the error has a method
// As(interface{}) bool such that As(target) returns true. In the
// latter case, the As method is responsible for setting target.
//
// As will panic if target is not a non-nil pointer to either a
// type that implements error, or to any interface type. "As"
// returns false if err is nil.
func As(err error, target interface{}) bool { //nolint:revive
	if target == nil {
		panic("errors: target cannot be nil")
	}
	val := reflect.ValueOf(target)
	typ := val.Type()
	if typ.Kind() != reflect.Ptr || val.IsNil() {
		panic("errors: target must be a non-nil pointer")
	}
	// e := typ.Elem()
	// k := e.Kind()
	// if k != reflect.Interface && k != reflect.Slice && !e.Implements(errorType) {
	// 	// panic("errors: *target must be interface or implement error")
	// 	return false
	// }
	targetType := typ.Elem()
	for err != nil {
		if x, ok := err.(interface{ As(interface{}) bool }); ok && x.As(target) { //nolint:revive
			return true
		}
		if reflect.TypeOf(err).AssignableTo(targetType) {
			val.Elem().Set(reflect.ValueOf(err))
			return true
		}
		err = Unwrap(err) //nolint:revive
	}
	return false
}

// AsSlice tests err.As for errs slice
func AsSlice(errs []error, target interface{}) bool { //nolint:revive
	if target == nil {
		panic("errors: target cannot be nil")
	}
	val := reflect.ValueOf(target)
	typ := val.Type()
	if typ.Kind() != reflect.Ptr || val.IsNil() {
		panic("errors: target must be a non-nil pointer")
	}
	if e := typ.Elem(); e.Kind() != reflect.Interface && !e.Implements(errorType) {
		// panic("errors: *target must be interface or implement error")
		return false
	}
	targetType := typ.Elem()
	for _, err := range errs {
		if reflect.TypeOf(err).AssignableTo(targetType) {
			val.Elem().Set(reflect.ValueOf(err))
			return true
		}
		if x, ok := err.(interface{ As(interface{}) bool }); ok && x.As(target) { //nolint:revive
			return true
		}
		err = Unwrap(err) //nolint:ineffassign,staticcheck
	}
	return false
}

var errorType = reflect.TypeOf((*error)(nil)).Elem()

// IsAnyOf tests whether any of `targets` is in `err`.
func IsAnyOf(err error, targets ...error) bool {
	for _, tgt := range targets {
		if Is(err, tgt) {
			return true
		}
	}
	return false
}

// Is reports whether any error in `err`'s chain matches target.
//
// The chain consists of err itself followed by the sequence of
// errors obtained by repeatedly calling Unwrap.
//
// An error is considered to match a target if it is equal to that
// target or if it implements a method Is(error) bool such that
// Is(target) returns true.
func Is(err, target error) bool { //nolint:revive
	if target == nil {
		return err == nil
	}

	isComparable := reflect.TypeOf(target).Comparable()
	tv := reflect.ValueOf(target)
	// target is not Code-based, try convert source err with target's type, and test whether its plain text message is equal
	var savedMsg string
	if !isNil(tv) {
		savedMsg = target.Error()
	}
	for {
		if isComparable && err == target {
			return true
		}
		if x, ok := err.(interface{ Is(error) bool }); ok && x.Is(target) {
			return true
		}
		if _, ok := target.(Code); !ok {
			var te Code
			if ok = As(err, &te); ok && !isNil(reflect.ValueOf(err)) && strings.EqualFold(te.Error(), savedMsg) {
				return true
			}
		}

		// // TODO: consider supporting target.Is(err). This would allow
		// // user-definable predicates, but also may allow for coping with sloppy
		// // APIs, thereby making it easier to get away with them.
		// if err = Unwrap(err); err == nil {
		// 	errors.Is()
		// 	return false
		// }

		switch x := err.(type) {
		case interface{ Unwrap() error }:
			err = x.Unwrap() //nolint:revive
			if err == nil {
				return false
			}
		case interface{ Unwrap() []error }:
			for _, err := range x.Unwrap() {
				if Is(err, target) {
					return true
				}
			}
			return false
		default:
			return false
		}
	}
}

func IsStd(err, target error) bool { //nolint:revive
	if target == nil {
		return err == target
	}

	isComparable := reflect.TypeOf(target).Comparable()
	for {
		if isComparable && err == target {
			return true
		}
		if x, ok := err.(interface{ Is(error) bool }); ok && x.Is(target) {
			return true
		}
		switch x := err.(type) {
		case interface{ Unwrap() error }:
			err = x.Unwrap() //nolint:revive
			if err == nil {
				return false
			}
		case interface{ Unwrap() []error }:
			for _, err := range x.Unwrap() {
				if Is(err, target) {
					return true
				}
			}
			return false
		default:
			return false
		}
	}
}

func Iss(err error, targets ...error) (matched bool) { //nolint:revive
	if targets == nil {
		return err == nil
	}
	if err == nil {
		return true
	}

	for _, target := range targets {
		isComparable := reflect.TypeOf(target).Comparable()
		tv := reflect.ValueOf(target)
		// target is not Code-based, try convert source err with target's type, and test whether its plain text message is equal
		var savedMsg string
		if !isNil(tv) {
			savedMsg = target.Error()
		}
		for {
			if isComparable && err == target {
				return true
			}
			if x, ok := err.(interface{ Is(error) bool }); ok && x.Is(target) {
				return true
			}
			if _, ok := target.(Code); !ok {
				if ok = As(err, &target); ok && !isNil(reflect.ValueOf(target)) && strings.EqualFold(target.Error(), savedMsg) {
					return true
				}
			}

			// // TODO: consider supporting target.Is(err). This would allow
			// // user-definable predicates, but also may allow for coping with sloppy
			// // APIs, thereby making it easier to get away with them.
			// if err = Unwrap(err); err == nil {
			// 	return false
			// }

			switch x := err.(type) {
			case interface{ Unwrap() error }:
				err = x.Unwrap() //nolint:revive
				if err == nil {
					return false
				}
			case interface{ Unwrap() []error }:
				for _, err := range x.Unwrap() {
					if Is(err, target) {
						return true
					}
				}
				return false
			default:
				return false
			}
		}
	}
	return
}

// isNil for go1.12+, the difference is it never panic on unavailable kinds.
// see also reflect.IsNil.
func isNil(v reflect.Value) bool {
	return isNilv(&v)
}

// IsNilv for go1.12+, the difference is it never panic on unavailable kinds.
// see also reflect.IsNil.
func isNilv(v *reflect.Value) bool {
	if v != nil {
		switch k := v.Kind(); k { //nolint:exhaustive //no need
		case reflect.Uintptr:
			if v.CanAddr() {
				return v.UnsafeAddr() == 0 // special: reflect.IsNil assumed nil check on an uintptr is illegal, faint!
			}
		case reflect.Chan, reflect.Func, reflect.Map, reflect.Ptr:
			return v.IsNil()
		case reflect.UnsafePointer:
			return v.Pointer() == 0 // for go1.11, this is a workaround even not bad
		case reflect.Interface, reflect.Slice:
			return v.IsNil()
			// case reflect.Array:
			//	// never true, for an array, it is never IsNil
			// case reflect.String:
			// case reflect.Struct:
		}
	}
	return false
}

// IsSlice tests err.Is for errs slice
func IsSlice(errs []error, target error) bool { //nolint:revive
	if target == nil {
		// for _, e := range errs {
		//	if e == target {
		//		return true
		//	}
		// }
		return false
	}

	isComparable := reflect.TypeOf(target).Comparable()
	for {
		if isComparable {
			for _, e := range errs {
				if e == target {
					return true
				}
			}
			// return false
		}

		for _, e := range errs {
			if x, ok := e.(interface{ Is(error) bool }); ok && x.Is(target) {
				return true
			}
			// if err := Unwrap(e); err == nil {
			//	return false
			// }
		}
		return false //nolint:staticcheck
	}
}

// TypeIs reports whether any error in `err`'s chain matches target.
//
// The chain consists of err itself followed by the sequence of errors obtained by
// repeatedly calling Unwrap.
//
// An error is considered to match a target if it is equal to that target or if
// it implements a method Is(error) bool such that Is(target) returns true.
func TypeIs(err, target error) bool { //nolint:revive
	if target == nil {
		return err == target
	}

	isComparable := reflect.TypeOf(target).Comparable()
	for {
		if isComparable {
			if reflect.TypeOf(target) == reflect.TypeOf(err) {
				return true
			}
		}
		if x, ok := err.(interface{ Is(error) bool }); ok && x.Is(target) {
			return true
		}

		// // TODO: consider supporting target.Is(err). This would allow
		// // user-definable predicates, but also may allow for coping with sloppy
		// // APIs, thereby making it easier to get away with them.
		// if err = Unwrap(err); err == nil {
		// 	return false
		// }

		switch x := err.(type) {
		case interface{ Unwrap() error }:
			err = x.Unwrap() //nolint:revive
			if err == nil {
				return false
			}
		case interface{ Unwrap() []error }:
			for _, err := range x.Unwrap() {
				if Is(err, target) {
					return true
				}
			}
			return false
		default:
			return false
		}
	}
}

// TypeIsSlice tests err.Is for errs slice
func TypeIsSlice(errs []error, target error) bool { //nolint:revive
	if target == nil {
		// for _, e := range errs {
		//	if e == target {
		//		return true
		//	}
		// }
		return false
	}

	isComparable := reflect.TypeOf(target).Comparable()
	for {
		if isComparable {
			tt := reflect.TypeOf(target)
			for _, e := range errs {
				// if e == target {
				//	return true
				// }
				if reflect.TypeOf(e) == tt {
					return true
				}
			}
			// return false
		}

		for _, e := range errs {
			if x, ok := e.(interface{ Is(error) bool }); ok && x.Is(target) {
				return true
			}
			// if err := Unwrap(e); err == nil {
			//	return false
			// }
		}
		return false //nolint:staticcheck
	}
}

// func As(err error, target interface{}) bool
// func Is(err, target error) bool
// func New(text string) error
// func Unwrap(err error) error

// Unwrap returns the result of calling the Unwrap method on err, if
// `err`'s type contains an Unwrap method returning error.
// Otherwise, Unwrap returns nil.
//
// An errors.Error is an unwrappable error object, all its inner errors
// can be unwrapped in turn. Therefore it maintains an internal unwrapping
// index and it can't be reset externally. The only approach to clear it
// and make Unwrap work from head is, to keep Unwrap till this turn ending
// by returning nil.
//
// Examples for Unwrap:
//
//	var err = errors.New("hello").WithErrors(io.EOF, io.ShortBuffers)
//	var e error = err
//	for e != nil {
//	    e = errors.Unwrap(err)
//	    // test if e is not nil and process it...
//	}
func Unwrap(err error) error {
	u, ok := err.(interface {
		Unwrap() error
	})
	if !ok {
		return nil
	}
	return u.Unwrap()
}

// Wrap returns an error annotating err with a Stack trace
// at the point Wrap is called, and the supplied message.
// If err is nil, Wrap returns nil.
func Wrap(err error, message string, args ...interface{}) *WithStackInfo { //nolint:revive
	if err == nil {
		return nil
	}

	if len(args) > 0 {
		message = fmt.Sprintf(message, args...) //nolint:revive
	}

	return &WithStackInfo{
		causes2: causes2{
			Causers: []error{err},
			msg:     message,
		},
		Stack: callers(1),
	}
}
