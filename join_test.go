// Copyright © 2023 Hedzr Yeh.

package errors_test

import (
	"errors"
	"fmt"
	"io"
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

func dummyV3(t *testing.T) error {
	a, b := 10, 0
	result, err := Divide(a, b)
	if err != nil {
		var divErr *DivisionError
		switch {
		case v3.As(err, &divErr):
			fmt.Printf("%d / %d is not mathematically valid: %s\n",
				divErr.IntA, divErr.IntB, divErr.Error())
		default:
			fmt.Printf("unexpected division error: %s\n", err)
			t.Fail()
		}
		return err
	}

	fmt.Printf("%d / %d = %d\n", a, b, result)
	return err
}

func TestCauses2_errorsV3(t *testing.T) {
	err := io.EOF

	if !v3.Is(err, io.EOF) {
		t.Fatal("FAILED: expecting err is io.EOF")
	}

	err = dummyV3(t)
	err = fmt.Errorf("wrapped: %w", err)
	t.Logf("divide: %v", err)
	t.Logf("Unwrap: %v", v3.Unwrap(err))
}
