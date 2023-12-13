// Copyright © 2023 Hedzr Yeh.

//go:build go1.20
// +build go1.20

package errors_test

import (
	"errors"
	"fmt"
	"testing"

	v3 "gopkg.in/hedzr/errors.v3"
)

func TestJoinErrorsStd(t *testing.T) {
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

	err3 := fmt.Errorf("error3: %w", err)
	fmt.Printf("%T, %v\n", err3, v3.Unwrap(err3))
	if v3.Is(err3, err1) {
		t.Log("err3 is err1")
	} else {
		t.Fatal("expecting err3 is err1")
	}
	if v3.Is(err3, err2) {
		t.Log("err3 is err2")
	} else {
		t.Fatal("expecting err3 is err1")
	}
}
