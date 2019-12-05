// Copyright © 2019 Hedzr Yeh.

package errors_test

import (
	"github.com/hedzr/errors"
	"io"
	"testing"
)

const (
	BUG1001 errors.Code = -1001
	BUG1002 errors.Code = -1002
	BUG1003 errors.Code = -1003
	BUG1004 errors.Code = -1004
	BUG1005 errors.Code = -1005
)

var (
	errBug1001 = errors.NewCodedError(BUG1001).Msg("something is wrong").Attach(io.EOF)
)

var (
	errBug1 = errors.New("bug 1")
	errBug2 = errors.New("bug 2, %v, %d", []string{"a", "b"}, 5)
	errBug3 = errors.New("bug 3")

	eb1  = errors.NewTemplate("xxbug 1, cause: %v")
	eb11 = errors.New("").Msg("first, %v", "ok").Template("xxbug11, cause")
	eb2  = errors.NewCodedError(BUG1004).Template("xxbug4, cause: %v")
	eb3  = errors.NewCodedError(BUG1005).Template("xxbug5, cause: none")
	eb31 = errors.NewCodedError(BUG1004).Msg("first, %v", "ok").Template("xxbug4.31, cause: %v")
	eb4  = errors.NewCodedError(BUG1005).Template("xxbug54, cause: none")
)

func init() {
	BUG1001.Register("BUG1001")
	BUG1002.Register("BUG1002")

	BUG1004.Register("BUG1004")
	BUG1005.Register("BUG1005")
}

func TestAll(t *testing.T) {
	t.Log(BUG1001.String())
	t.Log(errBug1001)
	t.Log(BUG1003)

	var err error
	err = errors.New("something").Attach(errBug1, errBug2).Nest(errBug1, errBug2).Msg("anything")

	err2 := errors.NewWithError(errBug3)
	err3 := errors.NewCodedError(BUG1002).Attach(errBug1).Code(BUG1002)
	err4 := errors.NewCodedError(BUG1002).Msg("xx").Nest(errBug1)
	t.Log(err)
	t.Log(err2.Error())
	t.Log(err3)
	t.Log(err4.Error())

	e := err4.Unwrap()
	t.Log(e)
	e = err3.Unwrap()
	t.Log(e)
	e = err2.Unwrap()
	t.Log(e)
	e = errors.Unwrap(err)
	t.Log(e)

	err = errors.New("something").Attach(io.ErrClosedPipe)
	e = errors.Unwrap(err)
	t.Log(e)
	if errors.As(e, &io.ErrClosedPipe) {
		t.Log("As() ok")
	}
	if errors.Is(e, io.ErrClosedPipe) {
		t.Log("Is() ok")
	}

	err = eb1.Format("resources exhausted")
	t.Log(err)
	err = eb11.Format()
	t.Log(err)
	err = eb2.Format("resources exhausted")
	t.Log(err)
	err = eb3.Format()
	t.Log(err)
	err = eb31.Format("resources exhausted")
	t.Log(err)
	err = eb4.Format()
	t.Log(err)
}
