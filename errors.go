package errors

import "fmt"

// New returns an error with the supplied message.
// New also records the Stack trace at the point it was called.
func New(args ...interface{}) *WithStackInfo {
	s := &builder{skip: 1}

	if len(args) > 0 {
		if message, ok := args[0].(string); ok {
			return s.WithMessage(message, args[1:]...).Build()
		}
		for _, opt := range args {
			if o, ok := opt.(Opt); ok {
				o(s)
			}
		}
		return s.Build()
	}

	return &WithStackInfo{Stack: callers(1)}
}

// Opt _
type Opt func(s *builder)

// WithErrors _
func WithErrors(errs ...error) Opt {
	return func(s *builder) {
		s.WithErrors(errs...)
	}
}

// Skip sets how many frames will be ignored while we are extracting the stacktrace info.
// Skip starts a builder with fluent API style, so you could continue
// build the error what you want:
//
//     err := errors.Skip(1).Message("hello %v", "you").Build()
//
func Skip(skip int) Builder {
	return &builder{skip: skip}
}

// Message formats a message and starts a builder to create the final error object.
//
//     err := errors.Message("hello %v", "you").Attach(causer).Build()
func Message(message string, args ...interface{}) Builder {
	return NewBuilder().WithMessage(message, args...)
}

// NewBuilder starts a new error builder.
//
// Typically, you could make an error with fluent calls:
//
//    err = errors.NewBuilder().
//    	WithCode(Internal).
//    	WithErrors(io.EOF).
//    	WithErrors(io.ErrShortWrite).
//    	Build()
//    t.Logf("failed: %+v", err)
//
func NewBuilder() Builder {
	return &builder{skip: 1}
}

// Builder provides a fluent calling interface to make error building easy.
type Builder interface {
	// WithSkip specifies a special number of stack frames that will be ignored.
	WithSkip(skip int) Builder
	// WithErrors attaches the given errs as inner errors.
	WithErrors(errs ...error) Builder
	// WithMessage formats the error message
	WithMessage(message string, args ...interface{}) Builder
	// WithCode specifies an error code.
	WithCode(code Code) Builder
	// Build builds the final error object (with *WithStackInfo type wrapped)
	Build() *WithStackInfo

	// Attach inner errors for backward compatibility to v2
	Attach(errs ...error)
}

type builder struct {
	skip    int
	causes2 causes2
}

// WithSkip specifies a special number of stack frames that will be ignored.
func (s *builder) WithSkip(skip int) Builder {
	s.skip = skip
	return s
}

// WithErrors attaches the given errs as inner errors.
func (s *builder) WithErrors(errs ...error) Builder {
	_ = s.causes2.WithErrors(errs...)
	return s
}

// Attach attaches the given errs as inner errors.
// For backward compatibility to v2
func (s *builder) Attach(errs ...error) {
	_ = s.WithErrors(errs...).Build()
}

// WithMessage formats the error message
func (s *builder) WithMessage(message string, args ...interface{}) Builder {
	if len(args) > 0 {
		message = fmt.Sprintf(message, args...)
	}
	s.causes2.msg = message
	return s
}

// WithCode specifies an error code.
func (s *builder) WithCode(code Code) Builder {
	s.causes2.Code = code
	return s
}

// Build builds the final error object (with *WithStackInfo type wrapped)
func (s *builder) Build() *WithStackInfo {
	w := &WithStackInfo{
		causes2: s.causes2,
		Stack:   callers(s.skip),
	}
	return w
}
