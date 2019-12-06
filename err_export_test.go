// Copyright © 2019 Hedzr Yeh.

package errors

import (
	"io"
	"testing"
)

func TestIsAsInner(t *testing.T) {
	e := &ExtErr{errs: []error{io.EOF}}
	ext := &ExtErr{inner: e, errs: []error{e, e}}
	if !Is(ext, e) {
		t.Fatalf("wrong! expect ext ==IS== &ExtErr{errs:[]error{io.EOF}}")
	}

	var ase error = e
	if !As(ext, &ase) {
		t.Fatalf("wrong! can't extract &ExtErr{errs: []error{io.EOF}} from ext")
	}

	ext = &ExtErr{errs: []error{&ExtErr{inner: e}, e}}
	if !ext.As(&ase) {
		t.Fatalf("wrong! can't extract &ExtErr{errs: []error{io.EOF}} from ext")
	}
}
