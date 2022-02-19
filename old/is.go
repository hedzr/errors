package old

import "reflect"

// Is reports whether any error in `err`'s chain matches target.
//
// The chain consists of err itself followed by the sequence of errors obtained by
// repeatedly calling Unwrap.
//
// An error is considered to match a target if it is equal to that target or if
// it implements a method Is(error) bool such that Is(target) returns true.
func Is(err, target error) bool {
	if target == nil {
		return err == target
	}

	isComparable := reflect.TypeOf(target).Comparable()
	for {
		if isComparable {
			if err == target {
				return true
			}
		}
		if x, ok := err.(interface{ Is(error) bool }); ok && x.Is(target) {
			return true
		}
		// TODO: consider supporting target.Is(err). This would allow
		// user-definable predicates, but also may allow for coping with sloppy
		// APIs, thereby making it easier to get away with them.
		if err = Unwrap(err); err == nil {
			return false
		}
	}
}

// IsSlice tests err.Is for errs slice
func IsSlice(errs []error, target error) bool {
	if target == nil {
		//for _, e := range errs {
		//	if e == target {
		//		return true
		//	}
		//}
		return false
	}

	isComparable := reflect.TypeOf(target).Comparable()
	for {
		if isComparable {
			for _, e := range errs {
				if e == target {
					return true
				}
			}
			// return false
		}

		for _, e := range errs {
			if x, ok := e.(interface{ Is(error) bool }); ok && x.Is(target) {
				return true
			}
			//if err := Unwrap(e); err == nil {
			//	return false
			//}
		}
		return false
	}
}

// TypeIs reports whether any error in `err`'s chain matches target.
//
// The chain consists of err itself followed by the sequence of errors obtained by
// repeatedly calling Unwrap.
//
// An error is considered to match a target if it is equal to that target or if
// it implements a method Is(error) bool such that Is(target) returns true.
func TypeIs(err, target error) bool {
	if target == nil {
		return err == target
	}

	isComparable := reflect.TypeOf(target).Comparable()
	for {
		if isComparable {
			if reflect.TypeOf(target) == reflect.TypeOf(err) {
				return true
			}
		}
		if x, ok := err.(interface{ Is(error) bool }); ok && x.Is(target) {
			return true
		}
		// TODO: consider supporting target.Is(err). This would allow
		// user-definable predicates, but also may allow for coping with sloppy
		// APIs, thereby making it easier to get away with them.
		if err = Unwrap(err); err == nil {
			return false
		}
	}
}

// TypeIsSlice tests err.Is for errs slice
func TypeIsSlice(errs []error, target error) bool {
	if target == nil {
		//for _, e := range errs {
		//	if e == target {
		//		return true
		//	}
		//}
		return false
	}

	isComparable := reflect.TypeOf(target).Comparable()
	for {
		if isComparable {
			tt := reflect.TypeOf(target)
			for _, e := range errs {
				//if e == target {
				//	return true
				//}
				if reflect.TypeOf(e) == tt {
					return true
				}
			}
			// return false
		}

		for _, e := range errs {
			if x, ok := e.(interface{ Is(error) bool }); ok && x.Is(target) {
				return true
			}
			//if err := Unwrap(e); err == nil {
			//	return false
			//}
		}
		return false
	}
}
