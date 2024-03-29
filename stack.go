package errors

import (
	"fmt"
	"io"
	"path"
	"runtime"
	"strings"
)

// Frame represents a program counter inside a Stack frame.
type Frame uintptr

// pc returns the program counter for this frame;
// multiple frames may have the same PC value.
func (f Frame) pc() uintptr { return uintptr(f) - 1 }

// file returns the full path to the file that contains the
// function for this Frame's pc.
func (f Frame) file() string {
	fn := runtime.FuncForPC(f.pc())
	if fn == nil {
		return "unknown"
	}
	file, _ := fn.FileLine(f.pc())
	return file
}

// line returns the line number of source code of the
// function for this Frame's pc.
func (f Frame) line() int {
	fn := runtime.FuncForPC(f.pc())
	if fn == nil {
		return 0
	}
	_, line := fn.FileLine(f.pc())
	return line
}

// Format formats the frame according to the fmt.Formatter interface.
//
//	%s    source file
//	%d    source line
//	%n    function name
//	%v    equivalent to %s:%d
//
// Format accepts flags that alter the printing of some verbs, as follows:
//
//	%+s   function name and path of source file relative to the
//	      compiling time.
//	      GOPATH separated by \n\t (<funcname>\n\t<path>)
//	%+v   equivalent to %+s:%d
func (f Frame) Format(s fmt.State, verb rune) {
	switch verb {
	case 's':
		switch {
		case s.Flag('+'):
			pc := f.pc()
			fn := runtime.FuncForPC(pc)
			if fn == nil {
				_, _ = io.WriteString(s, "unknown")
			} else {
				file, _ := fn.FileLine(pc)
				_, _ = fmt.Fprintf(s, "%s\n\t%s", fn.Name(), file)
			}
		default:
			_, _ = io.WriteString(s, path.Base(f.file()))
		}
	case 'd':
		_, _ = fmt.Fprintf(s, "%d", f.line())
	case 'n':
		name := runtime.FuncForPC(f.pc()).Name()
		_, _ = io.WriteString(s, funcname(name))
	case 'v':
		f.Format(s, 's')
		_, _ = io.WriteString(s, ":")
		f.Format(s, 'd')
	}
}

// StackTrace is Stack of Frames from innermost (newest) to outermost (oldest).
type StackTrace []Frame

// Format formats the Stack of Frames according to the fmt.Formatter interface.
//
//	%s	lists source files for each Frame in the Stack
//	%v	lists the source file and line number for each Frame in the Stack
//
// Format accepts flags that alter the printing of some verbs, as follows:
//
//	%+v   Prints filename, function, and line number for each Frame in the Stack.
func (st StackTrace) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		switch {
		case s.Flag('+'):
			for _, f := range st {
				_, _ = fmt.Fprintf(s, "\n%+v", f)
			}
		case s.Flag('#'):
			_, _ = fmt.Fprintf(s, "%#v", []Frame(st))
		default:
			_, _ = fmt.Fprintf(s, "%v", []Frame(st))
		}
	case 's':
		_, _ = fmt.Fprintf(s, "%s", []Frame(st))
	}
}

// Stack represents a Stack of program counters.
type Stack []uintptr

// Format formats the stack of Frames according to the fmt.Formatter interface.
//
//	%s	lists source files for each Frame in the stack
//	%v	lists the source file and line number for each Frame in the stack
//
// Format accepts flags that alter the printing of some verbs, as follows:
//
//	%+v   Prints filename, function, and line number for each Frame in the stack.
func (s *Stack) Format(st fmt.State, verb rune) {
	if verb == 'v' && st.Flag('+') {
		for _, pc := range *s {
			f := Frame(pc)
			_, _ = fmt.Fprintf(st, "\n%+v", f)
		}
	}
}

// StackTrace returns the stacktrace frames
func (s *Stack) StackTrace() StackTrace {
	f := make([]Frame, len(*s))
	for i := 0; i < len(f); i++ {
		f[i] = Frame((*s)[i])
	}
	return f
}

func callers(skip int) *Stack {
	const depth = 32
	var pcs [depth]uintptr
	n := runtime.Callers(2+skip, pcs[:]) // by default, we skip these frames: callers(), and runtime.Callers()
	var st Stack = pcs[0:n]
	return &st
}

// funcname removes the path prefix component of a function's name reported by func.Name().
func funcname(name string) string {
	i := strings.LastIndex(name, "/")
	name = name[i+1:] //nolint:revive
	i = strings.Index(name, ".")
	return name[i+1:]
}
