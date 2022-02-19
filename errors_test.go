package errors

import (
	"io"
	"testing"
)

func TestNew(t *testing.T) {

	err := New("hello %v", "world")

	t.Logf("failed: %+v", err)

	err = Skip(1).
		WithSkip(0).
		WithMessage("bug skip 0").
		Build()
	t.Logf("failed: %+v", err)

	err = Message("1").
		WithSkip(0).
		WithMessage("bug msg").
		Build()
	t.Logf("failed: %+v", err)

	err = NewBuilder().
		WithCode(Internal).
		WithErrors(io.EOF).
		WithErrors(io.ErrShortWrite).
		Build()
	t.Logf("failed: %+v", err)

	err = New(WithErrors(io.EOF, io.ErrShortWrite))
	t.Logf("failed: %+v", err)

	err = New()
	t.Logf("failed: %+v", err)

	err = New("hello").Attach(io.EOF)
	t.Logf("failed: %+v", err)

}

func TestWithStackInfo_New(t *testing.T) {

	err := New("hello %v", "world")

	err.WithErrors(io.EOF, io.ErrShortWrite).
		WithErrors(io.ErrClosedPipe).
		WithCode(Internal).
		End()
	t.Logf("failed: %+v", err)

}
