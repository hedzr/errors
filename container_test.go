// Copyright © 2020 Hedzr Yeh.

package errors_test

import (
	"gopkg.in/hedzr/errors.v2"
	"io"
	"testing"
)

func sample(simulate bool) (err error) {
	c := errors.NewContainer("sample error")
	if simulate {
		errors.AttachTo(c, io.EOF, io.ErrUnexpectedEOF, io.ErrShortBuffer, io.ErrShortWrite)
	}
	err = c.Error()
	return
}

func TestContainer(t *testing.T) {
	err := sample(false)
	if err != nil {
		t.Fatal(err)
	} else {
		t.Logf("%+v", err)
	}

	err = sample(true)
	if err == nil {
		t.Fatal("want error")
	} else {
		t.Logf("%+v", err)
	}
}
