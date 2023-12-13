package errors

import (
	"bytes"
	"fmt"
	"io"
	"strconv"
)

// A Code is a signed 32-bit error code copied from gRPC spec
// but negatived.
//
// And more builtin error codes added since hedzr/errors.v3 (v3.1.6).
//
// You may register any application-level error codes by calling
// RegisterCode(codeInt, desc).
type Code int32

const (
	// OK is returned on success. [HTTP/non-HTTP]
	OK Code = 0

	// Canceled indicates the operation was canceled (typically by the caller). [HTTP/non-HTTP]
	Canceled Code = -1

	// Unknown error. [HTTP/non-HTTP]
	// An example of where this error may be returned is
	// if a Status value received from another address space belongs to
	// an error-space that is not known in this address space. Also
	// errors raised by APIs that do not return enough error information
	// may be converted to this error.
	Unknown Code = -2

	// InvalidArgument indicates client specified an invalid argument. [HTTP]
	//
	// Note that this differs from FailedPrecondition. It indicates
	// arguments that are problematic regardless of the state of the
	// system (e.g., a malformed file name).
	//
	// And this also differs from IllegalArgument, who's applied and
	// identify an application error or a general logical error.
	InvalidArgument Code = -3

	// DeadlineExceeded means operation expired before completion. [HTTP]
	// For operations that change the state of the system, this error
	// might be returned even if the operation has completed
	// successfully. For example, a successful response from a server
	// could have been delayed long enough for the deadline to expire.
	//
	// = HTTP 408 Timeout
	DeadlineExceeded Code = -4

	// NotFound means some requested entity (e.g., file or directory)
	// wasn't found. [HTTP]
	//
	// = HTTP 404
	NotFound Code = -5

	// AlreadyExists means an attempt to create an entity failed
	// because one already exists. [HTTP]
	AlreadyExists Code = -6

	// PermissionDenied indicates the caller does not have permission to
	// execute the specified operation. [HTTP]
	//
	// It must not be used for rejections
	// caused by exhausting some resource (use ResourceExhausted
	// instead for those errors). It must not be
	// used if the caller cannot be identified (use Unauthenticated
	// instead for those errors).
	PermissionDenied Code = -7

	// ResourceExhausted indicates some resource has been exhausted,
	// perhaps a per-user quota, or perhaps the entire file system
	// is out of space. [HTTP]
	ResourceExhausted Code = -8

	// FailedPrecondition indicates operation was rejected because the
	// system is not in a state required for the operation's execution.
	// For example, directory to be deleted may be non-empty, a rmdir
	// operation is applied to a non-directory, etc. [HTTP]
	//
	// A litmus test that may help a service implementor in deciding
	// between FailedPrecondition, Aborted, and Unavailable:
	//  (a) Use Unavailable if the client can retry just the failing call.
	//  (b) Use Aborted if the client should retry at a higher-level
	//      (e.g., restarting a read-modify-write sequence).
	//  (c) Use FailedPrecondition if the client should not retry until
	//      the system state has been explicitly fixed. E.g., if a "rmdir"
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
	// etc. [HTTP]
	//
	// See litmus test above for deciding between FailedPrecondition,
	// Aborted, and Unavailable.
	Aborted Code = -10

	// OutOfRange means operation was attempted past the valid range. [HTTP]
	//
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
	// supported/enabled in this service. [HTTP]
	Unimplemented Code = -12

	// Internal errors [HTTP].
	//
	// Means some invariants expected by underlying
	// system has been broken. If you see one of these errors,
	// something is very broken.
	Internal Code = -13

	// Unavailable indicates the service is currently unavailable. [HTTP]
	// This is a most likely a transient condition and may be corrected
	// by retrying with a backoff. Note that it is not always safe to
	// retry non-idempotent operations.
	//
	// See litmus test above for deciding between FailedPrecondition,
	// Aborted, and Unavailable.
	Unavailable Code = -14

	// DataLoss indicates unrecoverable data loss or corruption. [HTTP]
	DataLoss Code = -15

	// Unauthenticated indicates the request does not have valid
	// authentication credentials for the operation. [HTTP]
	//
	// = HTTP 401 Unauthorized
	Unauthenticated Code = -16

	// RateLimited indicates some flow control algorithm is running. [HTTP]
	// and applied.
	//
	// = HTTP Code 429
	RateLimited Code = -17

	// BadRequest generates a 400 error. [HTTP]
	//
	// = HTTP 400
	BadRequest Code = -18

	// Conflict generates a 409 error. [HTTP]
	//
	// = hTTP 409
	Conflict Code = -19

	// Forbidden generates a 403 error. [HTTP]
	Forbidden Code = -20

	// InternalServerError generates a 500 error. [HTTP]
	InternalServerError Code = -21

	// MethodNotAllowed generates a 405 error. [HTTP]
	MethodNotAllowed Code = -22

	// Timeout generates a Timeout error.
	Timeout Code = -23

	// IllegalState is used for the application is entering a
	// bad state.
	IllegalState Code = -24

	// IllegalFormat can be used for Format failed, user input parsing
	// or analysis failed, etc.
	IllegalFormat Code = -25

	// IllegalArgument is like InvalidArgument but applied on application.
	IllegalArgument = -26

	// InitializationFailed is used for application start up unsuccessfully.
	InitializationFailed = -27

	// DataUnavailable is used for the data fetching failed.
	DataUnavailable = -28

	// UnsupportedOperation is like MethodNotAllowed but applied on application.
	UnsupportedOperation = -29

	// UnsupportedVersion can be used for production continuously iteration.
	UnsupportedVersion = -30

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
	"Illegal Format":        IllegalFormat,
	"Illegal State":         IllegalState,
	"Illegal Argument":      IllegalArgument,
	"Initialization Failed": InitializationFailed,
	"Data Unavailable":      DataUnavailable,
	"Unsupported Operation": UnsupportedOperation,
	"Unsupported Version":   UnsupportedVersion,
}

