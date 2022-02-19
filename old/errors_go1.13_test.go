// Copyright © 2020 Hedzr Yeh.

//go:build go1.13
// +build go1.13

package old

import (
	errorsstd "errors"
	"fmt"

	// pkgerrors "github.com/pkg/errors"
	"io"
	"os"
	"testing"
)

func geneof13() error {
	return fmt.Errorf("xxx %w wrapped at go1.%v+", io.EOF, 13)
}

func Test1go113(t *testing.T) {
	var err error

	err = geneof()
	if Cause(err) == io.EOF {
		t.Logf("ok: %v", err)
	} else {
		t.Fatal("expect it is a EOF")
	}

	//err = geneof2()
	//if pkgerrors.Cause(err) == io.EOF {
	//	t.Logf("ok: %v", err)
	//} else {
	//	t.Fatal("expect it is a EOF")
	//}

	err = geneof13()
	if errorsstd.Is(err, io.EOF) {
		t.Logf("ok: %v", err)
	} else {
		t.Fatal("expect it is a EOF")
	}
}

func Test2go113(t *testing.T) {
	var err error

	err = geneofx()
	if Cause(err) == io.EOF {
		t.Logf("ok: %v", err)
	} else {
		t.Fatal("expect it is a EOF")
	}
	if Cause(err) == io.EOF {
		t.Logf("ok: %+v", err)
	} else {
		t.Fatal("expect it is a EOF")
	}

	err = geneofxs()
	if Cause(err) == io.EOF {
		t.Logf("ok: %v", err)
	} else {
		t.Fatal("expect it is a EOF")
	}

	// errorx tests -------------------------------

	// errorx.Cause() and Cause1()
	if Cause(err) == io.EOF {
		// Wrap(err, msg): the error object has stacktrace info
		t.Logf("ok: %+v", err)
	} else {
		t.Fatal("Cause() failed: expect it is a EOF")
	}
	if Is(err, io.EOF) {
		// Wrap(err, msg): the error object has stacktrace info
		t.Logf("ok: %+v", err)
	} else {
		t.Fatal("Is() failed: expect it is a EOF")
	}
	if Unwrap(Unwrap(err)) == io.EOF {
		// Wrap(err, msg): the error object has stacktrace info
		t.Logf("ok: %+v", err)
	} else {
		t.Fatal("Unwrap() failed: expect it is a EOF")
	}

	var perr *os.PathError
	err = Wrap(&os.PathError{Err: io.EOF, Op: "find", Path: "/"}, "wrong path and rights")
	if As(err, &perr) {
		t.Logf("ok: %+v", *perr)
	} else {
		t.Fatal("As() failed: expect it is a os.PathError{}")
	}

	// var c = NewContainer("container")
	// AttachTo(c, io.EOF, io.ErrShortBuffer, io.ErrUnexpectedEOF)
	// t.Logf("ok: %+v | container is empty: %v", c, ContainerIsEmpty(c))
	// if ContainerIsEmpty(c) != false {
	// 	t.Fatal("ContainerIsEmpty(c) failed: expect it is false.")
	// }
}
