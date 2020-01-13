package canned

import (
	"io"
	"testing"
)

func TestPkgWrap(t *testing.T) {
	err := New("canned errors")
	for _, e := range []struct{ err error }{
		//
	} {
		err.Attach(e.err)
	}
	if !err.IsEmpty() {
		t.Fatal("expect empty")
	}

	err = New("canned errors #%d", 1)
	for _, e := range []struct{ err error }{
		{io.EOF},
		{io.ErrShortWrite},
		{io.ErrClosedPipe},
	} {
		err.Attach(e.err)
	}
	if err.IsEmpty() {
		t.Fatal("expect not empty")
	}
	t.Logf("%+v", err)
	t.Logf("%v", err)
	t.Logf("%s", err)
	t.Logf("%q", err)
}
