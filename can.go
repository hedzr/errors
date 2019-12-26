// Copyright © 2019 Hedzr Yeh.

package errors

// CanWalk tests if err is walkable
func CanWalk(err error) (ok bool) {
	_, ok = err.(Walkable)
	return
}

// CanRange tests if err is range-able
func CanRange(err error) (ok bool) {
	_, ok = err.(Ranged)
	return
}

// CanUnwrap tests if err is unwrap-able
func CanUnwrap(err error) (ok bool) {
	_, ok = err.(interface{ Unwrap() error })
	return
}

// CanIs tests if err is is-able
func CanIs(err error) (ok bool) {
	_, ok = err.(interface{ Is(error) bool })
	return
}

// CanAs tests if err is as-able
func CanAs(err error) (ok bool) {
	_, ok = err.(interface{ As(interface{}) bool })
	return
}
