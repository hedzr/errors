package errors

// Error object
type Error interface {
	// Buildable _
	Buildable

	// Data returns the wrapped common user data by Buildable.WithData.
	// The error objects with passed Buildable.WithData will be moved
	// into inner errors set, so its are excluded from Data().
	Data() []interface{}
	// TaggedData returns the wrapped tagged user data by
	// Buildable.WithTaggedData.
	TaggedData() TaggedData
	// Cause returns the underlying cause of the error, if possible.
	// An error value has a cause if it implements the following
	// interface:
	//
	//     type causer interface {
	//            Cause() error
	//     }
	//
	// If an error object does not implement Cause interface, the
	// original error object will be returned.
	// If the error is nil, nil will be returned without further
	// investigation.
	Cause() error
	// Causes simply returns the wrapped inner errors.
	// It doesn't consider an wrapped Code entity is an inner error too.
	// So if you wanna to extract any inner error objects, use
	// errors.Unwrap for instead. The errors.Unwrap could extract all
	// of them one by one:
	//
	//      var err = errors.New("hello").WithErrors(io.EOF, io.ShortBuffers)
	//      var e error = err
	//      for e != nil {
	//          e = errors.Unwrap(err)
	//      }
	//
	Causes() []error
}

// Buildable provides a fluent calling interface to make error building easy.
// Buildable is an error interface too.
type Buildable interface {
	// error interface
	error

	// WithSkip specifies a special number of stack frames that will be ignored.
	WithSkip(skip int) Buildable
	// WithMessage formats the error message
	WithMessage(message string, args ...interface{}) Buildable
	// WithCode specifies an error code.
	// An error code `Code` is a integer number with error interface
	// supported.
	WithCode(code Code) Buildable
	// WithErrors attaches the given errs as inner errors.
	// WithErrors is like our old Attach().
	// It wraps the inner errors into underlying container and
	// represents them all in a singular up-level error object.
	// The wrapped inner errors can be retrieved with errors.Causes:
	//
	//      var err = errors.New("hello").WithErrors(io.EOF, io.ShortBuffers)
	//      var errs []error = errors.Causes(err)
	//
	// Or, use As() to extract its:
	//
	//      var errs []error
	//      errors.As(err, &errs)
	//
	// Or, use Unwrap() for its:
	//
	//      var e error = err
	//      for e != nil {
	//          e = errors.Unwrap(err)
	//      }
	//
	WithErrors(errs ...error) Buildable
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
	WithData(errs ...interface{}) Buildable
	// WithTaggedData appends user data with tag into internal container.
	// These data can be retrieved by
	WithTaggedData(siteScenes TaggedData) Buildable
	// WithCause sets the underlying error manually if necessary.
	WithCause(cause error) Buildable

	// End could terminate the with-build stream calls without any return value.
	End()

	// Container _
	Container
}

type causer interface {
	// Cause returns the underlying cause of the error, if possible.
	// An error value has a cause if it implements the following
	// interface:
	//
	//     type causer interface {
	//            Cause() error
	//     }
	//
	// If an error object does not implement Cause interface, the
	// original error object will be returned.
	// If the error is nil, nil will be returned without further
	// investigation.
	Cause() error
}

// causers is a tool interface. In your scene, use errors.Causes(err)
// to extract the inner errors. Or, use As():
//
//      err := New("many inner errors").WithErrors(e1,e2,e3)
//      var errs []error
//      errors.As(err, &errs)
//      errs = errors.Causes(err)
//
// You may extract the inner errors one by one:
//
//      var e error = err
//      for e != nil {
//          e = errors.Unwrap(err)
//      }
//
type causers interface {
	// Causes _
	Causes() []error
}

// Container represents an error container which can hold a group
// of inner errors.
type Container interface {
	// IsEmpty tests has attached errors
	IsEmpty() bool
	// Defer can be used as a defer function to simplify your codes.
	//
	// The codes:
	//
	//     func some(){
	//       // as a inner errors container
	//       child := func() (err error) {
	//      	errContainer := errors.New("")
	//      	defer errContainer.Defer(&err)
	//
	//      	for _, r := range []error{io.EOF, io.ErrClosedPipe, errors.Internal} {
	//      		errContainer.Attach(r)
	//      	}
	//
	//      	return
	//       }
	//
	//       err := child()
	//       t.Logf("failed: %+v", err)
	//    }
	//
	Defer(err *error)
	// Attachable _
	Attachable
}

// Attachable _
type Attachable interface {
	// Attach collects the errors except it's nil
	Attach(errs ...error)
}

// TaggedData _
type TaggedData map[string]interface{}