package errors

import (
	"io"
	"testing"
)

func TestCauses2_New(t *testing.T) {
	err := &causes2{
		Code:    Internal,
		Causers: nil,
		msg:     "ui err",
	}

	t.Logf("failed: %+v", err)

	err = &causes2{
		Code:    0,
		Causers: nil,
		msg:     "ui err",
	}
	_ = err.WithAttach(io.EOF, io.ErrClosedPipe)
	t.Logf("failed: %+v", err)

	err = &causes2{
		Code:    Internal,
		Causers: nil,
		msg:     "ui err",
	}
	err.WithAttach(io.EOF, io.ErrClosedPipe).End()
	t.Logf("failed: %+v", err)

}
