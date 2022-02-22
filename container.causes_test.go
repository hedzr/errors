package errors

import (
	"io"
	"testing"
)

func TestCauses(t *testing.T) {
	c := &causes{Causers: []error{io.EOF}}

	t.Logf("1. %%v : %v", c)
	t.Logf("2. %%+v : %+v", c)
	t.Logf("3. %%s : %s", c)
	t.Logf("4. %%q : %q", c)

	c1 := c.Cause()
	t.Logf("Cause(): %v", c1)

	c2 := c.Causes()
	t.Logf("Causes(): %v", c2)
}

func TestCausesZeroLength(t *testing.T) {
	c := &causes{Causers: []error{io.EOF}}
	c1 := c.Cause()
	t.Logf("Cause(): %v", c1)
	c2 := c.Causes()
	t.Logf("Causes(): %v", c2)
}

func TestCausesUnwrap(t *testing.T) {
	c := &causes{Causers: []error{io.EOF}}
	e := c.Unwrap()
	if e != io.EOF {
		t.Fatal("expecting Unwrap to io.EOF")
	}

	if !c.Is(io.EOF) {
		t.Fatal("expecting Is() io.EOF")
	}

	d := &causes{Causers: []error{c, io.ErrClosedPipe}}
	if !d.Is(io.EOF) {
		t.Fatal("expecting d.Is() io.EOF")
	}
	if !d.Is(io.ErrClosedPipe) {
		t.Fatal("expecting d.Is() io.ErrClosedPipe")
	}

	var e2 *causes
	if !d.As(&e2) || e2 != c {
		t.Fatal("As() failed")
	}
}

func TestCausesEmptyCausers(t *testing.T) {
	c := &causes{Causers: nil}

	t.Logf("1. %%v : %v", c)
	t.Logf("2. %%+v : %+v", c)
	t.Logf("3. %%s : %s", c)
	t.Logf("4. %%q : %q", c)

	c1 := c.Cause()
	t.Logf("Cause(): %v", c1)

	c2 := c.Causes()
	t.Logf("Causes(): %v", c2)

	if nil != c.Unwrap() {
		t.Fatal("expecting the return result is nil")
	}
}
