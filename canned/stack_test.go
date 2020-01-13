package canned

import (
	"fmt"
	"github.com/pkg/errors"
	"regexp"
	"runtime"
	"strings"
	"testing"
)

var initpc, _, _, _ = runtime.Caller(0)

func TestFrameLine(t *testing.T) {
	var tests = []struct {
		Frame
		want int
	}{{
		Frame(initpc),
		12, // 9,
	}, {
		func() Frame {
			var pc, _, _, _ = runtime.Caller(0)
			return Frame(pc)
		}(),
		23, // 20,
	}, {
		func() Frame {
			var pc, _, _, _ = runtime.Caller(1)
			return Frame(pc)
		}(),
		25, // 28,
	}, {
		Frame(0), // invalid PC
		0,
	}}

	for _, tt := range tests {
		got := tt.Frame.line()
		want := tt.want
		if want != got {
			t.Errorf("Frame(%v): want: %v, got: %v", uintptr(tt.Frame), want, got)
		}
	}
}

type X struct{}

func (x X) val() Frame {
	var pc, _, _, _ = runtime.Caller(0)
	return Frame(pc)
}

func (x *X) ptr() Frame {
	var pc, _, _, _ = runtime.Caller(0)
	return Frame(pc)
}

func TestFrameFormat(t *testing.T) {
	var tests = []struct {
		Frame
		format string
		want   string
	}{{
		Frame(initpc),
		"%s",
		"stack_test.go",
	}, {
		Frame(initpc),
		"%+s",
		"github.com/hedzr/errors/canned.init\n" +
			"\t.+/github.com/hedzr/errors/canned/stack_test.go",
	}, {
		Frame(0),
		"%s",
		"unknown",
	}, {
		Frame(0),
		"%+s",
		"unknown",
	}, {
		Frame(initpc),
		"%d",
		"12",
	}, {
		Frame(0),
		"%d",
		"0",
	}, {
		Frame(initpc),
		"%n",
		"init",
	}, {
		func() Frame {
			var x X
			return x.ptr()
		}(),
		"%n",
		`\(\*X\).ptr`,
	}, {
		func() Frame {
			var x X
			return x.val()
		}(),
		"%n",
		"X.val",
	}, {
		Frame(0),
		"%n",
		"",
	}, {
		Frame(initpc),
		"%v",
		"stack_test.go:12",
	}, {
		Frame(initpc),
		"%+v",
		"github.com/hedzr/errors/canned.init\n" +
			"\t.+/github.com/hedzr/errors/canned/stack_test.go:12",
	}, {
		Frame(0),
		"%v",
		"unknown:0",
	}}

	for i, tt := range tests {
		testFormatRegexp(t, i, tt.Frame, tt.format, tt.want)
	}
}

func TestFuncname(t *testing.T) {
	tests := []struct {
		name, want string
	}{
		{"", ""},
		{"runtime.main", "main"},
		{"github.com/hedzr/errors/canned.funcname", "funcname"},
		{"funcname", "funcname"},
		{"io.copyBuffer", "copyBuffer"},
		{"main.(*R).Write", "(*R).Write"},
	}

	for _, tt := range tests {
		got := funcname(tt.name)
		want := tt.want
		if got != want {
			t.Errorf("funcname(%q): want: %q, got %q", tt.name, want, got)
		}
	}
}

