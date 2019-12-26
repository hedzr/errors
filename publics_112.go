// Copyright © 2019 Hedzr Yeh.

//+build !go1.13

package errors

// HasWrappedError detects if nested or wrapped errors present
//
// nested error: ExtErr.inner
// wrapped error: fmt.Errorf("... %w ...", err)
func HasWrappedError(err error) (yes bool) {
	if ex, ok := err.(interface{ GetNestedError() *ExtErr }); ok {
		return ex.GetNestedError() != nil
	}
	return false
}
