# errors

nestable errors for golang dev.

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

### register your error codes:

The user-defined error codes could be registered into `errors.Code` with its codename.

For example (run it at paly-ground: https://play.golang.org/p/N-P7lqdJPzy):

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



## LICENSE

MIT
