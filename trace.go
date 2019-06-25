package trace

import (
	"bytes"
	"fmt"
	"runtime"
)

const (
	Lfile     = 1 << iota         // full file name
	Lline                         // line number
	Lfunction                     // name of the function
	LstdFlags = Lfile | Lfunction // initial values
)

// Here returns string representation of a reference
// default flags LstdFlags
func Here(flag int) string {
	if flag == 0 {
		flag = LstdFlags
	}

	var buf bytes.Buffer

	frame := getFrame(2)
	if frame != nil {
		if flag&Lfile != 0 {
			fmt.Fprintf(&buf, "%s", frame.File)
		}
		if flag&Lline != 0 {
			fmt.Fprintf(&buf, ":%d", frame.Line)
		}
		if flag&Lfunction != 0 {
			fmt.Fprintf(&buf, " %s", frame.Function)
		}
	}

	return buf.String()
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
