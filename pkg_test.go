package errors

import (
	"io"
	"testing"
)

// func err1() error {
// 	return errors.Wrap(io.EOF, "wrong")
// }
//
// func err2() error {
// 	return errors.WithStack(io.EOF)
// }

func TestPkgWrap(t *testing.T) {
	// t.Log(err2())
	//
	// err1 := err1()
	// t.Log(err1)

	err2 := New("fmt err")
	t.Logf("%+v", err2)

	err3 := NewCodedError(Canceled, io.EOF)
	t.Logf("%+v", err3)
	t.Logf("%v", err3)
	t.Logf("%s", err3)
	t.Logf("%q", err3)
}
