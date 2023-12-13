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

	err := v3.Join(err1, err2)

	fmt.Printf("%T, %v\n", err, err)

	if v3.Is(err, err1) {
		t.Log("err is err1")
	} else {
		t.Fatal("FAILED: expecting err is err1")
	}

	if v3.Is(err, err2) {
		t.Log("err is err2")
	} else {
		t.Fatal("FAILED: expecting err is err2")
	}

	err3 := fmt.Errorf("error3: %w", err)
	fmt.Printf("%T, %v\n", err3, v3.Unwrap(err3))

	if v3.Is(err3, err1) {
		t.Log("err3 is err1")
	} else {
		t.Fatal("FAILED: expecting err3 is err1")
	}

	if v3.Is(err3, err2) {
		t.Log("err3 is err2")
	} else {
		t.Fatal("FAILED: expecting err3 is err2")
	}

	if !v3.Is(err2, err3) {
		t.Log("err2 isn't err3")
	} else {
		t.Fatal("FAILED: expecting err2 is err3")
	}
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
