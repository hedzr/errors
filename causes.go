package errors

import (
	"bytes"
	"fmt"
	"io"
)

type causes2 struct {
	Code
	Causers []error
	msg     string
}

// WithCode for error interface
func (c *causes2) WithCode(code Code) *causes2 {
	c.Code = code
	return c
}

func (w *causes2) WithMessage(message string, args ...interface{}) *causes2 {
	if len(args) > 0 {
		message = fmt.Sprintf(message, args...)
	}
	w.msg = message
	return w
}

// End ends the WithXXX stream calls while you dislike unwanted `err =`.
func (w *causes2) End() {}

// WithAttach appends errs
// WithStackInfo.Attach() can only wrap and hold one child error object.
func (w *causes2) WithAttach(errs ...error) *causes2 {
	if len(errs) > 0 {
		for _, e := range errs {
			if e != nil {
				w.Causers = append(w.Causers, e)
			}
		}
	}
	return w
}

func (w *causes2) Error() string {
	var buf bytes.Buffer
	if w.msg != "" {
		buf.WriteString(w.msg)
	}
	var needclose, needsep bool
	if w.Code != OK {
		if buf.Len() > 0 {
			buf.WriteRune(' ')
		}
		buf.WriteString("[")
		buf.WriteString(w.Code.String())
		needclose = true
		needsep = true
	} else if len(w.Causers) > 0 {
		if buf.Len() > 0 {
			buf.WriteRune(' ')
		}
		buf.WriteString("[")
		needclose = true
	}
	for i, c := range w.Causers {
		if i > 0 || needsep {
			buf.WriteString(" | ")
		}
		buf.WriteString(c.Error())
	}
	if needclose {
		buf.WriteString("]")
	}
	// buf.WriteString(w.Stack)
	return buf.String()
}

// Format formats the stack of Frames according to the fmt.Formatter interface.
//
//    %s	lists source files for each Frame in the stack
//    %v	lists the source file and line number for each Frame in the stack
//
// Format accepts flags that alter the printing of some verbs, as follows:
//
//    %+v   Prints filename, function, and line number for each Frame in the stack.
func (w *causes2) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		if s.Flag('+') {
			_, _ = fmt.Fprintf(s, "%+v", w.Error())
			return
		}
		fallthrough
	case 's':
		_, _ = io.WriteString(s, w.Error())
	case 'q':
		_, _ = fmt.Fprintf(s, "%q", w.Error())
	}
}

func (w *causes2) Cause() error {
	if len(w.Causers) == 0 {
		return nil
	}
	return w.Causers[0]
}

func (w *causes2) Causes() []error {
	if len(w.Causers) == 0 {
		return nil
	}
	return w.Causers
}

func (w *causes2) Unwrap() error {
	// return w.Cause()

	for _, err := range w.Causers {
		//u, ok := err.(interface {
		//	Unwrap() error
		//})
		//if ok {
		//	return u.Unwrap()
		//}
		return err // just return the first cause
	}
	return nil
}

func (w *causes2) Is(target error) bool {
	return IsSlice(w.Causers, target)
}

func (w *causes2) TypeIs(target error) bool {
	return TypeIsSlice(w.Causers, target)
}

// As finds the first error in `err`'s chain that matches target, and if so, sets
// target to that error value and returns true.
func (w *causes2) As(target interface{}) bool {
	return AsSlice(w.Causers, target)
}

// IsEmpty tests has attached errors
func (w *causes2) IsEmpty() bool {
	return len(w.Causers) == 0
}

//
//
//

type causes struct {
	Causers []error
	*Stack
}

func (w *causes) Error() string {
	var buf bytes.Buffer
	buf.WriteString("[")
	for i, c := range w.Causers {
		if i > 0 {
			buf.WriteString(", ")
		}
		buf.WriteString(c.Error())
	}
	buf.WriteString("]")
	// buf.WriteString(w.Stack)
	return buf.String()
}

// Format formats the stack of Frames according to the fmt.Formatter interface.
//
//    %s	lists source files for each Frame in the stack
//    %v	lists the source file and line number for each Frame in the stack
//
// Format accepts flags that alter the printing of some verbs, as follows:
//
//    %+v   Prints filename, function, and line number for each Frame in the stack.
func (w *causes) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		if s.Flag('+') {
			_, _ = fmt.Fprintf(s, "%+v", w.Error())
			w.Stack.Format(s, verb)
			return
		}
		fallthrough
	case 's':
		_, _ = io.WriteString(s, w.Error())
	case 'q':
		_, _ = fmt.Fprintf(s, "%q", w.Error())
	}
}

func (w *causes) Cause() error {
	if len(w.Causers) == 0 {
		return nil
	}
	return w.Causers[0]
}

func (w *causes) Causes() []error {
	if len(w.Causers) == 0 {
		return nil
	}
	return w.Causers
}

func (w *causes) Unwrap() error {
	// return w.Cause()

	for _, err := range w.Causers {
		//u, ok := err.(interface {
		//	Unwrap() error
		//})
		//if ok {
		//	return u.Unwrap()
		//}
		return err // just return the first cause
	}
	return nil
}

func (w *causes) Is(target error) bool {
	return IsSlice(w.Causers, target)
}

func (w *causes) TypeIs(target error) bool {
	return TypeIsSlice(w.Causers, target)
}

// As finds the first error in `err`'s chain that matches target, and if so, sets
// target to that error value and returns true.
func (w *causes) As(target interface{}) bool {
	return AsSlice(w.Causers, target)
}
