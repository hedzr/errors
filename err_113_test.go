// Copyright © 2019 Hedzr Yeh.

package errors_test

import (
	"fmt"
	"github.com/hedzr/errors"
	"io"
	"testing"
)

func TestIsAs113(t *testing.T) {
	var err error
	err = errors.New("something").Attach(errBug1, errBug2).Nest(errBug1, errBug2).Msg("anything")
	if !errors.Is(err, errBug1) {
		t.Fatalf("wrong! expect err ==IS== errBug1")
	}
	if !errors.Is(err, errBug2) {
		t.Fatalf("wrong! expect err ==IS== errBug2")
	}

	err2 := errors.NewWithError(io.ErrShortWrite).Nest(io.EOF)
	if !errors.Is(err2, io.ErrShortWrite) {
		t.Fatalf("wrong! expect err2 ==IS== io.ErrShortWrite")
	}
	if !errors.Is(err2, io.EOF) {
		t.Fatalf("wrong! expect err2 ==IS== io.EOF")
	}

	var ase error = io.EOF
	err = errors.New("x").Attach(fmt.Errorf("error is %w", io.EOF), errors.NewWithError(io.EOF))
	if !errors.Is(err, io.EOF) {
		t.Fatalf("wrong! expect err ==IS== io.EOF")
	}
	if !errors.As(err, &ase) {
		t.Fatalf("wrong! can't extract EOF from err")
	}

	ex := errors.NewWithError(io.EOF)
	err = errors.New("x").Nest(fmt.Errorf("error is %w", io.EOF), ex)
	if !errors.Is(err, io.EOF) {
		t.Fatalf("wrong! expect err ==IS== io.EOF")
	}
	if !errors.Is(err, ex) {
		t.Fatalf("wrong! expect err ==IS== ExtErr{io.EOF}")
	}
	if !errors.As(err, &ase) {
		t.Fatalf("wrong! can't extract EOF from err")
	}

}

