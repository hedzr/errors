package errors

import "fmt"

// NewContainer wraps a group of errors and msg as one and return it.
// The returned error object is a container to hold many sub-errors.
//
// Examples:
//
//
//
func NewContainer(message string, args ...interface{}) *withCauses {
	if len(args) > 0 {
		message = fmt.Sprintf(message, args...)
	}
	err := &withCauses{
		msg:   message,
		stack: callers(),
	}
	return err
}

// ContainerIsEmpty appends more errors into 'container' error container.
func ContainerIsEmpty(container error) bool {
	if x, ok := container.(interface{ IsEmpty() bool }); ok {
		return x.IsEmpty()
	}
	return false
}

// AttachTo appends more errors into 'container' error container.
func AttachTo(container *withCauses, errs ...error) {
	if container == nil {
		panic("nil error container not allowed")
	}
	container.Attach(errs...)
}
