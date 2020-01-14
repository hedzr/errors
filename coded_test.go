package errors

import (
	"io"
	"testing"
)

func TestCodes(t *testing.T) {
	err := InvalidArgument.New("wrong").Attach(io.ErrShortWrite)
	t.Log(err)
	t.Logf("%+v", err)

	if !Is(err, io.ErrShortWrite) {
		t.Fatal("wrong Is()")
	}
	if Is(err, io.EOF) {
		t.Fatal("wrong Is()")
	}
}

func TestCodesEqual(t *testing.T) {
	err := InvalidArgument.New("wrong").Attach(io.ErrShortWrite)

	ok := EqualR(err, InvalidArgument)
	if !ok {
		t.Fatal("want Equal() return true but got false")
	}
}
