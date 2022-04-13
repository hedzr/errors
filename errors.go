package errors

import "fmt"

// New returns an error with the supplied message.
// New also records the Stack trace at the point where it was called.
func New(args ...interface{}) Error {
	s := &builder{skip: 1}

	if len(args) > 0 {
		if message, ok := args[0].(string); ok {
			return s.WithSkip(2).WithMessage(message, args[1:]...).Build()
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

// WithErrors attach child errors into an error container.
// For a container which has IsEmpty() interface, it would not be attach if it is empty (i.e. no errors).
// For a nil error object, it will be ignored.
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
	// For a container which has IsEmpty() interface, it would not be attach if it is empty (i.e. no errors).
	// For a nil error object, it will be ignored.
	WithErrors(errs ...error) Builder
	// WithMessage formats the error message
	WithMessage(message string, args ...interface{}) Builder
	// WithCode specifies an error code.
	WithCode(code Code) Builder

	// Build builds the final error object (with Buildable interface bound)
	Build() Error

	// BREAK - Use WithErrors() for instead
	// Attach inner errors for backward compatibility to v2
	// Attach(errs ...error)
}

type builder struct {
	skip        int
	causes2     causes2
	sites       []interface{}
	taggedSites TaggedData
}

// WithSkip specifies a special number of stack frames that will be ignored.
func (s *builder) WithSkip(skip int) Builder {
	s.skip = skip
	return s
}

// WithCode specifies an error code.
func (s *builder) WithCode(code Code) Builder {
	s.causes2.Code = code
	return s
}

// // Attach attaches the given errs as inner errors.
// // For backward compatibility to v2
// func (s *builder) Attach(errs ...error) Buildable {
//	return s.WithErrors(errs...).Build()
// }

// WithMessage formats the error message
func (s *builder) WithMessage(message string, args ...interface{}) Builder {
	if len(args) > 0 {
		message = fmt.Sprintf(message, args...)
	}
	s.causes2.msg = message
	return s
}

// WithErrors attaches the given errs as inner errors.
// For a container which has IsEmpty() interface, it would not be attach if it is empty (i.e. no errors).
// For a nil error object, it will be ignored.
func (s *builder) WithErrors(errs ...error) Builder {
	_ = s.causes2.WithErrors(errs...)
	return s
}

// WithData appends errs if the general object is a error object.
// It can be used in defer-recover block typically. For example:
//
//    defer func() {
//      if e := recover(); e != nil {
//        err = errors.New("[recovered] copyTo unsatisfied ([%v] %v -> [%v] %v), causes: %v",
//          c.indirectType(from.Type()), from, c.indirectType(to.Type()), to, e).
//          WithData(e)
//        n := log.CalcStackFrames(1)   // skip defer-recover frame at first
//        log.Skip(n).Errorf("%v", err) // skip go-lib frames and defer-recover frame, back to the point throwing panic
//      }
//    }()
//
func (s *builder) WithData(errs ...interface{}) Builder {
	s.sites = append(s.sites, errs...)
	return s
}

// WithTaggedData appends user data with tag into internal container.
// These data can be retrieved by
func (s *builder) WithTaggedData(siteScenes TaggedData) Builder {
	if s.taggedSites == nil {
		s.taggedSites = make(TaggedData)
	}
	for k, v := range siteScenes {
		s.taggedSites[k] = v
	}
	return s
}

// WithCause sets the underlying error manually if necessary.
func (s *builder) WithCause(cause error) Builder {
	_ = s.causes2.WithErrors(cause)
	return s
}

// Build builds the final error object (with *WithStackInfo type wrapped)
func (s *builder) Build() Error {
	w := &WithStackInfo{
		causes2: s.causes2,
		Stack:   callers(s.skip),
	}
	return w
}
