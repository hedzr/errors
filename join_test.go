// Copyright © 2023 Hedzr Yeh.

package errors_test

import (
	"errors"
	"fmt"
	"testing"

	v3 "gopkg.in/hedzr/errors.v3"
)

func TestJoinErrors(t *testing.T) {
	err1 := errors.New("err1")
	err2 := errors.New("err2")
	err := v3.Join(err1, err2)
	fmt.Printf("%T, %v\n", err, err)
	if v3.Is(err, err1) {
		t.Log("err is err1")
	} else {
		t.Fatal("expecting err is err1")
	}
	if v3.Is(err, err2) {
		t.Log("err is err2")
	} else {
		t.Fatal("expecting err is err2")
	}
}

type DivisionError struct {
	IntA int
	IntB int
	Msg  string
}

func (e *DivisionError) Error() string {
	return e.Msg
}

func Divide(a, b int) (int, error) {
	if b == 0 {
		return 0, &DivisionError{
			Msg:  fmt.Sprintf("cannot divide '%d' by zero", a),
			IntA: a, IntB: b,
		}
	}
	return a / b, nil
}
