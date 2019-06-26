package trace

import (
	"bytes"
	"fmt"
	"io"
	"runtime"
)

// Output flags
const (
	Lfile     = 1 << iota         // full file name
	Lline                         // line number
	Lfunction                     // name of the function
	LstdFlags = Lfile | Lfunction // initial values
)

// Here returns string representation of a reference
// default flags LstdFlags
func Here(flags int) string {
	return FromParent(1, flags)
}

// FromParent returns string representation of a parent reference
// default flags LstdFlags
func FromParent(calldepth int, flags int) string {
	var buf bytes.Buffer

	frame := getFrame(calldepth + 2)
	outputFrame(flags, frame, &buf)

	return buf.String()
}

func outputFrame(flags int, frame *runtime.Frame, w io.Writer) {
	if flags == 0 {
		flags = LstdFlags
	}

	if frame != nil {
		if flags&Lfile != 0 {
			fmt.Fprintf(w, "%s", frame.File)
		}
		if flags&Lline != 0 {
			fmt.Fprintf(w, ":%d", frame.Line)
		}
		if flags&Lfunction != 0 {
			fmt.Fprintf(w, " %s", frame.Function)
		}
	}
}

func getFrame(calldepth int) *runtime.Frame {
	pc, file, line, ok := runtime.Caller(calldepth)
	if !ok {
		return nil
	}

	frame := &runtime.Frame{
		PC:   pc,
		File: file,
		Line: line,
	}

	funcForPc := runtime.FuncForPC(pc)
	if funcForPc != nil {
		frame.Func = funcForPc
		frame.Function = funcForPc.Name()
		frame.Entry = funcForPc.Entry()
	}

	return frame
}
