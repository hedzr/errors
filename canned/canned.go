package canned

import (
	"bytes"
	"fmt"
	"io"
)

// New returns a canned error object, later you may `Attach` a group of errors into it.
func New(format string, args ...interface{}) *Canned {
	var str string
	if len(args) > 0 {
		str = fmt.Sprintf(format, args...)
	} else {
		str = format
	}
	return &Canned{
		msg:   str,
		stack: callers(),
	}
}

// Canned is a holder to store a group of inner errors.
type Canned struct {
	msg   string
	errs  []error
	stack *Stack
}

// Attach could add an error into the Canned container
func (e *Canned) Attach(err error) {
	e.errs = append(e.errs, err)
}

// IsEmpty returns false if any inner errors exists
func (e *Canned) IsEmpty() bool {
	return len(e.errs) == 0
}

// Format implements Formatter interface for fmt.Printf("%+v", err)
func (e *Canned) Format(st fmt.State, verb rune) {
	switch verb {
	case 'v':
		if st.Flag('+') {
			io.WriteString(st, e.Error())
			e.stack.Format(st, verb)
			return
		}
		fallthrough
	case 's':
		io.WriteString(st, e.Error())
	case 'q':
		fmt.Fprintf(st, "%q", e.Error())
	}
}

// Error returns error message string presentation
func (e *Canned) Error() string {
	var buf bytes.Buffer

	buf.WriteString(e.msg)

	if len(e.errs) > 0 {
		buf.WriteString("[")
	}
	for i, ee := range e.errs {
		if i > 0 {
			buf.WriteString(", ")
		}
		buf.WriteString(ee.Error())
	}
	if len(e.errs) > 0 {
		buf.WriteString("]")
	}

	return buf.String()
}
