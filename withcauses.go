package errors

// WithCauses holds a group of errors object.
type WithCauses struct {
	causers []error
	msg     string
	*Stack
}

func (w *WithCauses) Defer(err *error) {
	*err = w.Error()
}

func (w *WithCauses) Error() error {
	if len(w.causers) == 0 {
		return nil
	}
	return w.wrap(w.causers...)
}

func (w *WithCauses) wrap(errs ...error) error {
	return &causes{
		Causers: errs,
		Stack:   w.Stack,
	}
}

// Attach appends errs
func (w *WithCauses) Attach(errs ...error) {
	for _, ex := range errs {
		if ex != nil {
			w.causers = append(w.causers, ex)
		}
	}
	w.Stack = callers(1)
}

// Cause returns the underlying cause of the error, if possible.
// An error value has a cause if it implements the following
// interface:
//
//     type causer interface {
//            Cause() error
//     }
//
// If the error does not implement Cause, the original error will
// be returned. If the error is nil, nil will be returned without further
// investigation.
func (w *WithCauses) Cause() error {
	if len(w.causers) == 0 {
		return nil
	}
	return w.causers[0]
}

// SetCause sets the underlying error manually if necessary.
func (w *WithCauses) SetCause(cause error) error {
	if cause == nil {
		return nil
	}
	if len(w.causers) == 0 {
		w.causers = append(w.causers, cause)
	} else {
		w.causers[0] = cause
	}
	return w.Cause()
}

// Causes returns the underlying cause of the errors.
func (w *WithCauses) Causes() []error {
	if len(w.causers) == 0 {
		return nil
	}
	return w.causers
}

// Unwrap returns the result of calling the Unwrap method on err, if
// `err`'s type contains an Unwrap method returning error.
// Otherwise, Unwrap returns nil.
func (w *WithCauses) Unwrap() error {
	return w.Cause()
}

// IsEmpty tests has attached errors
func (w *WithCauses) IsEmpty() bool {
	return len(w.causers) == 0
}

//
// ----------
//

// Is reports whether any error in `err`'s chain matches target.
func (w *WithCauses) Is(target error) bool {
	return IsSlice(w.causers, target)
	//if target == nil {
	//	//for _, e := range w.causers {
	//	//	if e == target {
	//	//		return true
	//	//	}
	//	//}
	//	return false
	//}
	//
	//isComparable := reflect.TypeOf(target).Comparable()
	//for {
	//	if isComparable {
	//		for _, e := range w.causers {
	//			if e == target {
	//				return true
	//			}
	//		}
	//		// return false
	//	}
	//
	//	for _, e := range w.causers {
	//		if x, ok := e.(interface{ Is(error) bool }); ok && x.Is(target) {
	//			return true
	//		}
	//		//if err := Unwrap(e); err == nil {
	//		//	return false
	//		//}
	//	}
	//	return false
	//}
}

func (w *WithCauses) TypeIs(target error) bool {
	return TypeIsSlice(w.causers, target)
}

// As finds the first error in `err`'s chain that matches target, and if so, sets
// target to that error value and returns true.
func (w *WithCauses) As(target interface{}) bool {
	return AsSlice(w.causers, target)
	//if target == nil {
	//	panic("errors: target cannot be nil")
	//}
	//val := reflect.ValueOf(target)
	//typ := val.Type()
	//if typ.Kind() != reflect.Ptr || val.IsNil() {
	//	panic("errors: target must be a non-nil pointer")
	//}
	//if e := typ.Elem(); e.Kind() != reflect.Interface && !e.Implements(errorType) {
	//	panic("errors: *target must be interface or implement error")
	//}
	//targetType := typ.Elem()
	//for _, err := range w.causers {
	//	for err != nil {
	//		if reflect.TypeOf(err).AssignableTo(targetType) {
	//			val.Elem().Set(reflect.ValueOf(err))
	//			return true
	//		}
	//		if x, ok := err.(interface{ As(interface{}) bool }); ok && x.As(target) {
	//			return true
	//		}
	//		err = Unwrap(err)
	//	}
	//}
	//return false
}
