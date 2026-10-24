// Command once warns that a function is deprecated the first time each line calls it, and never again for that line.
//
//	go run ./examples/once
package main

import (
	"fmt"
	"io"
	"os"
	"sync"

	"github.com/vardius/trace/v2"
)

// warned holds the frames of the calls already warned about. A Frame is an integer, so it is a cheap map key.
var warned sync.Map

// warnOnce writes msg the first time a line calls the function that calls warnOnce.
func warnOnce(w io.Writer, msg string) {
	at := trace.Caller(2) // 0 is this line, 1 the call to warnOnce, 2 the call to that function
	if _, seen := warned.LoadOrStore(at, struct{}{}); !seen {
		fmt.Fprintf(w, "%s: %s\n", at, msg)
	}
}

// Sum adds a and b.
//
// Deprecated: use a + b.
func Sum(a, b int) int {
	warnOnce(os.Stderr, "Sum is deprecated: use a + b")
	return a + b
}

func main() {
	total := 0
	for i := range 3 {
		total = Sum(total, i) // warns the first time only
	}
	fmt.Println(Sum(total, 10)) // another line, so another warning
}
