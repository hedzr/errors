package errors_test

import (
	"fmt"
	"gopkg.in/hedzr/errors.v2"
	"io"
)

func Example_container() {
	err := sample(false)
	if err != nil {
		panic(err)
	} else {
		fmt.Printf("%+v\n", err)
	}

	err = sample(true)
	if err == nil {
		panic("want error")
	} else {
		fmt.Printf("%+v\n", err)
	}
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
	fmt.Printf(err)
	fmt.Printf("%+v\n", err)

	if !errors.Is(err, io.ErrShortWrite) {
		panic("wrong Is()")
	}
	if errors.Is(err, io.EOF) {
		panic("wrong Is()")
	}
}

func Example_customErrorCode() {
	const ErrnoMyFault errors.Code = 1001
	ErrnoMyFault.Register("MyFault")
	fmt.Printf("%+v\n", ErrnoMyFault)

	err := ErrnoMyFault.New("my fault message")
	fmt.Printf("%+v\n", err)

	tmpl := ErrnoMyFault.NewTemplate("my fault message: %v")
	err = tmpl.FormatNew("whoops")
	fmt.Printf("%+v\n", err)
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
	err := fn()
	fmt.Printf("%+v\n", err)

	// Example output:
	// error
	// github.com/pkg/errors_test.fn
	//         /home/dfc/src/github.com/pkg/errors/example_test.go:47
	// github.com/pkg/errors_test.ExampleCause_printf
	//         /home/dfc/src/github.com/pkg/errors/example_test.go:63
	// testing.runExample
	//         /home/dfc/go/src/testing/example.go:114
	// testing.RunExamples
	//         /home/dfc/go/src/testing/example.go:38
	// testing.(*M).Run
	//         /home/dfc/go/src/testing/testing.go:744
	// main.main
	//         /github.com/pkg/errors/_test/_testmain.go:104
	// runtime.main
	//         /home/dfc/go/src/runtime/proc.go:183
	// runtime.goexit
	//         /home/dfc/go/src/runtime/asm_amd64.s:2059
	// github.com/pkg/errors_test.fn
	// 	  /home/dfc/src/github.com/pkg/errors/example_test.go:48: inner
	// github.com/pkg/errors_test.fn
	//        /home/dfc/src/github.com/pkg/errors/example_test.go:49: middle
	// github.com/pkg/errors_test.fn
	//      /home/dfc/src/github.com/pkg/errors/example_test.go:50: outer
}
