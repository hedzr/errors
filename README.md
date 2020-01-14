# errors

[![Build Status](https://travis-ci.org/hedzr/errors.svg?branch=master)](https://travis-ci.org/hedzr/errors)
[![GitHub tag (latest SemVer)](https://img.shields.io/github/tag/hedzr/errors.svg?label=release)](https://github.com/hedzr/errors/releases)
[![GoDoc](https://img.shields.io/badge/godoc-reference-blue.svg?style=flat)](https://godoc.org/github.com/hedzr/errors) 
[![Go Report Card](https://goreportcard.com/badge/github.com/hedzr/errors)](https://goreportcard.com/report/github.com/hedzr/errors)
[![codecov](https://codecov.io/gh/hedzr/errors/branch/master/graph/badge.svg)](https://codecov.io/gh/hedzr/errors)


Nestable errors for golang dev (both go1.13+ and lower now).

Take a look at: <https://play.golang.org/p/P0kk4NhAbd3>



## Import

### v1 (Archived for legacy projects)

```go
import "github.com/hedzr/errors"
```


### v2

The new [`v2` branch](https://github.com/hedzr/errors/tree/v2) is cleaning rewroten version, preview at `v2.0.x` (*gopkg v2*).

```go
import "gopkg.in/hedzr/errors.v2"
```





## Enh for go `errors`

### Adapted for golang 1.13+

- `Is(err, target) bool`
- `As(err, target) bool`
- `Unwrap(err) error`

```go
func TestIsAs(t *testing.T) {
	var err error
	err = errors.New("something").Attach(errBug1, errBug2).Nest(errBug3, errBug4).Msg("anything")
	if errors.Is(err, errBug1) {
		fmt.Println(" err ==IS== errBug1")
	}
	if errors.Is(err, errBug3) {
		fmt.Println(" err ==IS== errBug3")
	}

	err2 := errors.NewWithError(io.ErrShortWrite).Nest(io.EOF)
	if errors.Is(err2, io.ErrShortWrite) {
		fmt.Println(" err2 ==IS== io.ErrShortWrite")
	}
	if errors.Is(err2, io.EOF) {
		fmt.Println(" err2 ==IS== io.EOF")
	}
}
```

> since v1.1.7, we works for golang 1.12 and lower.



### Replace stdlib `errors`

In most cases, smoothly migrated from stdlib `errors` is possible: just replace the import statements.

Adapted and enhanced:

- `New(fmt, ...)`

More extendings:

- `NewTemplate(tmpl)`
- `NewWithError(errs...)`
- `NewCodedError(code, errs...)`



### enhancements for go `errors`

1. Walkable: `errors.Walk(fn)`
2. Ranged: `errors.Range(fn)`
3. Tests:
   1. `CanWalk(err)`, `CanRange(err)`, `CanIs(err)`, `CanAs(err)`, `CanUnwrap(err)`
   2. `Equal(err, code)`, `IsAny(err, codes...)`, `IsBoth(err, codes...)`, 
   3. `TextContains(err, text)`
   4. `HasAttachedErrors(err)`
   5. `HasWrappedError(err)`
   6. `HasInnerErrors(err)`
4. `Attach(err, errs...)`, `Nest(err, errs...)`
5. `DumpStacksAsString(allRoutines bool) string`


### for `pkg/errors`

- `Wrap(err, format, args...)`
- `Wrapf(err, format, args...)`
- `WithMessage(err, msg)`
- `WithStack(err)`



## Err Objects

### `ExtErr` object

#### `error` with message

```go
var (
  errBug1 = errors.New("bug 1")
  errBug2 = errors.New("bug 2, %v, %d", []string{"a","b"}, 5)
  errBug3 = errors.New("bug 3")
)

func main() {
  err := errors.New("something grouped").Attach(errBug1, errBug2)
  err1 := errors.New("something nested").Nest(errBug1, errBug2)
  err2 := errors.New(errBug3)

  err2.Msg("a %d", 1)

  log.Println(err, err1, err2)
}
```

#### Attachable, and Nestable (Wrapable)

##### Attach

`Attach(...)` can package a group of errors into the receiver `ExtErr`.

For example:

```go
// or: errors.New("1").Attach(io.EOF,io.ErrShortWrite, io.ErrShortBuffer)
err := errors.New("1").Attach(io.EOF).Attach(io.ErrShortWrite).Attach(io.ErrShortBuffer)
fmt.Println(err)
fmt.Printf("%#v\n", err)
// result:
// 1, EOF, short write, short buffer
// &errors.ExtErr{inner:(*errors.ExtErr)(nil), errs:[]error{(*errors.errorString)(0xc000042040), (*errors.errorString)(0xc000042020), (*errors.errorString)(0xc000042030)}, msg:"1", tmpl:""}
```

The structure looks like:

```
&ExtErr{
  msg: "1",
  errs: [
    io.EOF,
    io.ErrShortWrite,
    io.ErrShortBuffer
  ],
}
```

##### Nest

`Nest(...)` could wrap the errors as a deep descendant child `ExtErr` of the receiver `ExtError`.

For example:

```go
err := errors.New("1").Nest(io.EOF).Nest(io.ErrShortWrite).Nest(io.ErrShortBuffer)
fmt.Println(err)
fmt.Printf("%#v\n", err)
// result:
// 1, EOF[error, short write[error, short buffer]]
// &errors.ExtErr{inner:(*errors.ExtErr)(0x43e2e0), errs:[]error{(*errors.errorString)(0x40c040)}, msg:"1", tmpl:""}
```

To make it clear:

```
&ExtErr{
  msg: "1",
  inner: &ExtErr {
    inner: &ExtErr {
      errs: [
        io.ErrShortBuffer
      ],
    },
    errs: [
      io.ErrShortWrite,
    ],
  },
  errs: [
    io.EOF,
  ],
}
```







#### Template

You could put a string template into both `ExtErr` and `CodedErr`, and format its till using.

For example:

<details>

```go
const (
	BUG1004 errors.Code = -1004
	BUG1005 errors.Code = -1005
)

var (
	eb1  = errors.NewTemplate("xxbug 1, cause: %v")
	eb11 = errors.New("").Msg("first, %v", "ok").Template("xxbug11, cause")
	eb2  = errors.NewCodedError(BUG1004).Template("xxbug4, cause: %v")
	eb3  = errors.NewCodedError(BUG1005).Template("xxbug5, cause: none")
	eb31 = errors.NewCodedError(BUG1004).Msg("first, %v", "ok").Template("xxbug4.31, cause: %v")
	eb4  = errors.NewCodedError(BUG1005).Template("xxbug54, cause: none")
)

func init() {
	BUG1004.Register("BUG1004")
	BUG1005.Register("BUG1005")
}

func TestAll(t *testing.T) {
	err = eb1.Formatf("resources exhausted")
	t.Log(err)
}
```

Yet another one:

```go
const (
	ErrNoNotFound errors.Code = -9710
)

var (
	ErrNotFound  = errors.NewCodedError(ErrNoNotFound).Template("'%v' not found")
	ErrNotFound2 = errors.NewCodedError(-9711).Template("'%v' not found")
)

func init() {
	ErrNoNotFound.Register("Not Found")
	errors.Code(-9711).Register("Not Found 2")
}

// ...
return ErrNoNotFound.Formatf(filename)
```

</details>







### `CodedErr` object

#### `error` with a code number

```go
var (
  errNotFound = errors.NewCodedError(errors.NotFound)
  errNotFoundMsg = errors.NewCodedError(errors.NotFound).Msg("not found")
)
```

#### Predefined error codes

The builtin error codes are copied from Google gRPC codes but negatived.

**The numbers -1..-999 are reserved.**

#### register your error codes:

The user-defined error codes (must be < -1000, or > 0) could be registered into `errors.Code` with its codename.

For example (run it at play-ground: https://play.golang.org/p/ifUvABaPEoJ):

```go
package main

import (
	"fmt"
	"github.com/hedzr/errors"
	"io"
)

const (
	BUG1001 errors.Code = 1001
	BUG1002 errors.Code = 1002
)

var (
	errBug1001 = errors.NewCodedError(BUG1001).Msg("something is wrong").Attach(io.EOF)
)

func init() {
	BUG1001.Register("BUG1001")
	BUG1002.Register("BUG1002")
}

func main() {
	fmt.Println(BUG1001.String())
	fmt.Println(BUG1001.Number())
	fmt.Println(errBug1001)
	fmt.Println(errBug1001.Equal(BUG1001))
	fmt.Println(errBug1001.EqualRecursive(BUG1001))
}
```

Result:

```
BUG1001
001001|BUG1001|something is wrong, EOF
```







## Extending `ExtErr` or `CodedErr`

You might want to extend the `ExtErr`/`CodedErr` with more fields, just like `CodedErr.code`.

For better cascaded calls, it might be crazy had you had to 
override some functions: `Template`, `Format`, `Msg`, `Attach`, and 
`Nest`. But no more worries, simply copy them from `CodedErr` 
and correct the return type for yourself.

For example:

<details>

```go
// newError formats a ErrorForCmdr object
func newError(ignorable bool, sourceTemplate *ErrorForCmdr, args ...interface{}) *ErrorForCmdr {
	e := sourceTemplate.Format(args...)
	e.Ignorable = ignorable
	return e
}

// newErrorWithMsg formats a ErrorForCmdr object
func newErrorWithMsg(msg string, inner error) *ErrorForCmdr {
	return newErr(msg).Attach(inner)
}

func newErr(msg string, args ...interface{}) *ErrorForCmdr {
	return &ErrorForCmdr{ExtErr: *errors.New(msg, args...)}
}

func newErrTmpl(tmpl string) *ErrorForCmdr {
	return &ErrorForCmdr{ExtErr: *errors.NewTemplate(tmpl)}
}

func (e *ErrorForCmdr) Error() string {
	var buf bytes.Buffer
	buf.WriteString(fmt.Sprintf("%v|%s", e.Ignorable, e.ExtErr.Error()))
	return buf.String()
}

// Template setup a string format template.
// Coder could compile the error object with formatting args later.
//
// Note that `ExtErr.Template()` had been overrided here
func (e *ErrorForCmdr) Template(tmpl string) *ErrorForCmdr {
	_ = e.ExtErr.Template(tmpl)
	return e
}

// Format compiles the final msg with string template and args
//
// Note that `ExtErr.Template()` had been overridden here
func (e *ErrorForCmdr) Format(args ...interface{}) *ErrorForCmdr {
	_ = e.ExtErr.Format(args...)
	return e
}

// Msg encodes a formattable msg with args into ErrorForCmdr
//
// Note that `ExtErr.Template()` had been overridden here
func (e *ErrorForCmdr) Msg(msg string, args ...interface{}) *ErrorForCmdr {
	_ = e.ExtErr.Msg(msg, args...)
	return e
}

// Attach attaches the nested errors into ErrorForCmdr
//
// Note that `ExtErr.Template()` had been overridden here
func (e *ErrorForCmdr) Attach(errors ...error) *ErrorForCmdr {
	_ = e.ExtErr.Attach(errors...)
	return e
}

// Nest attaches the nested errors into ErrorForCmdr
//
// Note that `ExtErr.Template()` had been overridden here
func (e *ErrorForCmdr) Nest(errors ...error) *ErrorForCmdr {
	_ = e.ExtErr.Nest(errors...)
	return e
}
```

</details>

A sample here: https://github.com/hedzr/errors-for-mqtt



## More

strings helpers:

- `Left(s, length)`
- `Right(s, length)`
- `LeftPad(s, pad, width)`





## LICENSE

MIT
