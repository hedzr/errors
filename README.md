# errors.v3

[![Go](https://github.com/hedzr/errors/actions/workflows/go.yml/badge.svg)](https://github.com/hedzr/errors/actions/workflows/go.yml)
[![GitHub tag (latest SemVer)](https://img.shields.io/github/tag/hedzr/errors.svg?label=release)](https://gopkg.in/hedzr/errors.v3)
[![GoDoc](https://img.shields.io/badge/godoc-reference-blue.svg?style=flat)](https://pkg.go.dev/gopkg.in/hedzr/errors.v3)
[![Go Report Card](https://goreportcard.com/badge/github.com/hedzr/errors)](https://goreportcard.com/report/github.com/hedzr/errors)
[![Coverage Status](https://coveralls.io/repos/github/hedzr/errors/badge.svg)](https://coveralls.io/github/hedzr/errors)
[![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2Fhedzr%2Ferrors.svg?type=shield)](https://app.fossa.com/projects/git%2Bgithub.com%2Fhedzr%2Ferrors?ref=badge_shield)

Wrapped errors and more for golang developing (not just for go1.13+).

`hedzr/errors` provides the compatbilities to your old project up to go 1.13.

`hedzr/errors` provides some extra enhancements for better context environment saving on error occurred.

## Import

```go
import "gopkg.in/hedzr/errors.v3"
```

## History

- v3.0.15
  - fix: make Is() work for go1.12 and below

- v3.0.13
  - fea: IsDescended for error template test
  - fix: code and causes are present at same time

- v3.0.11
  - changed the `FormatWith` interface to support error template feature.

- v3.0.10
  - support go1.11-17,18+

- v3.0.9
  - fix: WithErrors will check IsEmpty on an error container and avoid adding it if empty.

- v3.0.8
  - restore error message template
  - all features in v2 are restored with a new fluent style

- v3.0.7
  - fix coverall upload
  - godoc
  - better sites/taggedsites output

- v3.0.6
  - back to master branch

- v3.0.5
  - break out `New(...).Attach(...)`, instead of `New(...).WithErrors(...)`, so that we can make the type architecture clearly and concisely.
  - `Builable` and `Error` interface are the abstract representations about our error objects.
  - bugs fixed
  - more godoc

- v3.0.3
  - review the backward compatibilities

- v3.0.0
  - rewrite most codes and cleanup multiple small types
  - use `New(...)` or `NewBuilder()` to make an error with code, message, inner error(s) and customizable stacktrace info.

## Features

These features are supported for compatibilities.

#### stdlib `errors' compatibilities

- `func As(err error, target interface{}) bool`
- `func Is(err, target error) bool`
- `func New(text string) error`
- `func Unwrap(err error) error`

#### `pkg/errors` compatibilities

- `func Wrap(err error, message string) error`
- `func Cause(err error) error`: unwraps recursively, just like Unwrap()
- [x] `func Cause1(err error) error`: unwraps just one level
- `func WithCause(cause error, message string, args ...interface{}) error`, = `Wrap`
- supports Stacktrace
  - in an error by `Wrap()`, stacktrace wrapped;
  - for your error, attached by `WithStack(cause error)`;

#### Others

- Codes
- Inner errors  
  We like the flatter inner errors more than the cascade chain, so the `Format("%w)` is a so-so approach to collect the errors. We believe the error slice is a better choice.
- Unwrap inner errors one by one
- Error Template

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
    errContainer := errors.New("multiple tasks have errors")

    defer errContainer.Defer(&err)
    for _, r := range []error{io.EOF, io.ErrShortWrite, io.ErrClosedPipe, errors.Internal} {
      errContainer.Attach(r)
    }
    
    doWithItem := func(item Item) (err error) {
      // ...
      return
    }
    for _, item := range SomeItems {
      // nil will be ignored safely, do'nt worry about invalid attaches.
      errContainer.Attach(doWithItem(item))
    }

    return
  }

  err := child()
  t.Logf("failed: %+v", err)
}
```

### Error Template

We could *declare* a message template at first and format it with live args
to build an error instantly.

```go
func TestErrorsTmpl(t *testing.T) {
var errTmpl = errors.New("expecting %v but got %v")

	var err error
	err = errTmpl.FormatWith("789", "123")
	t.Logf("The error is: %v", err)
	err = errTmpl.FormatWith(true, false)
	t.Logf("The error is: %v", err)
}
```

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

## LICENSE

MIT

### Scan

[![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2Fhedzr%2Ferrors.svg?type=large)](https://app.fossa.com/projects/git%2Bgithub.com%2Fhedzr%2Ferrors?ref=badge_large)