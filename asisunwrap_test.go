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

func TestWrap(t *testing.T) {
	err := Wrap(nil, "hello, %v", "world")
	t.Logf("failed: %+v", err)

	err = Wrap(Internal, "hello, %v", "world")
	t.Logf("failed: %+v", err)

	err = Wrap(Internal, "hello, %v", "world")
	t.Logf("failed: %+v", err)

	Unwrap(io.EOF)
}

func TestTypeIsSlice(t *testing.T) {
	TypeIsSlice(nil, nil)
	TypeIsSlice(nil, io.EOF)

	TypeIsSlice([]error{io.EOF}, io.EOF)

	err := Wrap(Internal, "hello, %v", "world")
	err.WithErrors(NotFound).End()

	if TypeIsSlice(err.Causes(), NotFound) == false {
		t.Fatalf("not ok")
	}

	if TypeIs(err, NotFound) == false {
		t.Fatalf("not ok")
	}

	err2 := New().WithErrors(NotFound, err)
	err3 := New().WithErrors(NotFound, err2)
	if TypeIsSlice(Causes(err3), err) == false {
		t.Fatalf("not ok")
	}

	if TypeIs(err2, NotFound) == false {
		t.Fatalf("not ok")
	}
	if TypeIs(err3, NotFound) == false {
		t.Fatalf("not ok")
	}
	if TypeIs(err3, err2) == false {
		t.Fatalf("not ok")
	}
	if TypeIs(err3, err) == false {
		t.Fatalf("not ok")
	}
	TypeIs(err3, nil)

	IsSlice(Causes(err3), nil)
	IsSlice(Causes(err3), io.ErrShortBuffer)
	IsSlice(Causes(err3), io.EOF)
	IsSlice(Causes(err3), NotFound)
	IsSlice(Causes(err3), Internal)
	IsSlice(Causes(err3), err2)
	IsSlice(Causes(err3), err)

	Is(err3, nil)
	Is(err3, io.ErrShortBuffer)
	Is(err3, io.EOF)
	Is(err3, NotFound)
	Is(err3, Internal)
	Is(err3, err2)
	Is(err3, err)
}

func TestAsRaisePanic(t *testing.T) {

	t.Run("1", func(t *testing.T) {
		defer func() { recover() }()
		As(nil, nil)
	})

	t.Run("2", func(t *testing.T) {
		defer func() { recover() }()
		var v int
		As(nil, &v)
	})

	t.Run("3", func(t *testing.T) {
		defer func() { recover() }()
		var err error
		As(nil, &err)
	})

	t.Run("4", func(t *testing.T) {
		defer func() { recover() }()
		var err int
		As(nil, err)
	})

}

func TestAsSliceRaisePanic(t *testing.T) {

	t.Run("1", func(t *testing.T) {
		defer func() { recover() }()
		AsSlice(nil, nil)
	})

	t.Run("2", func(t *testing.T) {
		defer func() { recover() }()
		var v int
		AsSlice(nil, &v)
	})

	t.Run("3", func(t *testing.T) {
		defer func() { recover() }()
		var err error
		AsSlice(nil, &err)
	})

	t.Run("4", func(t *testing.T) {
		defer func() { recover() }()
		var err int
		AsSlice(nil, err)
	})

}
