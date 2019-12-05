// Copyright © 2019 Hedzr Yeh.

package errors_test

import (
	"github.com/hedzr/errors"
	"io"
	"testing"
)

const (
	BUG1001 errors.Code = 1001
	BUG1002 errors.Code = 1002
	BUG1003 errors.Code = 1003
)

var (
	errBug1001 = errors.NewWithCodeMsg(BUG1001, "something is wrong", io.EOF)
)

var (
	errBug1 = errors.New("bug 1")
	errBug2 = errors.New("bug 2")
	errBug3 = errors.New("bug 3")
)

func init() {
	BUG1001.Register("BUG1001")
	BUG1002.Register("BUG1002")
}

func TestAll(t *testing.T) {
	t.Log(BUG1001.String())
	t.Log(errBug1001)
	t.Log(BUG1003)

	err := errors.New("something", errBug1, errBug2)
	err2 := errors.NewWithError(errBug3)
	err3 := errors.NewWithCode(BUG1002, errBug1)
	err4 := errors.NewWithCodeMsg(BUG1002, "xx", errBug1)
	t.Log(err)
	t.Log(err2)
	t.Log(err3)
	t.Log(err4)

}
