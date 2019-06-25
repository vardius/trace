package trace_test

import (
	"errors"
	"fmt"

	"github.com/vardius/trace"
)

func Example_here() {
	c := func() {
		fmt.Printf("%s\n%s", "Hello from:", trace.Here())
	}

	b := func() { c() }
	a := func() { b() }

	a()

	// Output:
	// Hello from:
	// /home/travis/gopath/src/github.com/vardius/trace/example_test.go:12 github.com/vardius/trace_test.Example_here.func1
}

func Example_here_second() {
	err := errors.New("Internal system error")

	fmt.Printf("%s %s\n%s", "Error:", err, trace.Here())

	// Output:
	// Error: Internal system error
	// /home/travis/gopath/src/github.com/vardius/trace/example_test.go:28 github.com/vardius/trace_test.Example_here_second
}
