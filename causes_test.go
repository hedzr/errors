package errors

import (
	"io"
	"testing"
)

func TestCauses2_New(t *testing.T) {
	err := &causes2{
		Code:    Internal,
		Causers: nil,
		msg:     "ui err",
	}

	t.Logf("failed: %+v", err)

	err = &causes2{
		Code:    0,
		Causers: nil,
		msg:     "ui err",
	}
	_ = err.WithErrors(io.EOF, io.ErrClosedPipe)
	t.Logf("failed: %+v", err)

	err = &causes2{
		Code:    Internal,
		Causers: nil,
		msg:     "ui err",
	}
	err.WithErrors(io.EOF, io.ErrClosedPipe).End()
	t.Logf("failed: %+v", err)

}

func TestAsAContainer(t *testing.T) {
	// as a inner errors container
	child := func() (err error) {
		errContainer := New("")

		defer errContainer.Defer(&err)
		for _, r := range []error{io.EOF, io.ErrShortWrite, io.ErrClosedPipe, Internal} {
			errContainer.Attach(r)
		}

		return
	}

	err := child()
	t.Logf("failed: %+v", err)
}
