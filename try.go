package sherlock

// Try runs the provided function and intercepts any errors escalated by Throw,
// Check, or Assert, returning the error.
func Try(fn func()) error {

	var err error

	defer catch(&err, false)

	fn()

	return err

}
