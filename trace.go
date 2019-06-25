package trace

import (
	"bytes"
	"fmt"
	"runtime"
)

// Here returns string representation of a reference
// to a File:Line Function calling this function
func Here() string {
	var buf bytes.Buffer

	// Ask for up to 3 pcs, including:
	// - runtime.Callers,
	// - trace.Here
	// - pcs for where golog.getPCs was called
	pc := make([]uintptr, 3)
	n := runtime.Callers(0, pc)

	if n < 3 {
		pc = pc[:]
	}

	// pass only valid pcs to runtime.CallersFrames
	// exclude irrelevant pcs
	pc = pc[2:n]

	if len(pc) > 0 {
		frames := runtime.CallersFrames(pc)
		// Loop to get frames.
		// A fixed number of pcs can expand to an indefinite number of Frames.
		for {
			frame, more := frames.Next()
			fmt.Fprintf(&buf, "%s:%d %s", frame.File, frame.Line, frame.Function)
			if !more {
				break
			}
		}
	}

	return buf.String()
}
