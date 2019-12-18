// Copyright © 2019 Hedzr Yeh.

package errors

import "strings"

// LeftPad adds the leading char `pad` to `s`, and truncate it to the length `width`.
//
// LeftPad("89", '0', 6) => "000089"
//
// LeftPad returns an empty string "" if width is negative or zero.
// LeftPad returns the source string `s` if its length is larger than `width`.
func LeftPad(s string, pad rune, width int) string {
	if width <= 0 {
		return ""
	}

	if len(s) >= width {
		return s
	}

	var b strings.Builder
	for i := 0; i < width-len(s); i++ {
		b.WriteRune(pad)
	}
	b.WriteString(s)
	return b.String()
}

// Left returns the left `length` substring of `s`.
// Left returns the whole source string `s` if its length is larger than `length`
//
// Left("12345",3) => "123"
func Left(s string, length int) string {
	if length <= 0 {
		return ""
	}
	if length < len(s) {
		return s[:length]
	}
	return s
}

// Right returns the right `length` substring of `s`.
// Right returns the whole source string `s` if its length is larger than `length`
//
// Right("12345",3) => "345"
func Right(s string, length int) string {
	if length <= 0 {
		return ""
	}
	if length < len(s) {
		return s[len(s)-length:]
	}
	return s
}
