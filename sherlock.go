/*
Package sherlock helps tidy up go code by reducing the substantial number of
"if err != nil" checks that are usually performed on code even in cases where
errors are unexpected. Instead of propagating error return values all the way up
the stack, sherlock can unwind the stack by using panic. An appropriately placed
Catch can be used to manage thrown errors higher up in the stack.

A typical use for sherlock will have exported package functions return errors as
is convention, but they will also contain a CatchAll, allowing unexported
functions to do away with this practice.

	func MyFunction() error {
		var err error
		func() {
			defer sherlock.CatchAll(&err)
			// function logic
		}()
		return err
	}

Written by Alan Murtagh
*/
package sherlock

import (
	"fmt"
	"os"
	"runtime"
	"runtime/debug"
	"strings"
)

type report struct {
	err   error
	stack string
	pkg   string
}

// Assert is used as a quick way to enforce contracts and to return custom
// errors when a situation is possible but not intended to be recoverable.
// If the condition is false, the provided error is thrown.
func Assert(condition bool, err error) {
	if condition {
		return
	}
	panic(&report{
		err:   err,
		stack: stacktrace(),
		pkg:   caller(),
	})
}

// Catch halts a sherlock panic and checks if the thrown error is the same error
// provided as an argument. If the errors match then the provided function is
// executed and the panic is recovered. If the errors do not match then the
// error is rethrown. Errors are equal only if they have the same address.
func Catch(err error, fn func()) {
	r := recover()
	if r == nil {
		return
	}
	x, ok := r.(*report)
	if !ok || x.pkg != caller() {
		fmt.Fprintf(os.Stderr, "%v", string(debug.Stack()))
		panic(r)
	}
	if err == x.err {
		fn()
	} else {
		panic(r)
	}
}

// CatchAll halts a sherlock panic and fills the provided error pointer with the
// error that was thrown.
//
// Sherlock can only catch sherlock thrown panics, and will rethrow a
// non-sherlock panic. This behaviour can result in final stack traces being
// difficult to use, but it is assumed that any non-sherlock panic is a bug, and
// so sherlock will dump a stacktrace into stderr.
//
// Sherlock can also only catch panics thrown within the same package. It is
// good practice to not let panics unwind beyond the boundaries of a package,
// and so this is considered a bug by sherlock.
func CatchAll(err *error) {
	r := recover()
	if r == nil {
		err = nil
		return
	}

	x, ok := r.(*report)
	if !ok {
		x, ok := r.(error)
		if ok {
			fmt.Fprintf(os.Stderr, "\n%v\n\n", x.Error())
		}
		fmt.Fprintf(os.Stderr, "%v\n", string(debug.Stack()))
		panic(r)
	} else if x.pkg != caller() {
		fmt.Fprintf(os.Stderr, "%v\n", x.err.Error())
	}
	*err = x.err
}

// Check takes an arbitrary number of arguments and checks only the final one.
// If the final argument is of type error and is non nil, it is thrown as a
// sherlock panic.
func Check(args ...interface{}) {
	l := len(args)
	if args[l-1] == nil {
		return
	}
	err, ok := args[l-1].(error)
	if !ok {
		return
	}
	panic(&report{
		err:   err,
		stack: stacktrace(),
		pkg:   caller(),
	})
}

// Throw simply throws the provided error as a sherlock panic.
func Throw(err error) {
	panic(&report{
		err:   err,
		stack: stacktrace(),
		pkg:   caller(),
	})
}

func stacktrace() string {
	// TODO: remove parts of stacktrace that exist due to this package.
	return string(debug.Stack())
}

// NOTE: caller determines the calling package by skipping up the stack and
// determining which package the calling function's calling function came from.
// Take care to ensure it is never used any further down the stack.
func caller() string {
	_, file, _, ok := runtime.Caller(1)
	if !ok {
		panic(nil)
	}
	i := strings.LastIndex(file, "/")
	return file[:i]
}
