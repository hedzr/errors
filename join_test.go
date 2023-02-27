// Copyright © 2023 Hedzr Yeh.

package errors_test

import (
	"fmt"
	"testing"

	"gopkg.in/hedzr/errors.v3"
)

func TestJoinErrors(t *testing.T) {
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
}