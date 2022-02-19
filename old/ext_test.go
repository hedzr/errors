package old_test

import (
	"gopkg.in/hedzr/errors.v2/old"
	"io"
	"testing"
)

func TestSkip(t *testing.T) {
	err := old.Skip(1).Build()
	t.Logf("failed: %+v", err)
}

func TestForExample(t *testing.T) {

	err := old.New("some tips %v", "here")

	// attaches much more error causing
	for _, e := range []error{io.EOF, io.ErrClosedPipe} {
		_ = err.Attach(e)
	}

	t.Logf("failed: %+v", err)

	// use another number different to default to skip the error frames
	err = old.Skip(3).Message("some tips %v", "here").Build()
	t.Logf("failed: %+v", err)
}
