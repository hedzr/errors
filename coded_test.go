package errors

import (
	"fmt"
	"io"
	"testing"
)

type bizErr struct {
	num int
}

func (e *bizErr) Error() string {
	return fmt.Sprintf("%v", e.num)
}

func TestCodeEqual(t *testing.T) {
	be := &bizErr{1}
	err := InvalidArgument.New("wrong").Attach(be)

	//var e *bizErr
	e1 := err.Unwrap().(*WithCodeInfo)

	if !e1.Equal(InvalidArgument) {
		t.Fatal("expecting e1 is equal to InvalidArgument")
	}
	if !Equal(e1, InvalidArgument) {
		t.Fatal("expecting e1 is equal to InvalidArgument")
	}
}

func TestCodeAsIsAndSoOn(t *testing.T) {
	be := &bizErr{1}
	err := InvalidArgument.New("wrong").Attach(be)

	var e *bizErr
	e1 := err.Unwrap().(*WithCodeInfo)
	if !e1.As(&e) {
		t.Fatal("WithCodeInfo.As() failed.")
	}

	if !err.Is(be) {
		t.Fatal("WithCodeInfo.Is() failed.")
	}
}

func TestCodes(t *testing.T) {
	be := &bizErr{1}
	err := InvalidArgument.New("wrong").Attach(be)
	t.Log(err)
	t.Logf("%+v", err)

	exm := Internal.New("msg")
	ex := exm.Unwrap()
	if x, ok := ex.(interface{ Code() Code }); ok {
		t.Log(x.Code())
		t.Logf("Internal: %q | cause = %v", x, ex.(*WithCodeInfo).Cause())
	} else {
		t.Fatalf("Internal: %v", ex)
	}

	if !Is(err, be) {
		t.Fatal("wrong Is(): expecting be")
	}
	if Is(err, io.EOF) {
		t.Fatal("wrong Is(): shouldn't be like to io.EOF")
	}
}

func TestCodesEqual(t *testing.T) {
	err := InvalidArgument.New("wrong").Attach(io.ErrShortWrite)

	ok := EqualR(err, InvalidArgument)
	if !ok {
		t.Fatal("want Equal() return true but got false")
	}
}
