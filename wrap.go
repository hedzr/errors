// Copyright © 2019 Hedzr Yeh.

//+build go1.13

package errors

import "errors"

// Is reports whether any error in err's chain matches target.
//
// The chain consists of err itself followed by the sequence of errors obtained by
// repeatedly calling Unwrap.
//
// An error is considered to match a target if it is equal to that target or if
// it implements a method Is(error) bool such that Is(target) returns true.
func Is(err, target error) bool {
	return errors.Is(err, target)
}

// As finds the first error in err's chain that matches target, and if so, sets
// target to that error value and returns true.
//
// The chain consists of err itself followed by the sequence of errors obtained by
// repeatedly calling Unwrap.
//
// An error matches target if the error's concrete value is assignable to the value
// pointed to by target, or if the error has a method As(interface{}) bool such that
// As(target) returns true. In the latter case, the As method is responsible for
// setting target.
//
// As will panic if target is not a non-nil pointer to either a type that implements
// error, or to any interface type. As returns false if err is nil.
func As(err error, target interface{}) bool {
	return errors.As(err, target)
}

// Unwrap returns the result of calling the Unwrap method on err, if err's
// type contains an Unwrap method returning error.
// Otherwise, Unwrap returns nil.
func Unwrap(err error) error {
	u, ok := err.(interface {
		Unwrap() error
	})
	if !ok {
		return nil
	}
	return u.Unwrap()
	// 中文是不是太不正常？这取决于内存余量够不够多。三个goland加上两个调试一对tcp server+client的话，在goland中输入中文会出现不可忍受的延迟，浮动窗口不弹出或者不消失。哎呀有趣有趣，搜狗不会出现这样de问题；但随即我已彻底卸载了搜狗输入法，这家伙搜集我的各种键入信息，代码、密码、不可描述的搜索关键词之类的，我还是不能交给他啊，起码从现在开始不能再交给它了，搜狗现在是和腾讯打得火热，重点是国内的互联网"大咖"，有一个算一个，没有一个值得我尊敬，他们聚敛财富的方式太脏了。
}
