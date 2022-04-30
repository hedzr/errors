package errors

// A Code is an signed 32-bit error code copied from gRPC spec but negatived.
type Code int32

const (
	// OK is returned on success.
	OK Code = 0

	// Canceled indicates the operation was canceled (typically by the caller).
	Canceled Code = -1

	// Unknown error. An example of where this error may be returned is
	// if a Status value received from another address space belongs to
	// an error-space that is not known in this address space. Also
	// errors raised by APIs that do not return enough error information
	// may be converted to this error.
	Unknown Code = -2

	// InvalidArgument indicates client specified an invalid argument.
	// Note that this differs from FailedPrecondition. It indicates arguments
	// that are problematic regardless of the state of the system
	// (e.g., a malformed file name).
	InvalidArgument Code = -3

	// DeadlineExceeded means operation expired before completion.
	// For operations that change the state of the system, this error may be
	// returned even if the operation has completed successfully. For
	// example, a successful response from a server could have been delayed
	// long enough for the deadline to expire.
	//
	// = HTTP 408 Timeout
	DeadlineExceeded Code = -4

	// NotFound means some requested entity (e.g., file or directory) was
	// not found.
	//
	// = HTTP 404
	NotFound Code = -5

	// AlreadyExists means an attempt to create an entity failed because one
	// already exists.
	AlreadyExists Code = -6

	// PermissionDenied indicates the caller does not have permission to
	// execute the specified operation. It must not be used for rejections
	// caused by exhausting some resource (use ResourceExhausted
	// instead for those errors). It must not be
	// used if the caller cannot be identified (use Unauthenticated
	// instead for those errors).
	PermissionDenied Code = -7

	// ResourceExhausted indicates some resource has been exhausted, perhaps
	// a per-user quota, or perhaps the entire file system is out of space.
	ResourceExhausted Code = -8

	// FailedPrecondition indicates operation was rejected because the
	// system is not in a state required for the operation's execution.
	// For example, directory to be deleted may be non-empty, an rmdir
	// operation is applied to a non-directory, etc.
	//
	// A litmus test that may help a service implementor in deciding
	// between FailedPrecondition, Aborted, and Unavailable:
	//  (a) Use Unavailable if the client can retry just the failing call.
	//  (b) Use Aborted if the client should retry at a higher-level
	//      (e.g., restarting a read-modify-write sequence).
	//  (c) Use FailedPrecondition if the client should not retry until
	//      the system state has been explicitly fixed. E.g., if an "rmdir"
	//      fails because the directory is non-empty, FailedPrecondition
	//      should be returned since the client should not retry unless
	//      they have first fixed up the directory by deleting files from it.
	//  (d) Use FailedPrecondition if the client performs conditional
	//      REST Get/Update/Delete on a resource and the resource on the
	//      server does not match the condition. E.g., conflicting
	//      read-modify-write on the same resource.
	FailedPrecondition Code = -9

	// Aborted indicates the operation was aborted, typically due to a
	// concurrency issue like sequencer check failures, transaction aborts,
	// etc.
	//
	// See litmus test above for deciding between FailedPrecondition,
	// Aborted, and Unavailable.
	Aborted Code = -10

	// OutOfRange means operation was attempted past the valid range.
	// E.g., seeking or reading past end of file.
	//
	// Unlike InvalidArgument, this error indicates a problem that may
	// be fixed if the system state changes. For example, a 32-bit file
	// system will generate InvalidArgument if asked to read at an
	// offset that is not in the range [0,2^32-1], but it will generate
	// OutOfRange if asked to read from an offset past the current
	// file size.
	//
	// There is a fair bit of overlap between FailedPrecondition and
	// OutOfRange. We recommend using OutOfRange (the more specific
	// error) when it applies so that callers who are iterating through
	// a space can easily look for an OutOfRange error to detect when
	// they are done.
	OutOfRange Code = -11

	// Unimplemented indicates operation is not implemented or not
	// supported/enabled in this service.
	Unimplemented Code = -12

	// Internal errors. Means some invariants expected by underlying
	// system has been broken. If you see one of these errors,
	// something is very broken.
	Internal Code = -13

	// Unavailable indicates the service is currently unavailable.
	// This is a most likely a transient condition and may be corrected
	// by retrying with a backoff. Note that it is not always safe to retry
	// non-idempotent operations.
	//
	// See litmus test above for deciding between FailedPrecondition,
	// Aborted, and Unavailable.
	Unavailable Code = -14

	// DataLoss indicates unrecoverable data loss or corruption.
	DataLoss Code = -15

	// Unauthenticated indicates the request does not have valid
	// authentication credentials for the operation.
	//
	// = HTTP 401 Unauthorized
	Unauthenticated Code = -16

	// RateLimited indicates some flow control algorithm is running and applied.
	// = HTTP Code 429
	RateLimited Code = -17

	// BadRequest generates a 400 error.
	// = HTTP 400
	BadRequest Code = -18

	// Conflict generates a 409 error.
	// = hTTP 409
	Conflict Code = -19

	// Forbidden generates a 403 error.
	Forbidden Code = -20

	// InternalServerError generates a 500 error.
	InternalServerError Code = -21

	// MethodNotAllowed generates a 405 error.
	MethodNotAllowed Code = -22

	// Timeout generates a Timeout error.
	Timeout Code = -23

	// MinErrorCode is the lower bound for user-defined Code.
	MinErrorCode Code = -1000
)

