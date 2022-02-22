package errors

import (
	"io"
	"os"
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

	err = NewBuilder().
		WithCode(Internal).
		WithErrors(io.EOF).
		WithErrors(io.ErrShortWrite).
		Build()
	// Attach(io.ErrClosedPipe)
	t.Logf("failed: %+v", err)

	err = NewBuilder().Build()
	err.Attach(os.ErrClosed, os.ErrInvalid, os.ErrPermission)
	t.Logf("failed: %+v", err)

	err = New(WithErrors(io.EOF, io.ErrShortWrite))
	t.Logf("failed: %+v", err)

	err = New()
	t.Logf("failed: %+v", err)

	// since v3.0.5, Attach() has no return value
	//err = New("hello").Attach(io.EOF)
	//t.Logf("failed: %+v", err)

}

func TestUnwrap(t *testing.T) {
	t.Log("unwrap all inner error (including Code object) one by one:")
	err := New()
	err. // WithCode(NotFound).
		WithErrors(io.EOF, io.ErrShortBuffer).
		WithMessage("has code and errors").
		Attach(os.ErrClosed, os.ErrInvalid, os.ErrPermission)
	var e error = err
	for e != nil {
		e = Unwrap(err)
		t.Logf("failed: %v", e)
	}
	if o, ok := err.(interface{ Reset() }); ok {
		o.Reset()
	}
	t.Log("again")
	e = err
	for e != nil {
		e = Unwrap(err)
		t.Logf("failed: %v", e)
	}
}
func TestWithStackInfo_New(t *testing.T) {

	err := New("hello %v", "world")

	err.WithErrors(io.EOF, io.ErrShortWrite).
		WithErrors(io.ErrClosedPipe).
		WithCode(Internal).
		End()
	t.Logf("failed: %+v", err)

	if CanCause(err) {
		e := err.(causer).Cause()
		t.Logf("failed: %v", e)
	}
	if CanCauses(err) {
		e := Causes(err)
		t.Logf("failed: %v", e)
	}

	err2 := New("hello %v", "world")
	err2.WithData(9, err, 10).WithTaggedData(TaggedData{"9": 9, "10": 10}).End()

	t.Logf("failed: %+v", err2)
	t.Logf("Data: %v", err2.Data())
	t.Logf("TaggedData: %v", err2.TaggedData())
}