var codeToStr = map[Code]string{
	OK:                   `OK`,
	Canceled:             `CANCELLED`,
	Unknown:              `UNKNOWN`,
	InvalidArgument:      `INVALID_ARGUMENT`,
	DeadlineExceeded:     `DEADLINE_EXCEEDED`,
	NotFound:             `NOT_FOUND`,
	AlreadyExists:        `ALREADY_EXISTS`,
	PermissionDenied:     `PERMISSION_DENIED`,
	ResourceExhausted:    `RESOURCE_EXHAUSTED`,
	FailedPrecondition:   `FAILED_PRECONDITION`,
	Aborted:              `ABORTED`,
	OutOfRange:           `OUT_OF_RANGE`,
	Unimplemented:        `UNIMPLEMENTED`,
	Internal:             `INTERNAL`,
	Unavailable:          `UNAVAILABLE`,
	DataLoss:             `DATA_LOSS`,
	Unauthenticated:      `UNAUTHENTICATED`,
	RateLimited:          `RATE_LIMITED`,
	BadRequest:           `BAD_REQUEST`,
	Conflict:             `CONFLICT`,
	Forbidden:            `FORBIDDEN`,
	InternalServerError:  `INTERNAL_SERVER_ERROR`,
	MethodNotAllowed:     `METHOD_NOT_ALLOWED`,
	Timeout:              `TIMEOUT`,
	IllegalState:         "Illegal State",
	IllegalFormat:        "Illegal Format",
	IllegalArgument:      "Illegal Argument",
	InitializationFailed: "Initialization Failed",
	DataUnavailable:      "Data Unavailable",
	UnsupportedOperation: "Unsupported Operation",
	UnsupportedVersion:   "Unsupported Version",
}

//
// ----------------------------
//

// New create a new *CodedErr object based an error code
func (c Code) New(msg string, args ...interface{}) Buildable { //nolint:revive
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

func (c Code) makeErrorString(line bool) string { //nolint:revive,unparam
	var buf bytes.Buffer
	_, _ = buf.WriteString(c.Error())
	_, _ = buf.WriteRune(' ')
	_, _ = buf.WriteRune('(')
	_, _ = buf.WriteString(strconv.Itoa(int(c)))
	_, _ = buf.WriteRune(')')
	return buf.String()
}

func (c Code) Is(other error) bool {
	if o, ok := other.(Code); ok && o == c {
		return true
	}
	return false
}

// Format formats the stack of Frames according to the fmt.Formatter interface.
//
//	%s	lists source files for each Frame in the stack
//	%v	lists the source file and line number for each Frame in the stack
//
// Format accepts flags that alter the printing of some verbs, as follows:
//
//	%+v   Prints filename, function, and line number for each Frame in the stack.
func (c Code) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		if s.Flag('+') {
			_, _ = fmt.Fprintf(s, "%+v", c.makeErrorString(true))
			// c.Stack.Format(s, verb)
			return
		}
		fallthrough
	case 's':
		_, _ = io.WriteString(s, c.makeErrorString(false))
	case 'q':
		_, _ = fmt.Fprintf(s, "%q", c.makeErrorString(false))
	}
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
// When a negative integer given, and it's less than
// errors.MinErrorCode, it'll be used as it.
// Or an errors.AlreadyExists returned.
//
// An existing code will be returned directly.
//
// RegisterCode provides a shortcut to declare a number as Code
// as your need.
//
// The best way is:
//
//	var ErrAck = errors.RegisterCode(3, "cannot ack")     // ErrAck will be -1003
//	var ErrAck = errors.RegisterCode(-1003, "cannot ack)  // equivalent with last line
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
