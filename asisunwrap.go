package errors

import (
	"fmt"
	"reflect"
)

// As finds the first error in `err`'s chain that matches target, and if so, sets
// target to that error value and returns true.
//
// The chain consists of err itself followed by the sequence of errors obtained by
// repeatedly calling Unwrap.
//
// An error matches target if the error's concrete value is assignable to the value
// pointed to by target, or if the error has a method As(interface{}) bool such that
// As(target) returns true. In the latter case, the As method is responsible for
// setting target.
//
// As will panic if target is not a non-nil pointer to either a type that implements
// error, or to any interface type. As returns false if err is nil.
func As(err error, target interface{}) bool {
	if target == nil {
		panic("errors: target cannot be nil")
	}
	val := reflect.ValueOf(target)
	typ := val.Type()
	if typ.Kind() != reflect.Ptr || val.IsNil() {
		panic("errors: target must be a non-nil pointer")
	}
	e := typ.Elem()
	k := e.Kind()
	if k != reflect.Interface && k != reflect.Slice && !e.Implements(errorType) {
		// panic("errors: *target must be interface or implement error")
		return false
	}
	targetType := typ.Elem()
	for err != nil {
		if x, ok := err.(interface{ As(interface{}) bool }); ok && x.As(target) {
			return true
		}
		if reflect.TypeOf(err).AssignableTo(targetType) {
			val.Elem().Set(reflect.ValueOf(err))
			return true
		}
		err = Unwrap(err)
	}
	return false
}

// AsSlice tests err.As for errs slice
func AsSlice(errs []error, target interface{}) bool {
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
		if x, ok := err.(interface{ As(interface{}) bool }); ok && x.As(target) {
			return true
		}
		err = Unwrap(err) //nolint:ineffassign,staticcheck
	}
	return false
}

var errorType = reflect.TypeOf((*error)(nil)).Elem()

// Is reports whether any error in `err`'s chain matches target.
//
// The chain consists of err itself followed by the sequence of errors obtained by
// repeatedly calling Unwrap.
//
// An error is considered to match a target if it is equal to that target or if
// it implements a method Is(error) bool such that Is(target) returns true.
func Is(err, target error) bool {
	if target == nil {
		return err == target
	}

	isComparable := reflect.TypeOf(target).Comparable()
	for {
		if isComparable {
			if err == target {
				return true
			}
		}
		if x, ok := err.(interface{ Is(error) bool }); ok && x.Is(target) {
			return true
		}
		if _, ok := target.(Code); !ok {
			if ok = As(err, &target); ok {
				return true
			}
		}
		// TODO: consider supporting target.Is(err). This would allow
		// user-definable predicates, but also may allow for coping with sloppy
		// APIs, thereby making it easier to get away with them.
		if err = Unwrap(err); err == nil {
			return false
		}
	}
}

// IsSlice tests err.Is for errs slice
func IsSlice(errs []error, target error) bool {
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
func TypeIs(err, target error) bool {
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
		// TODO: consider supporting target.Is(err). This would allow
		// user-definable predicates, but also may allow for coping with sloppy
		// APIs, thereby making it easier to get away with them.
		if err = Unwrap(err); err == nil {
			return false
		}
	}
}

// TypeIsSlice tests err.Is for errs slice
func TypeIsSlice(errs []error, target error) bool {
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
//      var err = errors.New("hello").WithErrors(io.EOF, io.ShortBuffers)
//      var e error = err
//      for e != nil {
//          e = errors.Unwrap(err)
//          // test if e is not nil and process it...
//      }
//
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
func Wrap(err error, message string, args ...interface{}) *WithStackInfo {
	if err == nil {
		return nil
	}

	if len(args) > 0 {
		message = fmt.Sprintf(message, args...)
	}

	return &WithStackInfo{
		causes2: causes2{
			Causers: []error{err},
			msg:     message,
		},
		Stack: callers(1),
	}
}
