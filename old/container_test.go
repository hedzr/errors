// Copyright © 2020 Hedzr Yeh.

package old_test

import (
	"bufio"
	"bytes"
	"gopkg.in/hedzr/errors.v2/old"
	"io"
	"testing"
)

func sampleC(simulate bool) (err error) {
	c := old.NewContainer("sample error")
	defer c.Defer(&err)
	if simulate {
		old.AttachTo(c, io.EOF, io.ErrUnexpectedEOF, io.ErrShortBuffer, io.ErrShortWrite)
	}
	err = c.Error()
	return
}

func TestContainer(t *testing.T) {
	err := sampleC(false)
	if err != nil {
		t.Fatal(err)
	} else {
		t.Logf("%+v", err)
	}

	err = sampleC(true)
	if err == nil {
		t.Fatal("want error")
	} else {
		t.Logf("%+v", err)
	}
}

type bizStrut struct {
	err old.Holder
	w   *bufio.Writer
}

func (bw *bizStrut) Write(b []byte) {
	_, err := bw.w.Write(b)
	bw.err.Attach(err)
}

func (bw *bizStrut) Flush() error {
	err := bw.w.Flush()
	bw.err.Attach(err)
	return bw.err.Error()
}

func TestContainer2(t *testing.T) {
	var bb bytes.Buffer
	var bw = &bizStrut{
		err: old.NewContainer("bizStruct have errors %v", "ext"),
		w:   bufio.NewWriter(&bb),
	}
	bw.Write([]byte("hello "))
	bw.Write([]byte("world "))
	if err := bw.Flush(); err != nil {
		t.Fatal(err)
	}
	if !bw.err.IsEmpty() {
		t.Fatal("non-empty container here")
	}
}

func TestContainer3(t *testing.T) {
	err := sampleC(true)
	if old.ContainerIsEmpty(err) {
		t.Fatal("non-empty container here")
	}

	c := old.New("sample error")
	if old.ContainerIsEmpty(c) {
		t.Fatal("non-empty container here")
	}
}
