package errors

import (
	"io"
	"testing"
)

func TestIs(t *testing.T) {

	series := []error{io.EOF, io.ErrShortWrite, io.ErrClosedPipe, Internal}

	err := &causes2{
		Code:    0,
		Causers: nil,
		msg:     "ui err",
	}
	_ = err.WithErrors(io.EOF, io.ErrClosedPipe)

	for _, e := range []error{io.EOF, io.ErrClosedPipe} {
		if !Is(err, e) {
			t.Fatalf("test for Is(%v) failed", e)
		}
	}

	err2 := New("hello %v", "world")

	// the old err2 (i.e. err3) will be moved into err2's slice
	// container, and more errors (io.EOF, io.ErrShortWrite, and
	// io.ErrClosedPipe) will be appended into the slice
	// container.
	err2.WithErrors(io.EOF, io.ErrShortWrite).
		WithErrors(io.ErrClosedPipe).
		WithCode(Internal).
		End()
	for _, e := range series {
		if !Is(err2, e) {
			t.Fatalf("test for Is(%v) failed", e)
		}
	}

	t.Logf("failed: %+v", err)

	var code Code
	if !(As(err2, &code) && code == Internal) {
		t.Fatalf("cannot extract coded error with As()")
	}

	// so As() will extract the first element in err2's slice container,
	// that is err3.
	var ee1 error
	if !(As(err2, &ee1) && ee1 == io.EOF) {
		t.Fatalf("cannot extract 'hello world' error with As(), ee1 = %v", ee1)
	}

	var ee2 []error
	if !(As(err2, &ee2) && len(ee2) == 3) {
		t.Fatalf("cannot extract []error error with As(), ee2 = %v", ee2)
	}

	var index int
	for ; ee1 != nil; index++ {
		ee1 = Unwrap(err2)
		if ee1 != nil && ee1 != series[index] {
			t.Fatalf("%d. cannot extract '%v' error with As(), ee1 = %v", index, series[index], ee1)
		}
	}
}
