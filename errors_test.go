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
	// err = New("hello").Attach(io.EOF)
	// t.Logf("failed: %+v", err)
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

func TestTemplateFormat(t *testing.T) {
	err := New("cannot set: %v (%v) -> %v (%v)")

	_ = err.FormatWith("a", "b", "c", "d")
	t.Logf("Error: %v", err)
	t.Logf("Error: %+v", err)

	_ = err.FormatWith("1", "2", "3", "4")
	t.Logf("Error: %v", err)
}

func TestContainerMore(t *testing.T) {
	var err error
	ec := New("copying got errors")
	ec.Attach(New("some error"))
	ec.Defer(&err)
	if err == nil {
		t.Fatal(`bad`)
	} else {
		t.Logf(`wanted err is non-nil: %v`, err)
	}
}

func TestIsDeep(t *testing.T) {
	var err error
	ec := New("copying got errors")
	in := New("unmatched %q")
	ec.Attach(in.FormatWith("demo"))
	ec.Defer(&err)
	if err == nil {
		t.Fatal(`bad`)
	} else {
		t.Logf(`wanted err is non-nil: %v`, err)
		if !Is(err, in) {
			t.Fatal("expecting Is() got returning true")
		}
	}
}
