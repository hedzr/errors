package old_test

import (
	"fmt"
	"gopkg.in/hedzr/errors.v2"
	"gopkg.in/hedzr/errors.v2/old"
	"io"
	"testing"
)

const (
	BUG1001 errors.Code = 1001
	BUG1002 errors.Code = 1002
	BUG1003 errors.Code = 1003
)

var (
	errBug1 = old.New("bug 1")
	errBug2 = old.New("bug 2")
	errBug3 = old.New("bug 3")
)

func init() {
	BUG1001.Register("BUG1001")
	BUG1002.Register("BUG1002")
}

func TestErrorTemplate(t *testing.T) {
	errTmpl1001 := BUG1001.NewTemplate("something is wrong %v")
	err4 := errTmpl1001.FormatNew("ok").Attach(errBug1)
	fmt.Println(err4)
	fmt.Printf("%+v\n", err4)
}

func TestAxMain(t *testing.T) {
	fmt.Println(BUG1002.String())
	fmt.Println(errBug1)
	fmt.Println(BUG1003)

	err := old.New("something").Attach(errBug1, errBug2)
	err2 := old.Wrap(errBug3, "")
	err3 := BUG1002.New("info wrong").Attach(errBug1)
	fmt.Println(err)
	fmt.Println(err2.Error())
	fmt.Println(err3)

	e := err3.Unwrap()
	fmt.Println(e)
	e = err3.Unwrap()
	fmt.Println(e)
	e = err2.Unwrap()
	fmt.Println(e)
	e = err.Unwrap()
	fmt.Println(e)

	err = old.New("something").Attach(io.ErrClosedPipe)
	e = err.Unwrap()
	fmt.Println(e)
}