var strToCode = map[string]Code{
	`OK`:                    OK,
	`CANCELLED`:             Canceled,
	`UNKNOWN`:               Unknown,
	`INVALID_ARGUMENT`:      InvalidArgument,
	`DEADLINE_EXCEEDED`:     DeadlineExceeded,
	`NOT_FOUND`:             NotFound,
	`ALREADY_EXISTS`:        AlreadyExists,
	`PERMISSION_DENIED`:     PermissionDenied,
	`RESOURCE_EXHAUSTED`:    ResourceExhausted,
	`FAILED_PRECONDITION`:   FailedPrecondition,
	`ABORTED`:               Aborted,
	`OUT_OF_RANGE`:          OutOfRange,
	`UNIMPLEMENTED`:         Unimplemented,
	`INTERNAL`:              Internal,
	`UNAVAILABLE`:           Unavailable,
	`DATA_LOSS`:             DataLoss,
	`UNAUTHENTICATED`:       Unauthenticated,
	`RATE_LIMITED`:          RateLimited,
	`BAD_REQUEST`:           BadRequest,
	`CONFLICT`:              Conflict,
	`FORBIDDEN`:             Forbidden,
	`INTERNAL_SERVER_ERROR`: InternalServerError,
	`METHOD_NOT_ALLOWED`:    MethodNotAllowed,
	`TIMEOUT`:               Timeout,
}

var codeToStr = map[Code]string{
	OK:                  `OK`,
	Canceled:            `CANCELLED`,
	Unknown:             `UNKNOWN`,
	InvalidArgument:     `INVALID_ARGUMENT`,
	DeadlineExceeded:    `DEADLINE_EXCEEDED`,
	NotFound:            `NOT_FOUND`,
	AlreadyExists:       `ALREADY_EXISTS`,
	PermissionDenied:    `PERMISSION_DENIED`,
	ResourceExhausted:   `RESOURCE_EXHAUSTED`,
	FailedPrecondition:  `FAILED_PRECONDITION`,
	Aborted:             `ABORTED`,
	OutOfRange:          `OUT_OF_RANGE`,
	Unimplemented:       `UNIMPLEMENTED`,
	Internal:            `INTERNAL`,
	Unavailable:         `UNAVAILABLE`,
	DataLoss:            `DATA_LOSS`,
	Unauthenticated:     `UNAUTHENTICATED`,
	RateLimited:         `RATE_LIMITED`,
	BadRequest:          `BAD_REQUEST`,
	Conflict:            `CONFLICT`,
	Forbidden:           `FORBIDDEN`,
	InternalServerError: `INTERNAL_SERVER_ERROR`,
	MethodNotAllowed:    `METHOD_NOT_ALLOWED`,
	Timeout:             `TIMEOUT`,
}

//
// ----------------------------
//

// New create a new *CodedErr object based an error code
func (c Code) New(msg string, args ...interface{}) Buildable {
	return Message(msg, args...).WithCode(c).Build()
}

// WithCode for error interface
func (c *Code) WithCode(code Code) *Code {
	*c = code
	return c
}

// Error for error interface
func (c Code) Error() string { return c.String() }

// String for stringer interface
func (c Code) String() string {
	if x, ok := codeToStr[c]; ok {
		return x
	}
	return codeToStr[Unknown]
}

// Register registers a code and its token string for using later
func (c Code) Register(codeName string) (errno Code) {
	errno = AlreadyExists
	if c <= MinErrorCode || c > 0 {
		if _, ok := strToCode[codeName]; !ok {
			if _, ok = codeToStr[c]; !ok {
				strToCode[codeName] = c
				codeToStr[c] = codeName
				errno = OK
			}
		}
	}
	return
}

// RegisterCode makes a code integer associated with a name, and
// returns the available Code number.
//
// When a positive integer given as codePositive (such as 3),
// it'll be negatived and added errors.MinErrorCode (-1000) so
// the final Code number will be -1003.
//
// When a negative integer given and it's less than
// errors.MinErrorCode, it'll be used as it.
// Or an errors.AlreadyExists returned.
//
// A existing code will be returned directly.
//
// RegisterCode provides a shortcut to declare a number as Code
// as your need.
//
// The best way is:
//
//     var ErrAck = errors.RegisterCode(3, "cannot ack")     // ErrAck will be -1003
//     var ErrAck = errors.RegisterCode(-1003, "cannot ack)  // equivelant with last line
//
func RegisterCode(codePositive int, codeName string) (errno Code) {
	errno = AlreadyExists
	applier := func(c Code) {
		if v, ok := strToCode[codeName]; !ok {
			if _, ok = codeToStr[c]; !ok {
				strToCode[codeName] = c
				codeToStr[c] = codeName
				errno = c
			}
		} else {
			errno = v
		}
	}
	if codePositive > 0 {
		c := MinErrorCode - Code(codePositive)
		applier(c)
	} else if c := Code(codePositive); c < MinErrorCode {
		applier(c)
	}
	return
}
