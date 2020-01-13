# errors.v2

[![Build Status](https://travis-ci.org/hedzr/errors.svg?branch=master)](https://travis-ci.org/hedzr/errors)
[![GitHub tag (latest SemVer)](https://img.shields.io/github/tag/hedzr/errors.svg?label=release)](https://github.com/hedzr/errors/releases)
[![GoDoc](https://img.shields.io/badge/godoc-reference-blue.svg?style=flat)](https://godoc.org/github.com/hedzr/errors) 
[![Go Report Card](https://goreportcard.com/badge/github.com/hedzr/errors)](https://goreportcard.com/report/github.com/hedzr/errors)
[![codecov](https://codecov.io/gh/hedzr/errors/branch/master/graph/badge.svg)](https://codecov.io/gh/hedzr/errors)


Attachable errors for golang dev (for go1.13+).



## Import

```go
// wrong: import "github.com/hedzr/errors/v2"
import "gopkg.in/hedzr/errors.v2"
```
To take affect after new version released right away, delete the Go Modules local cache:

```bash
rm -rf $GOPATH/pkg/mod/*
```

Or, try go get the exact version just like:

```bash
go get -v gopkg.in/hedzr/errors.v2@v2.0.3
```




## Features




#### stdlib `errors' compatibilities

- `func As(err error, target interface{}) bool`
- `func Is(err, target error) bool`
- `func New(text string) error`
- `func Unwrap(err error) error`

#### `pkg/errors` compatibilities

- `func Wrap(err error, message string) error`
- `func Cause(err error) error`
- [x] `func Cause1(err error) error`
- `func WithCause(cause error, message string, args ...interface{}) error`, = `Wrap`
- supports Stacktrace
  - in an error by `Wrap()`, stacktrace wrapped;
  - for your error, attached by `WithStack(cause error)`;

#### enhancements

- `New(msg, args...)` combines New and `Newf`(if there is a name), WithMessage, WithMessagef, ...
- `WithCause(cause error, message string, args...interface{})`
- `Wrap(err error, message string, args ...interface{}) error`
- `DumpStacksAsString(allRoutines bool)`: returns stack tracing information like debug.PrintStack()
- `CanXXX`:
   - `CanAttach(err interface{}) bool`
   - `CanCause(err interface{}) bool`
   - `CanUnwrap(err interface{}) bool`
   - `CanIs(err interface{}) bool`
   - `CanAs(err interface{}) bool`
     



#### error Container and sub-errors (wrapped, attached or nested)

- `NewContainer(message string, args ...interface{}) *withCauses`
- `ContainerIsEmpty(container error) bool`
- `AttachTo(container *withCauses, errs ...error)`
- `withCauses.Attach(errs ...error)`

For example:

```go
func a() (err error){
    holder := errors.NewContainer("errors in a()")
    // ...
    for {
        // ...
        // errors.AttachTo(holder, io.EOF, io.ShortWrite)
        holder.Attach(io.EOF, io.ShortWrite)
    }
    err = holder.Error()
    return
}

func b(){
    err := a()
    if errors.Is(err, io.ShortWrite) {
        panic(err)
    }
}
```



#### Coded error

- `Code` is a generic type of error codes
- `WithCode(code, err, msg, args...)` can format an error object with error code, attached inner err, msg, and stack info.
- `Code.New(msg, args...)` is like `WithCode`.
- `Code.Register(codeNameString)` declares the name string of an error code yourself.

```go
err := InvalidArgument.New("wrong").Attach(io.ErrShortWrite)

const MyCode001 Code=1001

MyCode001.Register("MyCode001")
err := MyCode001.New("wrong 001: no config file")
```


## ACK

- stack.go is an copy from pkg/errors
- withStack is an copy from pkg/errors
- Is, As, Unwrap are inspired from go1.13 errors
- Cause, Wrap are inspired from pkg/errors

## LICENSE

MIT
