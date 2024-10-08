# errors.v3

[![Go](https://github.com/hedzr/errors/actions/workflows/go.yml/badge.svg)](https://github.com/hedzr/errors/actions/workflows/go.yml)
[![GitHub tag (latest SemVer)](https://img.shields.io/github/tag/hedzr/errors.svg?label=release)](https://gopkg.in/hedzr/errors.v3)
[![GoDoc](https://img.shields.io/badge/godoc-reference-blue.svg?style=flat)](https://pkg.go.dev/gopkg.in/hedzr/errors.v3)
[![Go Report Card](https://goreportcard.com/badge/github.com/hedzr/errors)](https://goreportcard.com/report/github.com/hedzr/errors)
[![Coverage Status](https://coveralls.io/repos/github/hedzr/errors/badge.svg)](https://coveralls.io/github/hedzr/errors)
[![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2Fhedzr%2Ferrors.svg?type=shield)](https://app.fossa.com/projects/git%2Bgithub.com%2Fhedzr%2Ferrors?ref=badge_shield)

Wrapped errors and more for golang developing (not just for go1.11, go1.13, and go1.20+).

`hedzr/errors` provides the compatibilities to your old project up to go 1.20.

`hedzr/errors` provides some extra enhancements for better context environment saving on error occurred.

## Features

- Simple migrating way from std errors: all of standard functions have been copied to
- Better `New()`:
  - format message inline: `err := errors.New("hello %s", "world")`
  - format with `WithXXX`: `err := errors.New(errors.WithErrors(errs...))`
  - cascade format: `err := errors.New().WithErrors(errs...)`
  - Stacktrace awareness
  - Container for canning errors: [Error Container (Inner/Nested)](#error-container-innernested)
  - error template: [Format message instantly but the text template can be given at beginning](#error-template)
- Codes: treat a number as an error object
- Unwrap inner canned errors one by one
- No mental burden

### Others

## History

- v3.3.3
  - ensure working on go 1.23+

- v3.3.2
  - fixed/improved Attach() - attaching itself is denied now

- v3.3.1
  - fixed Iss() couldn't test the others except the first error.

- v3.3.0
  - added `Iss(err, errs...)` to test if any of errors are included in 'err'.
  - improved As/Is/Unwrap to fit for new joint error since go1.20
  - added causes2.Clear, ...
  - improved Error() string
  - reviewed and re-published this repo from v3.3

- v3.1.9
  - fixed error.Is deep test to check two errors' message text contents if matched
  - fixed errors.v3.Join when msg is not empty in an err obj
  - fixed causes.WithErrors(): err obj has been ignored even if its message is not empty

- OLDER in [CHANGELOG](https://github.com/hedzr/errors/blob/master/CHANGELOG)

## Compatibilities

These features are supported for compatibilities.

### stdlib `errors' compatibilities

- `func As(err error, target interface{}) bool`
- `func Is(err, target error) bool`
- `func New(text string) error`
- `func Unwrap(err error) error`
- `func Join(errs ...error) error`

### `pkg/errors` compatibilities

- `func Wrap(err error, message string) error`
- `func Cause(err error) error`: unwraps recursively, just like Unwrap()
- [x] `func Cause1(err error) error`: unwraps just one level
- `func WithCause(cause error, message string, args ...interface{}) error`, = `Wrap`
- supports Stacktrace
  - in an error by `Wrap()`, stacktrace wrapped;
  - for your error, attached by `WithStack(cause error)`;

### Some Enhancements

- `Iss(err error, errs ...error) bool`
- `AsSlice(errs []error, target interface{}) bool`
- `IsAnyOf(err error, targets ...error) bool`
- `IsSlice(errs []error, target error) bool`
- `TypeIs(err, target error) bool`
- `TypeIsSlice(errs []error, target error) bool`
- `Join(errs ...error) error`
- `DumpStacksAsString(allRoutines bool) string`
- `CanAttach(err interface{}) (ok bool)`
- `CanCause(err interface{}) (ok bool)`
- `CanCauses(err interface{}) (ok bool)`
- `CanUnwrap(err interface{}) (ok bool)`
- `CanIs(err interface{}) (ok bool)`
- `CanAs(err interface{}) (ok bool)`
- `Causes(err error) (errs []error)`

## Best Practices

### Basics

```go
package test

import (
    "gopkg.in/hedzr/errors.v3"
    "io"
    "reflect"
    "testing"
)

func TestForExample(t *testing.T) {
  fn := func() (err error) {
    ec := errors.New("some tips %v", "here")
    defer ec.Defer(&err)

    // attaches much more errors
    for _, e := range []error{io.EOF, io.ErrClosedPipe} {
      ec.Attach(e)
    }
  }

  err := fn()
  t.Logf("failed: %+v", err)

  // use another number different to default to skip the error frames
  err = errors.
        Skip(3). // from on Skip()
        WithMessage("some tips %v", "here").Build()
  t.Logf("failed: %+v", err)

  err = errors.
        Message("1"). // from Message() on
        WithSkip(0).
        WithMessage("bug msg").
        Build()
  t.Logf("failed: %+v", err)

  err = errors.
        NewBuilder(). // from NewBuilder() on
        WithCode(errors.Internal). // add errors.Code
        WithErrors(io.EOF). // attach inner errors
        WithErrors(io.ErrShortWrite, io.ErrClosedPipe).
        Build()
  t.Logf("failed: %+v", err)

  // As code
  var c1 errors.Code
  if errors.As(err, &c1) {
    println(c1) // = Internal
  }

  // As inner errors
  var a1 []error
  if errors.As(err, &a1) {
    println(len(a1)) // = 3, means [io.EOF, io.ErrShortWrite, io.ErrClosedPipe]
  }
  // Or use Causes() to extract them:
  if reflect.DeepEqual(a1, errors.Causes(err)) {
    t.Fatal("unexpected problem")
  }

  // As error, the first inner error will be extracted
  var ee1 error
  if errors.As(err, &ee1) {
    println(ee1) // = io.EOF
  }

  series := []error{io.EOF, io.ErrShortWrite, io.ErrClosedPipe, errors.Internal}
  var index int
  for ; ee1 != nil; index++ {
    ee1 = errors.Unwrap(err) // extract the inner errors one by one
    if ee1 != nil && ee1 != series[index] {
      t.Fatalf("%d. cannot extract '%v' error with As(), ee1 = %v", index, series[index], ee1)
    }
  }
}
```

### Error Container (Inner/Nested)

```go
func TestContainer(t *testing.T) {
  // as a inner errors container
  child := func() (err error) {
    ec := errors.New("multiple tasks have errors")
    defer ec.Defer(&err) // package the attached errors as a new one and return it as `err`

    for _, r := range []error{io.EOF, io.ErrShortWrite, io.ErrClosedPipe, errors.Internal} {
      ec.Attach(r)
    }
    
    doWithItem := func(item Item) (err error) {
      // ...
      return
    }
    for _, item := range SomeItems {
      // nil will be ignored safely, do'nt worry about invalid attaches.
      ec.Attach(doWithItem(item))
    }

    return
  }

  err := child() // get the canned errors as a packaged one
  t.Logf("failed: %+v", err)
}
```

### Error Template

We could *declare* a message template at first and format it with live args
to build an error instantly.

```go
func TestErrorsTmpl(t *testing.T) {
  errTmpl := errors.New("expecting %v but got %v")

  var err error
  err = errTmpl.FormatWith("789", "123")
  t.Logf("The error is: %v", err)
  err = errTmpl.FormatWith(true, false)
  t.Logf("The error is: %v", err)
}
```

`FormatWith` will make new clone from errTmpl so you can use multiple cloned errors thread-safely.

The derived error instance is the descendant of the error template.
This relation can be tested by `errors.IsDescent(errTempl, err)`

```go
func TestIsDescended(t *testing.T) {
  err3 := New("any error tmpl with %v")
  err4 := err3.FormatWith("huahua")
  if !IsDescended(err3, err4) {
    t.Fatalf("bad test on IsDescended(err3, err4)")
  }
}
```

### Better format for a nested error

Since v3.1.1, the better message format will be formatted at Printf("%+v").

```go
func TestAs_betterFormat(t *testing.T) {
  var err = New("Have errors").WithErrors(io.EOF, io.ErrShortWrite, io.ErrNoProgress)
  t.Logf("%v\n", err)
  
  var nestNestErr = New("Errors FOUND:").WithErrors(err, io.EOF)
  var nnnErr = New("Nested Errors:").WithErrors(nestNestErr, strconv.ErrRange)
  t.Logf("%v\n", nnnErr)
  t.Logf("%+v\n", nnnErr)
}
```

The output is:

```bash
=== RUN   TestAs_betterFormat
    causes_test.go:23: Have errors [EOF | short write | multiple Read calls return no data or error]
    causes_test.go:27: Nested Errors: [Errors FOUND: [Have errors [EOF | short write | multiple Read calls return no data or error] | EOF] | value out of range]
    causes_test.go:28: Nested Errors:
          - Errors FOUND:
            - Have errors
              - EOF
              - short write
              - multiple Read calls return no data or error
            - EOF
          - value out of range
        
        gopkg.in/hedzr/errors%2ev3.TestAs_betterFormat
          /Volumes/VolHack/work/godev/cmdr-series/libs/errors/causes_test.go:26
        testing.tRunner
          /usr/local/go/src/testing/testing.go:1576
        runtime.goexit
          /usr/local/go/src/runtime/asm_amd64.s:1598
--- PASS: TestAs_betterFormat (0.00s)
PASS
```

## LICENSE

MIT

### Scan

[![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2Fhedzr%2Ferrors.svg?type=large)](https://app.fossa.com/projects/git%2Bgithub.com%2Fhedzr%2Ferrors?ref=badge_large)
