package sherlock

import (
	"fmt"
	"os"
	"runtime/debug"
	"strings"
)

var unhandled interface{}

func catch(err *error, loud bool) {

	r := recover()

	// no panic intercepted
	if r == nil {
		err = nil
		return
	}

	// panic intercepted

	// check if panic was thrown by pi
	report, ok := r.(*report)

	if ok {

		// panic was thrown by pi
		*err = report.err

		if loud {
			output(*err, true)
		}

	} else {

		// panic not thrown by pi

		// put helpful information to the screen only if no other catch
		// has already printed the helpful information.
		if unhandled != r {

			unhandled = r

			output(r, false)

		}

		// continue escalating the panic
		panic(r)

	}

}

func output(x interface{}, expected bool) {

	cut := 11
	str := "intercepted unexpected panic"
	if expected {
		cut += 2
		str = "intercepted thrown error"
	}

	stack := string(debug.Stack())
	lines := strings.Split(string(debug.Stack()), "\n")
	stack = strings.Join(lines[cut:], "\n")

	fmt.Fprintf(os.Stderr, "%v: %v\n\n%v\n%v\n", str, x, lines[0], stack)

}
