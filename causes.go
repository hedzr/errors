// Copyright © 2020 Hedzr Yeh.

package errors

import (
	"bytes"
	"fmt"
	"io"
)

type causes struct {
	Causers []error
	*stack
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
	// buf.WriteString(w.stack)
	return buf.String()
}

func (w *causes) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		if s.Flag('+') {
			fmt.Fprintf(s, "%+v", w.Error())
			w.stack.Format(s, verb)
			return
		}
		fallthrough
	case 's':
		io.WriteString(s, w.Error())
	case 'q':
		fmt.Fprintf(s, "%q", w.Error())
	}
}
