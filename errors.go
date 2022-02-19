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

// Skip sets how many frames will be ignored while we are extracting the stacktrace info.
// Skip starts a builder with fluent API style, so you could continue
// build the error what you want:
//
//     err := errors.Skip(1).Message("hello %v", "you").Build()
//
func Skip(skip int) *builder {
	return &builder{skip: skip}
}

// Message formats a message and starts a builder to create the final error object.
//
//     err := errors.Message("hello %v", "you").Attach(causer).Build()
func Message(message string, args ...interface{}) *builder {
	return NewBuilder().WithMessage(message, args...)
}

func NewBuilder() *builder {
	return &builder{skip: 1}
}

type builder struct {
	skip    int
	causes2 causes2
}

func (s *builder) WithSkip(skip int) *builder {
	s.skip = skip
	return s
}

func (s *builder) WithErrors(errs ...error) *builder {
	_ = s.causes2.WithErrors(errs...)
	return s
}

func (s *builder) WithMessage(message string, args ...interface{}) *builder {
	if len(args) > 0 {
		message = fmt.Sprintf(message, args...)
	}
	s.causes2.msg = message
	return s
}

func (s *builder) WithCode(code Code) *builder {
	s.causes2.Code = code
	return s
}

func (s *builder) Build() *WithStackInfo {
	w := &WithStackInfo{
		causes2: s.causes2,
		Stack:   callers(s.skip),
	}
	return w
}
