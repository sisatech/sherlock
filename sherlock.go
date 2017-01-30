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

Written by Sisa-Tech Pty Ltd
*/
package sherlock
