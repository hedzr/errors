package errors

import (
	"fmt"
	"io"
)

// WithStackInfo is exported now
type WithStackInfo struct {
	causes2

	*Stack

	sites       []interface{}
	taggedSites map[string]interface{}
}

// WithStack annotates err with a Stack trace at the point WithStack was called.
// If err is nil, WithStack returns nil.
func WithStack(cause error) error {
	if cause == nil {
		return nil
	}
	return &WithStackInfo{causes2: causes2{Causers: []error{cause}}, Stack: callers(1)}
}

// End ends the WithXXX stream calls while you dislike unwanted `err =`.
//
// For instance, the construction of an error without warnings looks like:
//
//      err := New("hello %v", "world")
//      _ = err.WithErrors(io.EOF, io.ErrShortWrite).
//          WithErrors(io.ErrClosedPipe).
//          WithCode(Internal)
//
// To avoid the `_ =`, you might belove with a End() call:
//
//      err := New("hello %v", "world")
//      err.WithErrors(io.EOF, io.ErrShortWrite).
//          WithErrors(io.ErrClosedPipe).
//          WithCode(Internal).
//          End()
//
func (w *WithStackInfo) End() {}

func (w *WithStackInfo) rebuild() *WithStackInfo {
	return w
}

// WithCode for error interface
func (w *WithStackInfo) WithCode(code Code) *WithStackInfo {
	w.Code = code
	return w.rebuild()
}

// WithSkip _
func (w *WithStackInfo) WithSkip(skip int) *WithStackInfo {
	w.Stack = callers(skip)
	return w
}

// WithMessage _
func (w *WithStackInfo) WithMessage(message string, args ...interface{}) *WithStackInfo {
	_ = w.causes2.WithMessage(message, args...)
	return w
}

// WithErrors appends errs
// WithStackInfo.Attach() can only wrap and hold one child error object.
func (w *WithStackInfo) WithErrors(errs ...error) *WithStackInfo {
	_ = w.causes2.WithErrors(errs...)
	return w
}

// WithData appends errs if the general object is a error object
func (w *WithStackInfo) WithData(errs ...interface{}) *WithStackInfo {
	if len(errs) > 0 {
		for _, e := range errs {
			if e1, ok := e.(error); ok {
				_ = w.WithErrors(e1)
			} else if e != nil {
				w.sites = append(w.sites, e)
			}
		}
	}
	return w
}

// TaggedData _
type TaggedData map[string]interface{}

// WithTaggedData appends errs if the general object is a error object
func (w *WithStackInfo) WithTaggedData(siteScenes TaggedData) *WithStackInfo {
	if w.taggedSites == nil {
		w.taggedSites = make(TaggedData)
	}
	for k, v := range siteScenes {
		w.taggedSites[k] = v
	}
	return w
}

// WithCause sets the underlying error manually if necessary.
func (w *WithStackInfo) WithCause(cause error) *WithStackInfo {
	w.causes2.Causers = append(w.causes2.Causers, cause)
	return w
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
func (w *WithStackInfo) Cause() error {
	return w.causes2.Cause()
}

// Defer can be used as a defer function to simplify your codes.
//
// The codes:
//
//     func some(){
//       // as a inner errors container
//       child := func() (err error) {
//      	errContainer := errors.New("")
//      	defer errContainer.Defer(&err)
//
//      	for _, r := range []error{io.EOF, io.ErrClosedPipe, errors.Internal} {
//      		errContainer.Attach(r)
//      	}
//
//      	return
//       }
//
//       err := child()
//       t.Logf("failed: %+v", err)
//    }
//
func (w *WithStackInfo) Defer(err *error) {
	*err = w
}

// Format formats the stack of Frames according to the fmt.Formatter interface.
//
//    %s	lists source files for each Frame in the stack
//    %v	lists the source file and line number for each Frame in the stack
//
// Format accepts flags that alter the printing of some verbs, as follows:
//
//    %+v   Prints filename, function, and line number for each Frame in the stack.
func (w *WithStackInfo) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		if s.Flag('+') {
			_, _ = fmt.Fprintf(s, "%+v", w.Error())
			if len(w.sites) > 0 {
				_, _ = fmt.Fprintf(s, "Sites: %+v", w.sites)
			}
			if len(w.taggedSites) > 0 {
				_, _ = fmt.Fprintf(s, "Tagged Sites: %+v", w.taggedSites)
			}
			w.Stack.Format(s, verb)
			return
		}
		_, _ = fmt.Fprintf(s, "%v", w.Error())
	case 's':
		_, _ = io.WriteString(s, w.Error())
	case 'q':
		_, _ = fmt.Fprintf(s, "%q", w.Error())
	}
}

//// Is reports whether any error in `err`'s chain matches target.
//func (w *WithStackInfo) Is(target error) bool {
//	if x, ok := w.error.(interface{ Is(error) bool }); ok && x.Is(target) {
//		return true
//	}
//	return w.error == target
//}

//// TypeIs reports whether any error in `err`'s chain matches target.
//func (w *WithStackInfo) TypeIs(target error) bool {
//	if x, ok := w.error.(interface{ TypeIs(error) bool }); ok && x.TypeIs(target) {
//		return true
//	}
//	return w.error == target
//}

//// As finds the first error in `err`'s chain that matches target, and if so, sets
//// target to that error value and returns true.
//func (w *WithStackInfo) As(target interface{}) bool {
//	return As(w.error, target)
//	//if target == nil {
//	//	panic("errors: target cannot be nil")
//	//}
//	//val := reflect.ValueOf(target)
//	//typ := val.Type()
//	//if typ.Kind() != reflect.Ptr || val.IsNil() {
//	//	panic("errors: target must be a non-nil pointer")
//	//}
//	//if e := typ.Elem(); e.Kind() != reflect.Interface && !e.Implements(errorType) {
//	//	panic("errors: *target must be interface or implement error")
//	//}
//	//targetType := typ.Elem()
//	//err := w.error
//	//for err != nil {
//	//	if reflect.TypeOf(err).AssignableTo(targetType) {
//	//		val.Elem().Set(reflect.ValueOf(err))
//	//		return true
//	//	}
//	//	if x, ok := err.(interface{ As(interface{}) bool }); ok && x.As(target) {
//	//		return true
//	//	}
//	//	err = Unwrap(err)
//	//}
//	//return false
//}

//// Unwrap returns the result of calling the Unwrap method on err, if
//// `err`'s type contains an Unwrap method returning error.
//// Otherwise, Unwrap returns nil.
//func (w *WithStackInfo) Unwrap() error {
//	if w.error != nil {
//		return w.error
//	}
//	//if x, ok := w.error.(interface{ Unwrap() error }); ok {
//	//	return x.Unwrap()
//	//}
//	return nil
//}

//// IsEmpty tests has attached errors
//func (w *WithStackInfo) IsEmpty() bool {
//	if x, ok := w.error.(interface{ IsEmpty() bool }); ok {
//		return x.IsEmpty()
//	}
//	return false
//}
