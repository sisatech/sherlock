package sherlock

// Yell has identical functionality to Try, except that it prints information
// about any errors it intercepts to stderr to make debugging easy.
func Yell(fn func()) error {

	var err error

	defer catch(&err, true)

	fn()

	return err

}
