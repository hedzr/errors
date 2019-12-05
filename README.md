# errors

[![Build Status](https://travis-ci.org/hedzr/errors.svg?branch=master)](https://travis-ci.org/hedzr/errors)
[![GitHub tag (latest SemVer)](https://img.shields.io/github/tag/hedzr/errors.svg?label=release)](https://github.com/hedzr/errors/releases)
[![GoDoc](https://img.shields.io/badge/godoc-reference-blue.svg?style=flat)](https://godoc.org/github.com/hedzr/errors) 
[![Go Report Card](https://goreportcard.com/badge/github.com/hedzr/errors)](https://goreportcard.com/report/github.com/hedzr/errors)
[![codecov](https://codecov.io/gh/hedzr/errors/branch/master/graph/badge.svg)](https://codecov.io/gh/hedzr/errors)


Nestable errors for golang dev.

Take a look at: https://play.golang.org/p/bsGjRAWJDOA

## Import

```go
import "gtihub.com/hedzr/errors"
```

## ExtErr

### `error` with message

```go
var (
  errBug1 = errors.New("bug 1")
  errBug2 = errors.New("bug 2")
  errBug3 = errors.New("bug 3")
)

func main() {
  err := errors.New("something", errBug1, errBug2)
  err2 := errors.New(errBug3)

  log.Println(err, err2)
}
```

## CodedErr

### `error` with a code

```go
var(
  errNotFound = errors.NewWithCode(errors.NotFound)
  errNotFoundMsg = errors.NewWithCodeMsg(errors.NotFound, "not found")
)
```

### Predefined error codes

The builtin error codes are copied from Google gRPC codes but negatived.

**The numbers -1..-999 are reserved.**


### register your error codes:

The user-defined error codes (must be < -1000, or > 0) could be registered into `errors.Code` with its codename.

For example (run it at play-ground: https://play.golang.org/p/N-P7lqdJPzy):

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
	errBug1001 = errors.NewWithCodeMsg(BUG1001, "something is wrong", io.EOF)
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


## replacement of go `errors`

As a replacement, the functions are copied from go `errors`, such as:

- `Is(err, target) bool`
- `As(err, target) bool`
- `Unwrap(err) error`

## LICENSE

MIT
