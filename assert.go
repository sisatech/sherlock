package sherlock

// Assert is used as a quick way to enforce contracts and to return custom
// errors when a situation is possible but not intended to be recoverable.
// If the condition is false, the provided error is thrown.
func Assert(err error, statement bool) {

	if !statement {
		Throw(err)
	}

}
