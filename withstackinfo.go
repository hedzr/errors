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

func (w *WithStackInfo) IsDescended(descendant error) bool {
	if e, ok := descendant.(*WithStackInfo); ok {
		return e.Code == w.Code && e.msg == w.msg
	}
	return false
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

// Data returns the wrapped common user data by WithData.
// The error objects with passed WithData will be moved into inner
// errors set, so its are excluded from Data().
func (w *WithStackInfo) Data() []interface{} { return w.sites }

// TaggedData returns the wrapped tagged user data by WithTaggedData.
func (w *WithStackInfo) TaggedData() TaggedData { return w.taggedSites }

// Cause returns the underlying cause of the error, if possible.
// An error value has a cause if it implements the following
// interface:
//
//     type causer interface {
//            Cause() error
//     }
//
// If an error object does not implement Cause interface, the
// original error object will be returned.
// If the error is nil, nil will be returned without further
// investigation.
func (w *WithStackInfo) Cause() error {
	return w.causes2.Cause()
}

func (w *WithStackInfo) rebuild() Buildable {
	return w
}

// WithSkip specifies a special number of stack frames that will be ignored.
func (w *WithStackInfo) WithSkip(skip int) Buildable {
	w.Stack = callers(skip)
	return w
}

// WithMessage formats the error message
func (w *WithStackInfo) WithMessage(message string, args ...interface{}) Buildable {
	_ = w.causes2.WithMessage(message, args...)
	return w
}

// WithCode specifies an error code.
// An error code `Code` is a integer number with error interface
// supported.
func (w *WithStackInfo) WithCode(code Code) Buildable {
	w.Code = code
	return w.rebuild()
}

// WithErrors attaches the given errs as inner errors.
// WithErrors is like our old Attach().
// It wraps the inner errors into underlying container and
// represents them all in a singular up-level error object.
// The wrapped inner errors can be retrieved with errors.Causes:
//
//      var err = errors.New("hello").WithErrors(io.EOF, io.ShortBuffers)
//      var errs []error = errors.Causes(err)
//
// Or, use As() to extract its:
//
//      var errs []error
//      errors.As(err, &errs)
//
func (w *WithStackInfo) WithErrors(errs ...error) Buildable {
	_ = w.causes2.WithErrors(errs...)

	//for _, e := range errs {
	//	if e1, ok := e.(*WithStackInfo); ok {
	//		w.Stack = e1.Stack
	//	}
	//}
	return w
}

// WithData appends errs if the general object is a error object.
// It can be used in defer-recover block typically. For example:
//
//    defer func() {
//      if e := recover(); e != nil {
//        err = errors.New("[recovered] copyTo unsatisfied ([%v] %v -> [%v] %v), causes: %v",
//          c.indirectType(from.Type()), from, c.indirectType(to.Type()), to, e).
//          WithData(e)
//        n := log.CalcStackFrames(1)   // skip defer-recover frame at first
//        log.Skip(n).Errorf("%v", err) // skip go-lib frames and defer-recover frame, back to the point throwing panic
//      }
//    }()
//
func (w *WithStackInfo) WithData(errs ...interface{}) Buildable {
	if len(errs) > 0 {
		for _, e := range errs {
			if e1, ok := e.(error); ok {
				_ = w.WithErrors(e1)
				if e1, ok := e.(*WithStackInfo); ok {
					w.Stack = e1.Stack
				}
			} else if e != nil {
				w.sites = append(w.sites, e)
			}
		}
	}
	return w
}

// WithTaggedData appends errs if the general object is a error object
func (w *WithStackInfo) WithTaggedData(siteScenes TaggedData) Buildable {
	if w.taggedSites == nil {
		w.taggedSites = make(TaggedData)
	}
	for k, v := range siteScenes {
		w.taggedSites[k] = v
	}
	return w
}

// WithCause sets the underlying error manually if necessary.
func (w *WithStackInfo) WithCause(cause error) Buildable {
	w.causes2.Causers = append(w.causes2.Causers, cause)
	return w
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
	if !w.IsEmpty() || w.Code != OK {
		*err = w
	}
}

// Attach collects the errors except it's nil
//
// Since v3.0.5, we break Attach() and remove its returning value.
// So WithStackInfo is a Container compliant type now.
func (w *WithStackInfo) Attach(errs ...error) {
	_ = w.WithErrors(errs...)
}

// FormatWith _
func (w *WithStackInfo) FormatWith(args ...interface{}) error {
	c := w.Clone()
	c.liveArgs = args
	return c
}

// Clone _
func (w *WithStackInfo) Clone() *WithStackInfo {
	c := &WithStackInfo{
		causes2: causes2{
			Code:        w.causes2.Code,
			Causers:     w.causes2.Causers,
			msg:         w.causes2.msg,
			unwrapIndex: w.causes2.unwrapIndex,
			liveArgs:    w.causes2.liveArgs,
		},
		Stack:       w.Stack,
		sites:       w.sites,
		taggedSites: w.taggedSites,
	}
	return c
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
			n, _ := fmt.Fprintf(s, "%+v", w.Error())
			if len(w.sites) > 0 {
				if n > 0 {
					n1, _ := fmt.Fprintf(s, "\n  ")
					n += n1
				}
				n1, _ := fmt.Fprintf(s, "Sites: %+v", w.sites)
				n += n1
			}
			if len(w.taggedSites) > 0 {
				if n > 0 {
					_, _ = fmt.Fprintf(s, "\n  ")
				}
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
