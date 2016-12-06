package sherlock

// Throw sends an error to be intercepted by Try. If it is not called from
// within a Try it will cause a panic.
func Throw(err error) {

	panic(&report{
		err: err,
	})

}
