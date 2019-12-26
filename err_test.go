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

func TestExtErr(t *testing.T) {

	var err error
	err = errors.New("something").Attach(errBug1, errBug2).Nest(errBug1, errBug2).Msg("anything")
	t.Log(err.(*errors.ExtErr).NoCannedError())
	t.Log(err.(*errors.ExtErr).Is(err))

	e1 := BUG1003.New("z")
	e2 := BUG1001.New("z1").Nest(e1)
	e3 := BUG1002.New("z2").Attach(e2).Nest(e1)

	errors.Attach(err, e3)
	errors.Nest(err, e1)

	err.(*errors.ExtErr).Range(func(err error) (stop bool) {
		return true
	})
}

func TestAll(t *testing.T) {
	t.Log(BUG1001.String())
	t.Log(errBug1001)
	t.Log(BUG1003)

	e1 := BUG1003.New("z")
	t.Log(e1.Number())
	if !e1.Equal(BUG1003) {
		t.Fatal("wrong equal")
	}
	if !e1.EqualRecursive(BUG1003) {
		t.Fatal("wrong equalr 1.3")
	}
	e2 := BUG1001.New("z1").Nest(e1)
	if !e2.EqualRecursive(BUG1003) {
		t.Fatal("wrong equalr 2.3")
	}
	if !e2.EqualRecursive(BUG1001) {
		t.Fatal("wrong equalr 2.1")
	}
	e3 := BUG1002.New("z2").Attach(e2).Nest(e1)
	errors.Walk(e3, func(err error) (stop bool) {
		t.Logf("  ..w.. : %+v", err)
		return false
	})
	if !e3.EqualRecursive(BUG1003) {
		t.Fatal("wrong equal 3.3")
	}
	if !e3.EqualRecursive(BUG1002) {
		t.Fatal("wrong equal 3.2")
	}
	if !e3.EqualRecursive(BUG1001) {
		t.Fatal("wrong equal 3.1")
	}
	errors.Range(e2, func(err error) (stop bool) {
		t.Logf("  ..r.. : %+v", err)
		return false
	})

	t.Log(e2.NoCannedError())
	t.Log(e2.IsBoth(BUG1003))
	t.Log(e2.IsBoth(BUG1001, BUG1002))
	t.Log(e2.IsAny(BUG1002, BUG1001))
	t.Log(e2.IsAny(BUG1002))
	t.Log(e3.GetAttachedErrors(), e3.GetNestedError(), e3.GetMsgString(), e3.GetTemplateString())
	t.Log(errors.CanAs(e3), errors.CanIs(e3), errors.CanUnwrap(e3), errors.CanRange(e3), errors.CanWalk(e3))

	//

	t.Log(errors.Equal(e2, BUG1003))
	t.Log(errors.IsBoth(e2, BUG1003))
	t.Log(errors.IsBoth(e2, BUG1001, BUG1002))
	t.Log(errors.IsAny(e2, BUG1002, BUG1001))
	t.Log(errors.IsAny(e2, BUG1002))

	t.Log(errors.Equal(io.EOF, BUG1003))
	t.Log(errors.IsBoth(io.EOF, BUG1003))
	t.Log(errors.IsAny(io.EOF, BUG1002, BUG1001))

	errors.Attach(e2, e3)
	errors.Nest(e2, e1)

	//

	var err error
	err = errors.New("something").Attach(errBug1, errBug2).Nest(errBug1, errBug2).Msg("anything")
	t.Log(err.(*errors.ExtErr).NoCannedError())
	t.Log(err.(*errors.ExtErr).Is(err))

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

func TestIsAs(t *testing.T) {
	var err error
	err = errors.New("something").Attach(errBug1, errBug2).Nest(errBug1, errBug2).Msg("anything")
	if !errors.Is(err, errBug1) {
		t.Fatalf("wrong! expect err ==IS== errBug1")
	}
	if !errors.Is(err, errBug2) {
		t.Fatalf("wrong! expect err ==IS== errBug2")
	}

	err2 := errors.NewWithError(io.ErrShortWrite).Nest(io.EOF)
	if !errors.Is(err2, io.ErrShortWrite) {
		t.Fatalf("wrong! expect err2 ==IS== io.ErrShortWrite")
	}
	if !errors.Is(err2, io.EOF) {
		t.Fatalf("wrong! expect err2 ==IS== io.EOF")
	}

}

func TestNest(t *testing.T) {
	var err error

	err = errors.New("1").Nest(io.EOF).Nest(io.ErrShortWrite).Nest(io.ErrShortBuffer)
	t.Log(err)
	t.Logf("%#v\n", err)

	err = errors.New("1").Attach(io.EOF).Attach(io.ErrShortWrite).Attach(io.ErrShortBuffer)
	t.Log(err)
	t.Logf("%#v\n", err)
}