func TestStackTrace(t *testing.T) {
	tests := []struct {
		err  error
		want []string
	}{{
		New("ooh"), []string{
			"github.com/hedzr/errors/canned.TestStackTrace\n" +
				"\t.+/github.com/hedzr/errors/canned/stack_test.go:154",
		},
	}, {
		errors.Wrap(errors.New("ooh"), "ahh"), []string{
			"github.com/hedzr/errors/canned.TestStackTrace\n" +
				"\t.+/github.com/hedzr/errors/canned/stack_test.go:159", // this is the stack of Wrap, not New
		},
	}, {
		errors.Cause(errors.Wrap(New("ooh"), "ahh")), []string{
			"github.com/hedzr/errors/canned.TestStackTrace\n" +
				"\t.+/github.com/hedzr/errors/canned/stack_test.go:164", // this is the stack of New
		},
	}, {
		func() error { return errors.New("ooh") }(), []string{
			`github.com/hedzr/errors/canned.(func·009|TestStackTrace.func1)` +
				"\n\t.+/github.com/hedzr/errors/canned/stack_test.go:169", // this is the stack of New
			"github.com/hedzr/errors/canned.TestStackTrace\n" +
				"\t.+/github.com/hedzr/errors/canned/stack_test.go:169", // this is the stack of New's caller
		},
	}, {
		errors.Cause(func() error {
			return func() error {
				return errors.Errorf("hello %s", fmt.Sprintf("world"))
			}()
		}()), []string{
			`github.com/hedzr/errors/canned.(func·010|TestStackTrace.func2.1)` +
				"\n\t.+/github.com/hedzr/errors/canned/stack_test.go:178", // this is the stack of Errorf
			`github.com/hedzr/errors/canned.(func·011|TestStackTrace.func2)` +
				"\n\t.+/github.com/hedzr/errors/canned/stack_test.go:179", // this is the stack of Errorf's caller
			"github.com/hedzr/errors/canned.TestStackTrace\n" +
				"\t.+/github.com/hedzr/errors/canned/stack_test.go:180", // this is the stack of Errorf's caller's caller
		},
	}}
	for i, tt := range tests {
		x, ok := tt.err.(interface {
			StackTrace() StackTrace
		})
		if !ok {
			t.Errorf("expected %#v to implement StackTrace() StackTrace", tt.err)
			continue
		}
		st := x.StackTrace()
		for j, want := range tt.want {
			testFormatRegexp(t, i, st[j], "%+v", want)
		}
	}
}

func stackTrace() StackTrace {
	const depth = 8
	var pcs [depth]uintptr
	n := runtime.Callers(1, pcs[:])
	var st Stack = pcs[0:n]
	return st.StackTrace()
}

func TestStackTraceFormat(t *testing.T) {
	tests := []struct {
		StackTrace
		format string
		want   string
	}{{
		nil,
		"%s",
		`\[\]`,
	}, {
		nil,
		"%v",
		`\[\]`,
	}, {
		nil,
		"%+v",
		"",
	}, {
		nil,
		"%#v",
		`\[\]errors.Frame\(nil\)`,
	}, {
		make(StackTrace, 0),
		"%s",
		`\[\]`,
	}, {
		make(StackTrace, 0),
		"%v",
		`\[\]`,
	}, {
		make(StackTrace, 0),
		"%+v",
		"",
	}, {
		make(StackTrace, 0),
		"%#v",
		`\[\]errors.Frame{}`,
	}, {
		stackTrace()[:2],
		"%s",
		`\[stack_test.go stack_test.go\]`,
	}, {
		stackTrace()[:2],
		"%v",
		`\[stack_test.go:207 stack_test.go:254\]`,
	}, {
		stackTrace()[:2],
		"%+v",
		"\n" +
			"github.com/hedzr/errors/canned.stackTrace\n" +
			"\t.+/github.com/hedzr/errors/canned/stack_test.go:207\n" +
			"github.com/hedzr/errors/canned.TestStackTraceFormat\n" +
			"\t.+/github.com/hedzr/errors/canned/stack_test.go:258",
	}, {
		stackTrace()[:2],
		"%#v",
		`\[\]errors.Frame{stack_test.go:207, stack_test.go:266}`,
	}}

	for i, tt := range tests {
		testFormatRegexp(t, i, tt.StackTrace, tt.format, tt.want)
	}
}

func testFormatRegexp(t *testing.T, n int, arg interface{}, format, want string) {
	got := fmt.Sprintf(format, arg)
	gotLines := strings.SplitN(got, "\n", -1)
	wantLines := strings.SplitN(want, "\n", -1)

	if len(wantLines) > len(gotLines) {
		t.Errorf("test %d: wantLines(%d) > gotLines(%d):\n got: %q\nwant: %q", n+1, len(wantLines), len(gotLines), got, want)
		return
	}

	for i, w := range wantLines {
		match, err := regexp.MatchString(w, gotLines[i])
		if err != nil {
			t.Fatal(err)
		}
		if !match {
			t.Errorf("test %d: line %d: fmt.Sprintf(%q, err):\n got: %q\nwant: %q", n+1, i+1, format, got, want)
		}
	}
}
