# errors

[![Build Status](https://travis-ci.org/hedzr/errors.svg?branch=master)](https://travis-ci.org/hedzr/errors)
[![GitHub tag (latest SemVer)](https://img.shields.io/github/tag/hedzr/errors.svg?label=release)](https://github.com/hedzr/errors/releases)
[![GoDoc](https://img.shields.io/badge/godoc-reference-blue.svg?style=flat)](https://godoc.org/github.com/hedzr/errors) 
[![Go Report Card](https://goreportcard.com/badge/github.com/hedzr/errors)](https://goreportcard.com/report/github.com/hedzr/errors)
[![codecov](https://codecov.io/gh/hedzr/errors/branch/master/graph/badge.svg)](https://codecov.io/gh/hedzr/errors)


Nestable errors for golang dev.

Take a look at: <https://play.golang.org/p/Yt-9dCSHX1Z>

## Import

```go
import "gtihub.com/hedzr/errors"
```

## `ExtErr` object

### `error` with message

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

### Attachable, and Nestable

#### Attach

`Attach(...)` could wrap a group of errors into the receiver `ExtErr`.

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

The structure is:

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

#### Nest

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



## `CodedErr` object

### `error` with a code

```go
var(
  errNotFound = errors.NewCodedError(errors.NotFound)
  errNotFoundMsg = errors.NewCodedError(errors.NotFound).Msg("not found")
)
```

### Predefined error codes

The builtin error codes are copied from Google gRPC codes but negatived.

**The numbers -1..-999 are reserved.**


### register your error codes:

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
	fmt.Println(errBug1001)
}
```

Result:

```
BUG1001
001001|BUG1001|something is wrong, EOF
```

### Extending from `ExtErr`

You might want to extend from `ExtErr` with more fields, just like `CodedErr.code`.

For better cascaded calls, it might be crazy had you had to 
override some functions: `Template`, `Format`, `Msg`, `Attach`, and 
`Nest`. But no more worries, simply copy them from `CodedErr` 
and correct the return type for yourself.

## Template

You could put a string template into both `ExtErr` and `CodedErr`, and format its till using:

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
	err = eb1.Format("resources exhausted")
	t.Log(err)
}
```

Another sample:

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
return ErrNoNotFound.Format(filename)
```

## replacement of go `errors`

Adapted for golang 1.13:

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

## LICENSE

MIT
