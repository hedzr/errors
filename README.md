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
  err := errors.New("something").Attach(errBug1, errBug2)
  err2 := errors.New(errBug3)

  log.Println(err, err2)
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

As a replacement, the functions are copied from go `errors`, such as:

- `Is(err, target) bool`
- `As(err, target) bool`
- `Unwrap(err) error`

## LICENSE

MIT
