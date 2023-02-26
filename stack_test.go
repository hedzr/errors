package errors

import (
	"fmt"
	"io"
	"testing"
)

func TestFuncname(t *testing.T) {
	fn := funcname("/home/fine/app.TestFuncname()")
	println(fn)
}

func TestStack_StackTrace(t *testing.T) {
	var s = callers(0)
	println(s.StackTrace())

	fmt.Printf("%+v", s)
	fmt.Printf("%#v", s)
	fmt.Printf("%v", s)
	fmt.Printf("%s", s)
	fmt.Printf("%q", s)

	fmt.Printf("%+v", s.StackTrace())
	fmt.Printf("%#v", s.StackTrace())
	fmt.Printf("%v", s.StackTrace())
	fmt.Printf("%s", s.StackTrace())
	fmt.Printf("%q", s.StackTrace())

	for _, frame := range s.StackTrace() {
		fmt.Printf("%+v", frame)
		fmt.Printf("%#v", frame)
		fmt.Printf("%v", frame)
		fmt.Printf("%s", frame)
		fmt.Printf("%q", frame)
		fmt.Printf("%n", frame)

	}
}

func TestWithStack(t *testing.T) {
	err := WithStack(nil)
	t.Logf("failed: %+v", err)

	err = WithStack(io.EOF)
	t.Logf("failed: %+v", err)
}

func TestWithStackInfo(t *testing.T) {
	err := &WithStackInfo{}
	err.WithErrors(io.EOF).
		WithCode(Internal).
		WithMessage("").
		WithMessage("%v", "").
		WithSkip(1).
		WithData(1, 2, io.ErrShortWrite).
		WithTaggedData(TaggedData{"1": 1}).
		End()
	err.WithCause(io.ErrNoProgress).End()
	t.Logf(" err is: %+v", err)
	t.Logf(" err.Cause() is: %v", err.Cause())
	t.Logf(" err.Causes() are: %v", err.Causes())

	fmt.Printf("%+v", err)
	fmt.Printf("%#v", err)
	fmt.Printf("%v", err)
	fmt.Printf("%s", err)
	fmt.Printf("%q", err)
	fmt.Printf("%n", err)
}
