package errors

import "testing"

func TestDumpStacksAsString(t *testing.T) {
	DumpStacksAsString(true)
}

func TestCanAttach(t *testing.T) {
	err := New("")
	t.Log(CanAttach(err))
	t.Log(CanAttach(Internal))
}

func TestCanCause(t *testing.T) {
	err := New("")
	t.Log(CanCause(err))
	t.Log(CanCause(Internal))
}

func TestCanUnwrap(t *testing.T) {
	err := New("")
	t.Log(CanUnwrap(err))
	t.Log(CanUnwrap(Internal))
}

func TestCanIs(t *testing.T) {
	err := New("")
	t.Log(CanIs(err))
	t.Log(CanIs(Internal))
}

func TestCanAs(t *testing.T) {
	err := New("")
	t.Log(CanAs(err))
	t.Log(CanAs(Internal))
}
