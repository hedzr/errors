package errors

import (
	"fmt"
	"io"
	"strconv"
	"testing"
)

func TestAs_e1(t *testing.T) {
	// As() on our error wrappers
	err := New("xx")
	var e1 *causes2
	if As(err, &e1) {
		t.Logf("e1: %v", e1)
	} else {
		t.Fail()
	}
}

func TestAs_betterFormat(t *testing.T) {
	var err = New("Have errors").WithErrors(io.EOF, io.ErrShortWrite, io.ErrNoProgress)
	t.Logf("%v\n", err)

	var nestNestErr = New("Errors FOUND:").WithErrors(err, io.EOF)
	var nnnErr = New("Nested Errors:").WithErrors(nestNestErr, strconv.ErrRange)
	t.Logf("%v\n", nnnErr)
	t.Logf("%+v\n", nnnErr)
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

	fmt.Printf("%+v\n", err)
	fmt.Printf("%#v\n", err)
	fmt.Printf("%v\n", err)
	fmt.Printf("%s\n", err)
	fmt.Printf("%q\n", err)

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
