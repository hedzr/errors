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
