package errors

import (
	// errors2 "errors"
	"io"
	"strconv"
	"testing"
)

func TestErrorCodeIs(t *testing.T) {
	var err error = BadRequest
	if !Is(err, BadRequest) {
		t.Fatalf("want is")
	}

	err = io.ErrClosedPipe
	if Is(err, BadRequest) {
		t.Fatalf("want not is")
	}

	err = NotFound
	if Is(err, BadRequest) {
		t.Fatalf("want not is (code)")
	}

	//

	_, err = strconv.ParseInt("hello", 10, 64)
	if !Is(err, strconv.ErrRange) {
		t.Logf("'%v' recoganized OK.", err)
	} else {
		t.Fatalf("'%+v' CANNOT be recoganized", err)
	}

	// for v1.11, Is(err, strconv.ErrSyntax) should be false
	// and for v1.12+, Is(err, strconv.ErrSyntax) must be true
}

// TestErrorsIs _
func TestErrorsIs(t *testing.T) {
	_, err := strconv.ParseFloat("hello", 64)
	t.Logf("err = %+v", err)

	// e1 := errors2.Unwrap(err)
	// t.Logf("e1 = %+v", e1)

	t.Logf("errors.Is(err, strconv.ErrSyntax): %v", Is(err, strconv.ErrSyntax))
	t.Logf("errors.Is(err, &strconv.NumError{}): %v", Is(err, &strconv.NumError{Err: strconv.ErrRange}))

	var e2 *strconv.NumError
	if As(err, &e2) {
		t.Logf("As() ok, e2 = %v", e2)
	} else {
		t.Logf("As() not ok")
	}
}

func TestIs(t *testing.T) {
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

	t.Logf("the error is: %+v", err)
}

func TestIss(t *testing.T) {
	err0 := &causes2{
		Code:    0,
		Causers: nil,
		msg:     "ui err",
	}
	err := err0.WithErrors(io.EOF, io.ErrClosedPipe)
	if !Iss(err, io.EOF) {
		t.Fatal(err)
	}
	if Iss(err, io.ErrNoProgress, io.ErrShortWrite) {
		t.Fatal(err)
	}
	if !Iss(err, io.ErrNoProgress, io.ErrShortWrite, err0) {
		t.Fatal(err)
	}
	t.Logf("err0 = %v, err = %v", err0, err)
}

func TestIs2(t *testing.T) { //nolint:revive
	series := []error{io.EOF, io.ErrShortWrite, io.ErrClosedPipe, Internal}

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

	var code Code
	if !(As(err2, &code) && code == Internal) { //nolint:revive
		t.Fatalf("cannot extract coded error with As()")
	}

	// so As() will extract the first element in err2's slice container,
	// that is err3.
	var ee1 error
	if !(As(err2, &ee1) && ee1 == io.EOF) { //nolint:revive
		t.Fatalf("cannot extract 'hello world' error with As(), ee1 = %v", ee1)
	}

	var ee2 []error
	if !(As(err2, &ee2) && len(ee2) == 3) {
		t.Fatalf("cannot extract []error error with As(), ee2 = %v (len=%v)", ee2, len(ee2))
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

	_ = Unwrap(io.EOF)
}

func TestTypeIsSlice(t *testing.T) {
	TypeIsSlice(nil, nil)
	TypeIsSlice(nil, io.EOF)

	TypeIsSlice([]error{io.EOF}, io.EOF)

	err := Wrap(Internal, "hello, %v", "world")
	err.WithErrors(NotFound).End()

	if !TypeIsSlice(err.Causes(), NotFound) {
		t.Fatalf("not ok")
	}

	if !TypeIs(err, NotFound) {
		t.Fatalf("not ok")
	}

	err2 := New().WithErrors(NotFound, err)
	err3 := New().WithErrors(NotFound, err2)
	if !TypeIsSlice(Causes(err3), err) {
		t.Fatalf("not ok")
	}

	if !TypeIs(err2, NotFound) {
		t.Fatalf("not ok")
	}
	if !TypeIs(err3, NotFound) {
		t.Fatalf("not ok")
	}
	if !TypeIs(err3, err2) {
		t.Fatalf("not ok")
	}
	if !TypeIs(err3, err) {
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
		defer func() { recover() }() //nolint:errcheck //no need
		As(nil, nil)
	})

	t.Run("2", func(t *testing.T) {
		defer func() { recover() }() //nolint:errcheck //no need
		var v int
		As(nil, &v)
	})

	t.Run("3", func(t *testing.T) {
		defer func() { recover() }() //nolint:errcheck //no need
		var err error
		As(nil, &err)
	})

	t.Run("4", func(t *testing.T) {
		defer func() { recover() }() //nolint:errcheck //no need
		var err int
		As(nil, err)
	})
}

func TestAsSliceRaisePanic(t *testing.T) {
	t.Run("1", func(t *testing.T) {
		defer func() { recover() }() //nolint:errcheck //no need
		AsSlice(nil, nil)
	})

	t.Run("2", func(t *testing.T) {
		defer func() { recover() }() //nolint:errcheck //no need
		var v int
		AsSlice(nil, &v)
	})

	t.Run("3", func(t *testing.T) {
		defer func() { recover() }() //nolint:errcheck //no need
		var err error
		AsSlice(nil, &err)
	})

	t.Run("4", func(t *testing.T) {
		defer func() { recover() }() //nolint:errcheck //no need
		var err int
		AsSlice(nil, err)
	})
}
