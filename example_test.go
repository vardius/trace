package trace_test

import (
	"errors"
	"fmt"

	"github.com/vardius/trace"
)

func Example_here() {
	c := func() {
		fmt.Printf("%s\n%s", "Hello from:", trace.Here(trace.Lline|trace.Lfunction))
	}

	b := func() { c() }
	a := func() { b() }

	a()

	// Output:
	// Hello from:
	// :12 github.com/vardius/trace_test.Example_here.func1
}

func Example_here_second() {
	err := errors.New("Internal system error")

	fmt.Printf("%s %s\n%s", "Error:", err, trace.Here(trace.Lline|trace.Lfunction))

	// Output:
	// Error: Internal system error
	// :28 github.com/vardius/trace_test.Example_here_second
}

func Example_fromParent() {
	c := func() {
		fmt.Printf("%s\n%s", "Hello from:", trace.FromParent(3, trace.Lline|trace.Lfunction))
	}

	b := func() { c() }
	a := func() { b() }

	a()

	// Output:
	// Hello from:
	// :43 github.com/vardius/trace_test.Example_fromParent
}
