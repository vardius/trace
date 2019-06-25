package trace_test

import (
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
	// /home/travis/gopath/src/github.com/vardius/trace/example_test.go:11 github.com/vardius/trace_test.Example_here.func1
}
