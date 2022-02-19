package errors

//import (
//	"bytes"
//	"fmt"
//	"gopkg.in/hedzr/errors.v2/old"
//	"io"
//)

//// NewTemplate create an error template so that you may `FormatNew(liveArgs...)` late.
//func (c Code) NewTemplate(tmpl string) *WithCodeInfo {
//	err := &WithCodeInfo{
//		code:   c,
//		causer: nil,
//		msg:    tmpl,
//	}
//	return err
//}

//// WithCodeInfo is a type integrating both error code, cause, message, and template
//type WithCodeInfo struct {
//	code      Code
//	causer    error
//	msg       string
//	livedArgs []interface{}
//}
//
//// Code returns the error code value
//func (w *WithCodeInfo) Code() Code {
//	return w.code
//}
//
//// Equal tests if equals with code 'c'
//func (w *WithCodeInfo) Equal(c Code) bool {
//	return w.code == c
//}
//
//func (w *WithCodeInfo) Error() string {
//	var buf bytes.Buffer
//	buf.WriteString(w.code.String())
//	if len(w.msg) > 0 {
//		buf.WriteRune('|')
//		if len(w.livedArgs) > 0 {
//			buf.WriteString(fmt.Sprintf(w.msg, w.livedArgs))
//		} else {
//			buf.WriteString(w.msg)
//		}
//	}
//	if w.causer != nil {
//		buf.WriteRune('|')
//		buf.WriteString(w.causer.Error())
//	}
//	return buf.String()
//}
//
//// Format formats the stack of Frames according to the fmt.Formatter interface.
////
////    %s	lists source files for each Frame in the stack
////    %v	lists the source file and line number for each Frame in the stack
////
//// Format accepts flags that alter the printing of some verbs, as follows:
////
////    %+v   Prints filename, function, and line number for each Frame in the stack.
//func (w *WithCodeInfo) Format(s fmt.State, verb rune) {
//	switch verb {
//	case 'v':
//		if s.Flag('+') {
//			msg := w.msg
//			if len(w.livedArgs) > 0 {
//				msg = fmt.Sprintf(w.msg, w.livedArgs...)
//			}
//			_, _ = fmt.Fprintf(s, "%d|%+v|%s", int(w.code), w.code.String(), msg)
//			if w.causer != nil {
//				_, _ = fmt.Fprintf(s, "|%+v", w.causer)
//			}
//			return
//		}
//		fallthrough
//	case 's':
//		_, _ = io.WriteString(s, w.Error())
//	case 'q':
//		_, _ = fmt.Fprintf(s, "%q", w.Error())
//	}
//}
//
//// FormatNew creates a new error object based on this error template 'w'.
////
//// Example:
////
//// 	   errTmpl1001 := BUG1001.NewTemplate("something is wrong %v")
//// 	   err4 := errTmpl1001.FormatNew("ok").Attach(errBug1)
//// 	   fmt.Println(err4)
//// 	   fmt.Printf("%+v\n", err4)
////
//func (w *WithCodeInfo) FormatNew(livedArgs ...interface{}) *old.WithStackInfo {
//	x := WithCode(w.code, w.causer, w.msg)
//	x.error.(*WithCodeInfo).livedArgs = livedArgs
//	return x
//}
//
//// Attach appends errs
//func (w *WithCodeInfo) Attach(errs ...error) {
//	for _, err := range errs {
//		if err != nil {
//			w.causer = err
//		}
//	}
//	if len(errs) > 1 {
//		panic("*WithCodeInfo.Attach() can only wrap one child error object.")
//	}
//}
//
//// Cause returns the underlying cause of the error recursively,
//// if possible.
//func (w *WithCodeInfo) Cause() error {
//	return w.causer
//}
//
//// Unwrap returns the result of calling the Unwrap method on err, if err's
//// type contains an Unwrap method returning error.
//// Otherwise, Unwrap returns nil.
//func (w *WithCodeInfo) Unwrap() error {
//	return w.causer
//}
//
//// As finds the first error in err's chain that matches target, and if so, sets
//// target to that error value and returns true.
//func (w *WithCodeInfo) As(target interface{}) bool {
//	return old.As(w.causer, target)
//}
//
//// Is reports whether any error in err's chain matches target.
//func (w *WithCodeInfo) Is(target error) bool {
//	return w.causer == target || old.Is(w.causer, target)
//}
//
//// TypeIs reports whether any error in err's chain matches target.
//func (w *WithCodeInfo) TypeIs(target error) bool {
//	return w.causer == target || old.TypeIs(w.causer, target)
//}
