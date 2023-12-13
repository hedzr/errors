package errors

import (
	"fmt"
	"io"
	"testing"
)

func TestFuncname(t *testing.T) {
	fn := funcname("/home/fine/app.TestFuncname()")
	t.Log(fn)
}

func TestStack_StackTrace(t *testing.T) {
	s := callers(0)
	t.Log(s.StackTrace())

	fmt.Printf("%+v\n", s)
	fmt.Printf("%#v\n", s)
	fmt.Printf("%v\n", s)
	fmt.Printf("%s\n", s)
	fmt.Printf("%q\n", s)

	fmt.Printf("%+v\n", s.StackTrace())
	fmt.Printf("%#v\n", s.StackTrace())
	fmt.Printf("%v\n", s.StackTrace())
	fmt.Printf("%s\n", s.StackTrace())
	fmt.Printf("%q\n", s.StackTrace())

	for _, frame := range s.StackTrace() {
		fmt.Printf("%+v\n", frame)
		fmt.Printf("%#v\n", frame)
		fmt.Printf("%v\n", frame)
		fmt.Printf("%s\n", frame)
		fmt.Printf("%q\n", frame)
		fmt.Printf("%n\n", frame)
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

	fmt.Printf("v: %v\n", err)
	fmt.Printf("s: %s\n", err)
	fmt.Printf("q: %q\n", err)
	fmt.Printf("n: %n\n", err)

	fmt.Printf("+v: %+v\n", err)
	fmt.Printf("#v: %#v\n", err)
}
