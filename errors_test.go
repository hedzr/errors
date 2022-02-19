package errors

import (
	"io"
	"testing"
)

func TestNew(t *testing.T) {

	err := New("hello %v", "world")

	t.Logf("failed: %+v", err)

	err = Skip(1).WithSkip(0).
		WithMessage("bug skip 0").Build()
	t.Logf("failed: %+v", err)

	err = Message("1").WithSkip(0).
		WithMessage("bug msg").Build()
	t.Logf("failed: %+v", err)

	err = NewBuilder().
		WithCode(Internal).
		WithAttach(io.EOF).
		WithAttach(io.ErrShortWrite).
		Build()
	t.Logf("failed: %+v", err)

}

func TestWithStackInfo_New(t *testing.T) {

	err := New("hello %v", "world")

	err.WithAttach(io.EOF, io.ErrShortWrite).
		WithAttach(io.ErrClosedPipe).
		WithCode(Internal).
		End()
	t.Logf("failed: %+v", err)

}
