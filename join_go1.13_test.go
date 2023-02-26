// Copyright © 2023 Hedzr Yeh.

//go:build go1.13
// +build go1.13

package errors_test

import (
	"errors"
	"fmt"
	"io"
	"testing"

	v3 "gopkg.in/hedzr/errors.v3"
)

func TestJoinErrorsStdFormat(t *testing.T) {
	err1 := errors.New("err1")
	err2 := errors.New("err2")
	err := errors.Join(err1, err2)
	fmt.Printf("%T, %v\n", err, err)
	if errors.Is(err, err1) {
		t.Log("err is err1")
	} else {
		t.Fatal("expecting err is err1")
	}
	if errors.Is(err, err2) {
		t.Log("err is err2")
	} else {
		t.Fatal("expecting err is err2")
	}

	err3 := fmt.Errorf("error3: %w", err)
	fmt.Printf("%T, %v\n", err3, errors.Unwrap(err3))
	if errors.Is(err3, err1) {
		t.Log("err3 is err1")
	} else {
		t.Fatal("expecting err3 is err1")
	}
	if errors.Is(err3, err2) {
		t.Log("err3 is err2")
	} else {
		t.Fatal("expecting err3 is err1")
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

func dummy(t *testing.T) error {
	a, b := 10, 0
	result, err := Divide(a, b)
	if err != nil {
		var divErr *DivisionError
		switch {
		case errors.As(err, &divErr):
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

func TestCauses2_errors(t *testing.T) {
	err := io.EOF

	if !errors.Is(err, io.EOF) {
		t.Fail()
	}

	err = dummy(t)
	err = fmt.Errorf("wrapped: %w", err)
	t.Logf("divide: %v", err)
	t.Logf("Unwrap: %v", errors.Unwrap(err))
}

func TestCauses2_errorsV3(t *testing.T) {
	err := io.EOF

	if !v3.Is(err, io.EOF) {
		t.Fail()
	}

	err = dummyV3(t)
	err = fmt.Errorf("wrapped: %w", err)
	t.Logf("divide: %v", err)
	t.Logf("Unwrap: %v", v3.Unwrap(err))
}
