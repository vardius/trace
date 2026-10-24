package trace

import (
	"log/slog"
	"math"
	"runtime"
	"strconv"
	"strings"
	"sync"
)

// Frame is a place in the code: the program counter of a call, in the form runtime.Callers reports it.
// Here and Caller capture one, and its methods say which function, file and line it is.
//
// A Frame is an integer, so it is cheap to keep, and frames can be compared and used as map keys:
// frames captured by the same call are equal. A slog.Record's PC converts to a Frame, and a Frame converts back
// with uintptr, for slog.NewRecord. The zero Frame is an unknown place.
type Frame uintptr

// Here returns the frame of the line that calls it.
func Here() Frame {
	var pc [1]uintptr
	runtime.Callers(2, pc[:]) // skips runtime.Callers and Here
	return Frame(pc[0])
}

// Caller returns the frame skip calls above the line that calls it: Caller(0) is the same as Here,
// and Caller(1) is the line that called the function that calls Caller.
// It returns the zero Frame when the stack is less deep than that. A negative skip counts as 0.
func Caller(skip int) Frame {
	var pc [1]uintptr
	// Skips runtime.Callers and Caller. The bound keeps skip+2 from overflowing into a negative skip,
	// which runtime.Callers would read as 0 and report itself.
	runtime.Callers(min(max(skip, 0), math.MaxInt-2)+2, pc[:])
	return Frame(pc[0])
}

// Function returns the package path-qualified name of the function the frame is in,
// such as github.com/vardius/trace/v2.Here, or "" if the frame is unknown.
func (f Frame) Function() string { return f.location().function }

// File returns the full path of the frame's source file, with forward slashes, or "" if the frame is unknown.
func (f Frame) File() string { return f.location().file }

// Line returns the frame's line number, or 0 if the frame is unknown.
func (f Frame) Line() int { return f.location().line }

// String returns the frame as "file:line function", with the file's base name and the function's package path
// cut to its last element, such as "service.go:42 users.(*Service).Create". A major version suffix keeps
// the element before it, as in "trace/v2.Here". An unknown frame is "???".
func (f Frame) String() string { return f.location().text }

// AppendText appends the frame's String to b. It implements encoding.TextAppender and never returns an error.
func (f Frame) AppendText(b []byte) ([]byte, error) { return append(b, f.location().text...), nil }

// MarshalText returns the frame's String. It implements encoding.TextMarshaler, so encoding/json and other
// encoders write a Frame as text rather than as a number.
func (f Frame) MarshalText() ([]byte, error) { return f.AppendText(nil) }

// LogValue returns the frame's String. It implements slog.LogValuer, so a Frame logged as an attribute
// is written as text, and LogValue allocates nothing.
func (f Frame) LogValue() slog.Value { return slog.StringValue(f.location().text) }

type location struct {
	function, file string
	line           int
	text           string // what String returns
}

// unknown is the location of the zero Frame and of any program counter outside Go code.
var unknown = &location{text: "???"}

// locations holds the location of every frame resolved so far, one per call site, so it is bounded by the size
// of the program. A program counter always means the same place, so the cache never goes stale.
// Resolving a frame and building its text allocates several times and takes about 200 ns (BenchmarkFirstUse);
// a cached lookup allocates nothing and takes a few.
var locations sync.Map // Frame → *location

func (f Frame) location() *location {
	if f == 0 {
		return unknown
	}
	if l, ok := locations.Load(f); ok {
		return l.(*location)
	}

	// CallersFrames, not FuncForPC: it is the documented way to resolve a program counter from runtime.Callers,
	// and it accounts for inlining.
	fr, _ := runtime.CallersFrames([]uintptr{uintptr(f)}).Next()
	if fr.File == "" {
		return unknown // not cached, so arbitrary integers can't grow the cache
	}
	file := fr.File[strings.LastIndexByte(fr.File, '/')+1:]
	l := &location{
		function: fr.Function,
		file:     fr.File,
		line:     fr.Line,
		text:     file + ":" + strconv.Itoa(fr.Line) + " " + shortName(fr.Function),
	}
	actual, _ := locations.LoadOrStore(f, l)
	return actual.(*location)
}

// shortName cuts a function name's package path to its last element: github.com/acme/users.(*Service).Create
// becomes users.(*Service).Create. When the last element is a version, such as v2 in github.com/vardius/trace/v2.Here,
// the element before it stays too: trace/v2.Here.
func shortName(function string) string {
	i := strings.LastIndexByte(function, '/')
	if i >= 0 && isVersion(function[i+1:]) {
		i = strings.LastIndexByte(function[:i], '/')
	}
	return function[i+1:]
}

// isVersion reports whether name, a function name less its package path, starts with a version element:
// v and digits, followed by the dot before the function or the _test of an external test package.
func isVersion(name string) bool {
	number, ok := strings.CutPrefix(name, "v")
	if !ok {
		return false
	}
	after := strings.TrimLeft(number, "0123456789")
	return len(after) < len(number) && (strings.HasPrefix(after, ".") || strings.HasPrefix(after, "_"))
}
