package errors

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"strings"
)

type causes2 struct {
	Code

	Causers []error
	msg     string

	unwrapIndex  int // simple index for iterating Unwrap
	maxStringLen int // the output string max-length for an object (see also sites/taggedSites), negative or zero means no limit.

	liveArgs []interface{} //nolint:revive // error message template ?
}

func (w *causes2) limitObj(obj interface{}) (s string) { //nolint:revive
	s = fmt.Sprintf("%+v", obj)
	if w.maxStringLen > 0 && len(s) > w.maxStringLen {
		s = s[0:w.maxStringLen-3] + "..."
	}
	return
}

// WithMaxObjectStringLength for error interface
func (w *causes2) WithMaxObjectStringLength(maxlen int) *causes2 {
	w.maxStringLen = maxlen
	return w
}

// WithCode for error interface
func (w *causes2) WithCode(code Code) *causes2 {
	w.Code = code
	return w
}

func (w *causes2) WithMessage(message string, args ...interface{}) *causes2 { //nolint:revive
	if len(args) > 0 {
		message = fmt.Sprintf(message, args...) //nolint:revive
	}
	w.msg = message
	return w
}

// End ends the WithXXX stream calls while you dislike unwanted `err =`.
func (w *causes2) End() {}

// Defer can be used as a defer function to simplify your codes.
//
// The codes:
//
//	 func some(){
//	   // as a inner errors container
//	   child := func() (err error) {
//	  	errContainer := errors.New("")
//	  	defer errContainer.Defer(&err)
//
//	  	for _, r := range []error{io.EOF, io.ErrClosedPipe, errors.Internal} {
//	  		errContainer.Attach(r)
//	  	}
//
//	  	return
//	   }
//
//	   err := child()
//	   t.Logf("failed: %+v", err)
//	}
func (w *causes2) Defer(err *error) { //nolint:gocritic
	if w.IsEmpty() {
		*err = nil
	} else {
		*err = w
	}
}

func (w *causes2) Clear() Container {
	w.msg = ""
	w.Causers = nil
	w.liveArgs = nil
	w.Code = OK
	w.unwrapIndex = 0
	w.maxStringLen = 0
	return w
}

// WithErrors appends errs
//
// WithStackInfo.Attach() can only wrap and hold one child error object.
//
// WithErrors attach child errors into an error container.
// For a container which has IsEmpty() interface, it would not
// be attached if it is empty (i.e. no errors).
//
// For a nil error object, it will be ignored.
func (w *causes2) WithErrors(errs ...error) *causes2 { //nolint:revive
	for _, e := range errs {
		if e != nil {
			if check, ok := e.(interface{ IsEmpty() bool }); ok {
				if !check.IsEmpty() {
					w.Causers = append(w.Causers, e)
				} else if e.Error() != "" {
					w.Causers = append(w.Causers, e)
				}
			} else {
				w.Causers = append(w.Causers, e)
			}
		}
	}
	return w
}

// Attach collects the errors except it's nil
func (w *causes2) Attach(errs ...error) {
	// _ = w.WithErrors(errs...)

	for _, e := range errs {
		if e != nil {
			if x, ok := e.(*causes2); ok && x == w {
				continue
			}
			w.Causers = append(w.Causers, e)
		}
	}
}

// FormatWith _
func (w *causes2) FormatWith(args ...interface{}) error { //nolint:revive
	c := w.Clone()
	c.liveArgs = args
	return c
}

// Clone _
func (w *causes2) Clone() *causes2 {
	c := &causes2{
		Code:        w.Code,
		Causers:     w.Causers,
		msg:         w.msg,
		unwrapIndex: w.unwrapIndex,
		liveArgs:    w.liveArgs,
	}
	return c
}

func (w *causes2) Error() string {
	return w.makeErrorString(false)
}

func (w *causes2) makeErrorString(line bool) string { //nolint:revive
	var buf bytes.Buffer
	if w.msg != "" {
		if len(w.liveArgs) > 0 {
			msg := fmt.Sprintf(w.msg, w.liveArgs...)
			_, _ = buf.WriteString(msg)
		} else {
			_, _ = buf.WriteString(w.msg)
		}
	}

	if line {
		_, _ = buf.WriteRune('\n')
		if w.Code != OK {
			_, _ = buf.WriteString(w.Code.String())
			_, _ = buf.WriteRune('\n')
		}

		for _, c := range w.Causers {
			_, _ = buf.WriteString("  - ")
			var xc *causes2
			if As(c, &xc) {
				_, _ = buf.WriteString(leftPad(xc.makeErrorString(line), "  ", false))
			} else {
				_, _ = buf.WriteString(leftPad(c.Error(), "    ", false))
			}
		}

		// buf.WriteRune('\n')
		return buf.String()
	}

	var needclose, needsep bool
	if w.Code != OK {
		if buf.Len() > 0 {
			_, _ = buf.WriteRune(' ')
		}
		_, _ = buf.WriteString("[")
		_, _ = buf.WriteString(w.Code.String())
		needclose = true
		needsep = true
	}
	if w.msg == "" {
		needsep = false
	}
	if len(w.Causers) > 0 {
		if buf.Len() > 0 {
			_, _ = buf.WriteRune(' ')
		}
		_, _ = buf.WriteString("[")
		needclose = true
	}

	for i, c := range w.Causers {
		if i > 0 || needsep {
			_, _ = buf.WriteString(" | ")
		}
		_, _ = buf.WriteString(c.Error())
	}
	if needclose {
		_, _ = buf.WriteString("]")
	}
	// buf.WriteString(w.Stack)
	return buf.String()
}

