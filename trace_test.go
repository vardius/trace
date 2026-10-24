package trace_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"math"
	"runtime"
	"slices"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/vardius/trace/v2"
)

const pkg = "github.com/vardius/trace/v2_test"

func TestFrame(t *testing.T) {
	tests := []struct {
		name     string
		frame    func() trace.Frame
		function string
		line     int
		text     string
	}{
		{"function", here, pkg + ".here", 11, "frames_test.go:11 trace/v2_test.here"},
		{"closure", inClosure, pkg + ".inClosure.func1", 14, "frames_test.go:14 trace/v2_test.inClosure.func1"},
		{"method", new(service).method, pkg + ".(*service).method", 19, "frames_test.go:19 trace/v2_test.(*service).method"},
		{"generic function", generic[int], pkg + ".generic[...]", 21, "frames_test.go:21 trace/v2_test.generic[...]"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := tt.frame()
			if got := f.Function(); got != tt.function {
				t.Errorf("Function() = %q, want %q", got, tt.function)
			}
			if got := f.File(); !strings.HasSuffix(got, "/frames_test.go") {
				t.Errorf("File() = %q, want the full path of frames_test.go", got)
			}
			if got := f.Line(); got != tt.line {
				t.Errorf("Line() = %d, want %d", got, tt.line)
			}
			if got := f.String(); got != tt.text {
				t.Errorf("String() = %q, want %q", got, tt.text)
			}
		})
	}
}

// helper in frames_test.go is small enough for the compiler to inline, and helperNoInline is kept from it.
// Either way, the frames must name the functions the source says.
func TestFrameInlined(t *testing.T) {
	tests := []struct {
		name         string
		call         func() (self, caller trace.Frame)
		wantInlined  bool
		self, caller string
	}{
		{"inlined", helped, true, "frames_test.go:24 trace/v2_test.helper", "frames_test.go:26 trace/v2_test.helped"},
		{"not inlined", helpedNoInline, false,
			"frames_test.go:29 trace/v2_test.helperNoInline", "frames_test.go:31 trace/v2_test.helpedNoInline"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.wantInlined && raceEnabled {
				// It inlines Caller into helper instead, which leaves helper too big to inline.
				t.Skip("the race detector changes what the compiler inlines")
			}
			self, caller := tt.call()
			if got := inlined(self); got != tt.wantInlined {
				t.Fatalf("helper inlined = %v, want %v: this build doesn't test what the case says", got, tt.wantInlined)
			}
			if got := self.String(); got != tt.self {
				t.Errorf("Caller(0) in the helper = %q, want %q", got, tt.self)
			}
			if got := caller.String(); got != tt.caller {
				t.Errorf("Caller(1) in the helper = %q, want %q", got, tt.caller)
			}
		})
	}
}

// inlined reports whether the compiler inlined the function f is in, from the runtime's own view of it.
func inlined(f trace.Frame) bool {
	fr, _ := runtime.CallersFrames([]uintptr{uintptr(f)}).Next()
	return fr.Func == nil
}

func TestCaller(t *testing.T) {
	tests := []struct {
		name string
		skip int
		want string
	}{
		{"0 is the line that calls Caller", 0, "frames_test.go:37 trace/v2_test.skipsInner"},
		{"1 is the line above", 1, "frames_test.go:34 trace/v2_test.skips"},
		{"2 is the line above that", 2, "frames_test.go:39 trace/v2_test.skipsFromHere"},
		// v1 reported its own internals for a negative depth.
		{"negative counts as 0", -1, "frames_test.go:37 trace/v2_test.skipsInner"},
		{"beyond the stack", 1 << 20, "???"},
		// v1 overflowed and reported runtime.Caller.
		{"the largest int", math.MaxInt, "???"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := skipsFromHere(tt.skip).String(); got != tt.want {
				t.Errorf("Caller(%d) = %q, want %q", tt.skip, got, tt.want)
			}
		})
	}
}

func TestCallerBeyondTheStackIsZero(t *testing.T) {
	if got := trace.Caller(1 << 20); got != 0 {
		t.Errorf("Caller(1 << 20) = %#x, want the zero Frame", uintptr(got))
	}
}

func TestUnknownFrame(t *testing.T) {
	// 1 is no Go code's program counter.
	for _, f := range []trace.Frame{0, 1} {
		t.Run(fmt.Sprint(uintptr(f)), func(t *testing.T) {
			if f.Function() != "" || f.File() != "" || f.Line() != 0 {
				t.Errorf("Function, File, Line = %q, %q, %d, want empty", f.Function(), f.File(), f.Line())
			}
			text, err := f.MarshalText()
			if f.String() != "???" || string(text) != "???" || err != nil || f.LogValue().String() != "???" {
				t.Errorf("String, MarshalText, LogValue = %q, %q, %v, %q, want ???", f.String(), text, err, f.LogValue())
			}
		})
	}
}

