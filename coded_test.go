package errors

import (
	"fmt"
	"io"
	"testing"
)

type bizErr struct { //nolint:unused //usable
	num int
}

func (e *bizErr) Error() string { //nolint:unused //usable
	return fmt.Sprintf("%v", e.num)
}

func TestCode_WithCode(t *testing.T) {
	c := Internal
	c1 := (&c).WithCode(NotFound)

	t.Logf("failed: %+v", c1)

	c = Code(111)
	t.Logf("failed: %+v", c)
}

func TestCode_Register(t *testing.T) {
	c := Code(111)
	t.Logf("failed: %+v", c)

	_ = c.Register("Code111")
	t.Logf("failed: %+v", c)
}

// func TestCodeEqual(t *testing.T) {
//	be := &bizErr{1}
//	err := InvalidArgument.New("wrong").Attach(be)
//
//	//var e *bizErr
//	e1 := err.Unwrap().(*WithCodeInfo)
//
//	if !e1.Equal(InvalidArgument) {
//		t.Fatal("expecting e1 is equal to InvalidArgument")
//	}
//	if !Equal(e1, InvalidArgument) {
//		t.Fatal("expecting e1 is equal to InvalidArgument")
//	}
// }
//
// func TestCodeAsIsAndSoOn(t *testing.T) {
//	be := &bizErr{1}
//	err := InvalidArgument.New("wrong").Attach(be)
//
//	var e *bizErr
//	e1 := err.Unwrap().(*WithCodeInfo)
//	if !e1.As(&e) {
//		t.Fatal("WithCodeInfo.As() failed.")
//	}
//
//	if !err.Is(be) {
//		t.Fatal("WithCodeInfo.Is() failed.")
//	}
// }

// func TestCodes(t *testing.T) {
//	be := &bizErr{1}
//	err := InvalidArgument.New("wrong").Attach(be)
//	t.Log(err)
//	t.Logf("%+v", err)
//
//	exm := Internal.New("msg")
//	ex := exm.Unwrap()
//	if x, ok := ex.(interface{ Code() Code }); ok {
//		t.Log(x.Code())
//		t.Logf("Internal: %q | cause = %v", x, ex.(*WithCodeInfo).Cause())
//	} else {
//		t.Fatalf("Internal: %v", ex)
//	}
//
//	if !old.Is(err, be) {
//		t.Fatal("wrong Is(): expecting be")
//	}
//	if old.Is(err, io.EOF) {
//		t.Fatal("wrong Is(): shouldn't be like to io.EOF")
//	}
// }

func TestCodesEqual(t *testing.T) {
	err := InvalidArgument.New("wrong").WithErrors(io.ErrShortWrite)

	ok := Is(err, InvalidArgument)
	if !ok {
		t.Fatal("want Equal() return true but got false")
	}
}

func TestCodesRegister(t *testing.T) {
	const illegalStateEx Code = MinErrorCode - 1
	_ = RegisterCode(int(illegalStateEx), "I'm in an illegal state (ext for testing).")
	t.Log(illegalStateEx)
	t.Logf("%+v", illegalStateEx)
}
