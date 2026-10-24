// Package trace tells where in the code a call came from: the function, the file and the line.
//
// Here captures the line that calls it, and Caller the line a number of calls above, for helpers that report
// where they were called from. Both return a Frame: a program counter, which becomes text only when it is used.
//
//	if err != nil {
//		return fmt.Errorf("%s: %w", trace.Here(), err)
//	}
//
// The error reads "service.go:42 users.(*Service).Create: user not found".
//
// # Frames
//
// A Frame's String is "file:line function", with the base name of the file and the function's package path
// cut to its last element; a major version suffix keeps the element before it, as in "trace/v2.Here".
// Function, File and Line return the full function name, the full file path and the line on their own.
// The zero Frame, which Caller returns when the stack is less deep than asked, reads "???".
//
// Frames are comparable, so they can be map keys: frames captured at the same call are equal. A function
// the compiler inlined into several others has a copy of its code in each, and each copy has frames of its own.
// So does a generic function compiled for type arguments of different shapes, such as int and string.
//
// # Inlining
//
// Frames are resolved with runtime.CallersFrames, which accounts for inlining: a frame in a function the compiler
// inlined names that function and its line, not the function it was inlined into, and Caller counts inlined calls
// like any other.
//
// # Logging
//
// A Frame implements slog.LogValuer, so an attribute holding one logs as its text, which the handler asks for
// only if the line is written, and which costs no allocation:
//
//	logger.Warn("slow query", "at", trace.Caller(1))
//
// A Frame holds the same program counter as slog.Record's PC field. A helper that logs for its caller passes
// uintptr(trace.Caller(1)) to slog.NewRecord, so the record's source is the caller's line, and a handler
// formats a record's PC with trace.Frame(r.PC).
//
// It also implements encoding.TextAppender and encoding.TextMarshaler, so encoders such as encoding/json write it
// as its String, not as a number.
//
// # Cost
//
// Here and Caller allocate nothing; the cost is runtime.Callers walking up the stack. The first time a frame
// is used, it is resolved and the result kept, which allocates a few times, once per call site. After that,
// String, AppendText, Function, File, Line and LogValue allocate nothing; MarshalText allocates the slice it returns.
// The zero Frame costs nothing to print. A Frame made from any other integer that isn't a program counter
// is resolved again on every use.
// The results are kept for as long as the program runs, one per call site used, so they are bounded by the size
// of the program.
package trace
