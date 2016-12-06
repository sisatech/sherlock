package sherlock

// Check takes an arbitrary number of arguments and checks only the final one.
// If the final argument is of type error and is non nil, it is thrown as a
// sherlock panic.
func Check(args ...interface{}) {

	last := args[len(args)-1]

	if last == nil {
		return
	}

	err, ok := last.(error)
	if ok {
		Throw(err)
	}

}

// CheckThrow looks at the first and last arguments provided to it. If the last
// argument is non-nil and is also type error then it Throws err.
func CheckThrow(err error, args ...interface{}) {

	last := args[len(args)-1]

	if last == nil {
		return
	}

	_, ok := last.(error)
	if ok {
		Throw(err)
	}

}