func leftPad(s, padStr string, firstLine bool) string { //nolint:revive
	if padStr == "" {
		return s
	}

	var ln int
	var sb strings.Builder
	scanner := bufio.NewScanner(bufio.NewReader(strings.NewReader(s)))
	for scanner.Scan() {
		if ln != 0 || firstLine {
			_, _ = sb.WriteString(padStr)
		}
		_, _ = sb.WriteString(scanner.Text())
		_, _ = sb.WriteRune('\n')
		ln++
	}
	return sb.String()
}

// Format formats the stack of Frames according to the fmt.Formatter interface.
//
//	%s	lists source files for each Frame in the stack
//	%v	lists the source file and line number for each Frame in the stack
//
// Format accepts flags that alter the printing of some verbs, as follows:
//
//	%+v   Prints filename, function, and line number for each Frame in the stack.
func (w *causes2) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		if s.Flag('+') {
			_, _ = fmt.Fprintf(s, "%+v", w.makeErrorString(true))
			return
		}
		fallthrough
	case 's':
		_, _ = io.WriteString(s, w.Error())
	case 'q':
		_, _ = fmt.Fprintf(s, "%q", w.Error())
	}
}

// String for stringer interface
func (w *causes2) String() string { return w.Error() }

func (w *causes2) Cause() error {
	if len(w.Causers) == 0 {
		return nil
	}
	return w.Causers[0]
}

// Causes simply returns the wrapped inner errors.
// It doesn't consider an wrapped Code entity is an inner error too.
// So if you wanna to extract any inner error objects, use
// errors.Unwrap for instead. The errors.Unwrap could extract all
// of them one by one:
//
//	var err = errors.New("hello").WithErrors(io.EOF, io.ShortBuffers)
//	var e error = err
//	for e != nil {
//	    e = errors.Unwrap(err)
//	}
func (w *causes2) Causes() []error {
	if len(w.Causers) == 0 {
		return nil
	}
	return w.Causers
}

func (w *causes2) Unwrap() error {
	defer func() { w.unwrapIndex++ }()

	if w.unwrapIndex >= 0 && w.unwrapIndex < len(w.Causers) {
		return w.Causers[w.unwrapIndex]
	}
	if w.unwrapIndex == len(w.Causers) {
		if w.Code != OK {
			return w.Code
		}
	}

	// reset index
	w.unwrapIndex = -1
	return nil
}

func (w *causes2) Reset() {
	w.unwrapIndex = 0
}

// IsDescended test if ancestor is an error template and descendant
// is derived from it by calling ancestor.FormatWith.
func IsDescended(ancestor, descendant error) bool {
	if a, ok := ancestor.(interface{ IsDescended(descendant error) bool }); ok {
		return a.IsDescended(descendant)
	}
	return false
}

func (w *causes2) IsDescended(descendant error) bool {
	if e, ok := descendant.(*causes2); ok {
		return e.Code == w.Code && e.msg == w.msg
	}
	return false
}

func (w *causes2) Is(target error) bool {
	if w.Code != OK {
		if c, ok := target.(Code); ok && c == w.Code {
			return true
		}
	}
	if IsSlice(w.Causers, target) {
		return true
	}
	if te, ok := target.(*causes2); ok {
		return w.equal(te)
	}
	return false
}

func (w *causes2) equal(target *causes2) bool {
	return w.Code == target.Code &&
		w.msg == target.msg &&
		w.arrayEqual(w.Causers, target.Causers)
}

func (w *causes2) arrayEqual(a, b []error) bool {
	if len(a) != len(b) {
		return false
	}
	yes := true
	for i := 0; i < len(a); i++ {
		if Is(a[i], b[i]) {
			yes = false
			break
		}
	}
	return yes
}

func (w *causes2) TypeIs(target error) bool {
	return TypeIsSlice(w.Causers, target)
}

// As finds the first error in `err`'s chain that matches target,
// and if so, sets target to that error value and returns true.
func (w *causes2) As(target interface{}) bool { //nolint:revive
	if c, ok := target.(*Code); ok {
		*c = w.Code
		return true
	}
	if c, ok := target.(*[]error); ok {
		*c = w.Causers
		return true
	}
	if c, ok := target.(**causes2); ok {
		*c = w
		return true
	}
	return AsSlice(w.Causers, target)
}

// IsEmpty tests has attached errors
func (w *causes2) IsEmpty() bool {
	return w.Code == OK && len(w.Causers) == 0 && len(w.liveArgs) == 0
}
