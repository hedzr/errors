// Copyright © 2020 Hedzr Yeh.

package errors

import (
	"bytes"
	"fmt"
	"io"
	"reflect"
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
		u, ok := err.(interface {
			Unwrap() error
		})
		if ok {
			return u.Unwrap()
		}
	}
	return nil
}

func (w *causes) Is(target error) bool {
	if target == nil {
		for _, e := range w.Causers {
			if e == target {
				return true
			}
		}
		return false
	}

	isComparable := reflect.TypeOf(target).Comparable()
	for {
		if isComparable {
			for _, e := range w.Causers {
				if e == target {
					return true
				}
			}
			return false
		}

		for _, e := range w.Causers {
			if x, ok := e.(interface{ Is(error) bool }); ok && x.Is(target) {
				return true
			}
			if err := Unwrap(e); err == nil {
				return false
			}
		}
		return false
	}
}
