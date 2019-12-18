// Copyright © 2019 Hedzr Yeh.

package errors_test

import (
	"github.com/hedzr/errors"
	"testing"
)

func TestStrings(t *testing.T) {
	var s string

	for _, tc := range []struct {
		str    string
		w      int
		expect string
	}{
		{"89", -6, ""},
		{"89", 6, "000089"},
		{"56789", 6, "056789"},
		{"456789", 6, "456789"},
		{"3456789", 6, "3456789"},
		{"123456789", 6, "123456789"},
	} {
		if s = errors.LeftPad(tc.str, '0', tc.w); s != tc.expect {
			t.Fatalf("wrong leftpad(%q,'0',%d), expect %q, but got %q !", tc.str, tc.w, tc.expect, s)
		}
	}

	for _, tc := range []struct {
		str    string
		w      int
		expect string
	}{
		{"3456789", -3, ""},
		{"3456789", 0, ""},
		{"3456789", 1, "3"},
		{"3456789", 5, "34567"},
		{"3456789", 6, "345678"},
		{"3456789", 7, "3456789"},
		{"3456789", 8, "3456789"},
		{"3456789", 9, "3456789"},
	} {
		if s = errors.Left(tc.str, tc.w); s != tc.expect {
			t.Fatalf("wrong left(%q,%d), expect %q, but got %q !", tc.str, tc.w, tc.expect, s)
		}
	}

	for _, tc := range []struct {
		str    string
		w      int
		expect string
	}{
		{"123456789", -7, ""},
		{"123456789", 0, ""},
		{"123456789", 1, "9"},
		{"123456789", 2, "89"},
		{"123456789", 6, "456789"},
		{"123456789", 8, "23456789"},
		{"123456789", 9, "123456789"},
		{"123456789", 10, "123456789"},
		{"123456789", 11, "123456789"},
	} {
		if s = errors.Right(tc.str, tc.w); s != tc.expect {
			t.Fatalf("wrong right(%q,'0',%d), expect %q, but got %q !", tc.str, tc.w, tc.expect, s)
		}
	}
}
