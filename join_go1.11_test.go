// Copyright © 2023 Hedzr Yeh.

//go:build go1.11
// +build go1.11

package errors_test

import (
	"errors"
	"fmt"
	"testing"

	v3 "gopkg.in/hedzr/errors.v3"
)

func TestJoinErrorsStdFormatGo111(t *testing.T) {
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
}
