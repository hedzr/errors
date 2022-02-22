package errors_test

import (
	"fmt"
	"gopkg.in/hedzr/errors.v2"
	"io"
	"testing"
)

func TestForExamples(t *testing.T) {
	// Example_container()
	// Example_errorCode()
	// Example_errorCodeCustom()
	Example_errorTemplate()
	// ExampleWrap_extended()
}

func Example_container() {
	// err := sample(false)
	c := errors.NewContainer("sample error")
	err := c.Error()
	if err != nil {
		panic(err)
	} else {
		fmt.Printf("1. want nil: %v\n", err)
	}

	// err = sample(true)
	c = errors.NewContainer("sample error")
	// in a long loop, we can add many sub-errors into container 'c'...
	errors.AttachTo(c, io.EOF, io.ErrUnexpectedEOF, io.ErrShortBuffer, io.ErrShortWrite)
	// and we extract all of them as a single parent error object now.
	err = c.Error()
	if err == nil {
		panic("want error")
	} else {
		fmt.Printf("2. %v\n", err)
	}

	// Example output:
	// 1. want nil: <nil>
	// 2. [EOF, unexpected EOF, short buffer, short write]
}

func sample(simulate bool) (err error) {
	c := errors.NewContainer("sample error")
	if simulate {
		errors.AttachTo(c, io.EOF, io.ErrUnexpectedEOF, io.ErrShortBuffer, io.ErrShortWrite)
	}
	err = c.Error()
	return
}

func Example_errorCode() {
	err := errors.InvalidArgument.New("wrong").Attach(io.ErrShortWrite)
	fmt.Println(err)
	fmt.Printf("%+v\n", err)

	if !errors.Is(err, io.ErrShortWrite) {
		panic("wrong Is()")
	}
	if errors.Is(err, io.EOF) {
		panic("wrong Is()")
	}

	// Example output:
	// INVALID_ARGUMENT|wrong|short write
	// -3|INVALID_ARGUMENT|wrong|short write
	// gopkg.in/hedzr/errors%2ev2.Code.New
	// /Users/hz/hzw/golang-dev/src/github.com/hedzr/errors/coded.go:365
	// gopkg.in/hedzr/errors%2ev2_test.Example_errorCode
	// /Users/hz/hzw/golang-dev/src/github.com/hedzr/errors/example_test.go:50
	// gopkg.in/hedzr/errors%2ev2_test.TestForExamples
	// /Users/hz/hzw/golang-dev/src/github.com/hedzr/errors/example_test.go:11
	// testing.tRunner
	// /usr/local/opt/go/libexec/src/testing/testing.go:909
	// runtime.goexit
	// /usr/local/opt/go/libexec/src/runtime/asm_amd64.s:1357
}

func Example_errorCodeCustom() {
	const ErrnoMyFault errors.Code = 1101
	ErrnoMyFault.Register("MyFault")
	fmt.Printf("%+v\n", ErrnoMyFault)

	err := ErrnoMyFault.New("my fault message")
	fmt.Printf("%+v\n", err)

	// Example output:
	// MyFault
	// 1101|MyFault|my fault message
	// gopkg.in/hedzr/errors%2ev2.Code.New
	// /Users/hz/hzw/golang-dev/src/github.com/hedzr/errors/coded.go:365
	// gopkg.in/hedzr/errors%2ev2_test.Example_errorCodeCustom
	// /Users/hz/hzw/golang-dev/src/github.com/hedzr/errors/example_test.go:80
	// gopkg.in/hedzr/errors%2ev2_test.TestForExamples
	// /Users/hz/hzw/golang-dev/src/github.com/hedzr/errors/example_test.go:11
	// testing.tRunner
	// /usr/local/opt/go/libexec/src/testing/testing.go:909
	// runtime.goexit
	// /usr/local/opt/go/libexec/src/runtime/asm_amd64.s:1357
}

func Example_errorTemplate() {
	const ErrnoMyFault errors.Code = 1101
	ErrnoMyFault.Register("MyFault")

	tmpl := ErrnoMyFault.NewTemplate("my fault message: %v")
	err := tmpl.FormatNew("whoops")
	fmt.Printf("%+v\n", err)

	// Example output:
	// 1101|MyFault|my fault message: whoops
	// gopkg.in/hedzr/errors%2ev2.(*WithCodeInfo).FormatNew
	// /Users/hz/hzw/golang-dev/src/github.com/hedzr/errors/coded.go:270
	// gopkg.in/hedzr/errors%2ev2_test.Example_errorTemplate
	// /Users/hz/hzw/golang-dev/src/github.com/hedzr/errors/example_test.go:106
	// gopkg.in/hedzr/errors%2ev2_test.TestForExamples
	// /Users/hz/hzw/golang-dev/src/github.com/hedzr/errors/example_test.go:14
	// testing.tRunner
	// /usr/local/opt/go/libexec/src/testing/testing.go:909
	// runtime.goexit
	// /usr/local/opt/go/libexec/src/runtime/asm_amd64.s:1357

}

func ExampleNew() {
	err := errors.New("whoops")
	fmt.Println(err)

	// Output: whoops
}

func fn() error {
	e1 := errors.New("error")
	e2 := errors.Wrap(e1, "inner")
	e3 := errors.Wrap(e2, "middle")
	return errors.Wrap(e3, "outer")
}

func ExampleWrap() {
	cause := errors.New("whoops")
	err := errors.Wrap(cause, "oh noes")
	fmt.Println(err)

	// Output: oh noes: whoops
}

func ExampleWrap_extended() {
	// err := fn()
	e1 := errors.New("error")
	e2 := errors.Wrap(e1, "inner")
	e3 := errors.Wrap(e2, "middle")
	err := errors.Wrap(e3, "outer")
	fmt.Printf("%v\n", err)
	fmt.Printf("%+v\n", err)

	// Example output:
	// outer: middle: inner: error
	// outer: middle: inner: error
	// gopkg.in/hedzr/errors%2ev2_test.fn
	// /Users/hz/hzw/golang-dev/src/github.com/hedzr/errors/example_test.go:136
	// gopkg.in/hedzr/errors%2ev2_test.ExampleWrap_extended
	// /Users/hz/hzw/golang-dev/src/github.com/hedzr/errors/example_test.go:148
	// gopkg.in/hedzr/errors%2ev2_test.TestForExamples
	// /Users/hz/hzw/golang-dev/src/github.com/hedzr/errors/example_test.go:15
	// testing.tRunner
	// /usr/local/opt/go/libexec/src/testing/testing.go:909
	// runtime.goexit
	// /usr/local/opt/go/libexec/src/runtime/asm_amd64.s:1357
}
