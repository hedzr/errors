// Copyright © 2019 Hedzr Yeh.

//+build go1.13

package errors_test

import (
	"fmt"
	"github.com/hedzr/errors"
	"io"
	"os"
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

func TestIsStd(t *testing.T) {
	err2 := fmt.Errorf("BUG %w BUG", os.ErrExist)
	err3 := errors.Unwrap(err2)

	t.Log(err2)
	t.Log(err3)
	if !errors.Is(err3, os.ErrExist) {
		t.Fatal("expect errors.Is(err3, os.ErrExist) returns true")
	}

	t.Log(errors.HasWrappedError(err2))
	t.Log(errors.HasWrappedError(err3))
}

func TestAsStd(t *testing.T) {
	err1 := &os.PathError{Err: os.ErrPermission}
	err2 := fmt.Errorf("BUG %w BUG", err1)
	err3 := errors.Unwrap(err2)

	t.Log(err1)
	t.Log(err2)
	t.Log(err3)
	t.Log(errors.Is(err3, os.ErrExist))

	var perr *os.PathError
	if errors.As(err3, &perr) {
		fmt.Println(perr.Path)
	} else {
		t.Fatal("expect errors.As(err3, &perr) returns true")
	}
}

func TestUnwrapStd(t *testing.T) {
	err1 := errors.New("1")
	err2 := fmt.Errorf("BUG %w BUG", err1)
	err3 := errors.Unwrap(err2)

	t.Log(err1)
	t.Log(err2)
	t.Log(err3)

	if err3 != err1 {
		t.Fatal("expect errors.Unwrap(err2) returns err1")
	}
}
