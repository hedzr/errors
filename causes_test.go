package errors

import (
	"errors"
	"fmt"
	"io"
	"testing"
)

type DivisionError struct {
	IntA int
	IntB int
	Msg  string
}

func (e *DivisionError) Error() string {
	return e.Msg
}

func Divide(a, b int) (int, error) {
	if b == 0 {
		return 0, &DivisionError{
			Msg:  fmt.Sprintf("cannot divide '%d' by zero", a),
			IntA: a, IntB: b,
		}
	}
	return a / b, nil
}

func dummy(t *testing.T) error {
	a, b := 10, 0
	result, err := Divide(a, b)
	if err != nil {
		var divErr *DivisionError
		switch {
		case errors.As(err, &divErr):
			fmt.Printf("%d / %d is not mathematically valid: %s\n",
				divErr.IntA, divErr.IntB, divErr.Error())
		default:
			fmt.Printf("unexpected division error: %s\n", err)
			t.Fail()
		}
		return err
	}

	fmt.Printf("%d / %d = %d\n", a, b, result)
	return err
}

func dummyV3(t *testing.T) error {
	a, b := 10, 0
	result, err := Divide(a, b)
	if err != nil {
		var divErr *DivisionError
		switch {
		case As(err, &divErr):
			fmt.Printf("%d / %d is not mathematically valid: %s\n",
				divErr.IntA, divErr.IntB, divErr.Error())
		default:
			fmt.Printf("unexpected division error: %s\n", err)
			t.Fail()
		}
		return err
	}

	fmt.Printf("%d / %d = %d\n", a, b, result)
	return err
}

func TestCauses2_errors(t *testing.T) {
	err := io.EOF

	if !errors.Is(err, io.EOF) {
		t.Fail()
	}

	err = dummy(t)
	err = fmt.Errorf("wrapped: %w", err)
	t.Logf("divide: %v", err)
	t.Logf("Unwrap: %v", errors.Unwrap(err))
}

func TestCauses2_errorsV3(t *testing.T) {
	err := io.EOF

	if !Is(err, io.EOF) {
		t.Fail()
	}

	err = dummyV3(t)
	err = fmt.Errorf("wrapped: %w", err)
	t.Logf("divide: %v", err)
	t.Logf("Unwrap: %v", Unwrap(err))

	// As() on our error wrappers
	err = New("xx")
	var e1 *causes2
	if As(err, &e1) {
		t.Logf("e1: %v", e1)
	} else {
		t.Fail()
	}
}

func TestCauses2_WithCode(t *testing.T) {
	err := &causes2{
		Code:    Internal,
		Causers: nil,
		msg:     "ui err",
	}
	err.WithCode(NotFound).End()
	t.Logf("failed: %+v", err)
	t.Logf("failed: %+v", err.Cause())
	t.Logf("failed: %+v", err.Causes())
	t.Logf("failed: %+v", err.TypeIs(io.EOF))
	t.Logf("failed: %+v", err.IsEmpty())
}

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

func TestContainer(t *testing.T) {
	// as a inner errors container
	child := func() (err error) {
		errContainer := &causes2{}

		defer errContainer.Defer(&err)
		for _, r := range []error{io.EOF, io.ErrShortWrite, io.ErrClosedPipe, Internal} {
			errContainer.Attach(r)
		}

		return
	}

	err := child()
	t.Logf("failed: %+v", err)

	fmt.Printf("%+v", err)
	fmt.Printf("%#v", err)
	fmt.Printf("%v", err)
	fmt.Printf("%s", err)
	fmt.Printf("%q", err)

	// fmt.Printf("%n", err) // need go1.13+
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

func TestIsDescended(t *testing.T) {
	err1 := &causes2{
		Code:    Internal,
		Causers: nil,
		msg:     "ui err %v",
	}

	err2 := err1.FormatWith("1st")
	if !IsDescended(err1, err2) {
		t.Fatalf("bad test on IsDescended(err1, err2)")
	}

	err3 := New("any error tmpl with %v")
	err4 := err3.FormatWith("huahua")
	if !IsDescended(err3, err4) {
		t.Fatalf("bad test on IsDescended(err3, err4)")
	}

}