func TestSameCallSiteSameFrame(t *testing.T) {
	var frames []trace.Frame
	for range 3 {
		frames = append(frames, here())
	}
	if frames[0] != frames[1] || frames[1] != frames[2] {
		t.Errorf("frames from one call site = %v, want them equal", frames)
	}
	if other := inClosure(); other == frames[0] {
		t.Errorf("frames from two call sites are equal: %v", other)
	}
}

func TestAppendText(t *testing.T) {
	b, err := here().AppendText([]byte("at "))
	if string(b) != "at frames_test.go:11 trace/v2_test.here" || err != nil {
		t.Errorf("AppendText() = %q, %v, want %q, nil", b, err, "at frames_test.go:11 trace/v2_test.here")
	}
}

// Without MarshalText, encoding/json would write the program counter as a number.
func TestJSON(t *testing.T) {
	b, err := json.Marshal(struct{ At trace.Frame }{here()})
	if want := `{"At":"frames_test.go:11 trace/v2_test.here"}`; string(b) != want || err != nil {
		t.Errorf("json.Marshal() = %s, %v, want %s, nil", b, err, want)
	}
}

func TestErrorf(t *testing.T) {
	errNotFound := errors.New("not found")
	err := fmt.Errorf("%s: %w", here(), errNotFound)
	if want := "frames_test.go:11 trace/v2_test.here: not found"; err.Error() != want {
		t.Errorf("Error() = %q, want %q", err, want)
	}
	if !errors.Is(err, errNotFound) {
		t.Errorf("errors.Is(%v, errNotFound) = false", err)
	}
}

// withoutTime drops the time attribute, so the output can be compared.
func withoutTime(groups []string, a slog.Attr) slog.Attr {
	if a.Key == slog.TimeKey && len(groups) == 0 {
		return slog.Attr{}
	}
	return a
}

func TestSlogAttr(t *testing.T) {
	tests := []struct {
		name    string
		handler func(*bytes.Buffer) slog.Handler
		want    string
	}{
		{"text", func(b *bytes.Buffer) slog.Handler {
			return slog.NewTextHandler(b, &slog.HandlerOptions{ReplaceAttr: withoutTime})
		}, `level=INFO msg=saved at="frames_test.go:11 trace/v2_test.here"` + "\n"},
		{"JSON", func(b *bytes.Buffer) slog.Handler {
			return slog.NewJSONHandler(b, &slog.HandlerOptions{ReplaceAttr: withoutTime})
		}, `{"level":"INFO","msg":"saved","at":"frames_test.go:11 trace/v2_test.here"}` + "\n"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var out bytes.Buffer
			slog.New(tt.handler(&out)).Info("saved", "at", here())
			if out.String() != tt.want {
				t.Errorf("logged %q, want %q", out.String(), tt.want)
			}
		})
	}
}

// recorder keeps the records it handles.
type recorder struct{ records []slog.Record }

func (*recorder) Enabled(context.Context, slog.Level) bool { return true }
func (h *recorder) Handle(_ context.Context, r slog.Record) error {
	h.records = append(h.records, r)
	return nil
}
func (h *recorder) WithAttrs([]slog.Attr) slog.Handler { return h }
func (h *recorder) WithGroup(string) slog.Handler      { return h }

// A Frame and a slog.Record's PC are the same kind of program counter, both ways round.
func TestSlogRecordPC(t *testing.T) {
	h := new(recorder)
	logHere(slog.New(h))
	if got, want := trace.Frame(h.records[0].PC).String(), "frames_test.go:41 trace/v2_test.logHere"; got != want {
		t.Errorf("Frame(Record.PC) = %q, want %q", got, want)
	}

	_, caller := helped()
	r := slog.NewRecord(time.Time{}, slog.LevelInfo, "logged for the caller", uintptr(caller))
	if src := r.Source(); src.Function != pkg+".helped" || src.Line != 26 {
		t.Errorf("Record.Source() = %s:%d, want helped's line 26", src.Function, src.Line)
	}
}

// Goroutines that format new frames at once agree on what they are.
func TestConcurrentFirstUse(t *testing.T) {
	sites := []func() trace.Frame{firstUseA, firstUseB, firstUseC}
	want := []string{
		"frames_test.go:44 trace/v2_test.firstUseA",
		"frames_test.go:45 trace/v2_test.firstUseB",
		"frames_test.go:46 trace/v2_test.firstUseC",
	}
	got := make([][]string, 8)
	var wg sync.WaitGroup
	for g := range got {
		got[g] = make([]string, len(sites))
		wg.Go(func() {
			for i, site := range sites {
				got[g][i] = site().String()
			}
		})
	}
	wg.Wait()

	for g := range got {
		if !slices.Equal(got[g], want) {
			t.Errorf("goroutine %d got %q, want %q", g, got[g], want)
		}
	}
}
