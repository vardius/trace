package trace_test

import (
	"errors"
	"fmt"
	"io"
	"log/slog"
	"testing"

	"github.com/vardius/trace/v2"
)

var (
	frameSink  trace.Frame
	stringSink string
	intSink    int
)

// TestAllocations guards the allocations per call. Lower a number when an optimisation lands;
// never raise one without a reason in the commit.
//
// A frame's first use resolves it and keeps the result, which allocates once per call site;
// AllocsPerRun's warm-up run pays that, so the counts are those of every use after the first.
func TestAllocations(t *testing.T) {
	if raceEnabled {
		t.Skip("the race detector changes allocation counts")
	}

	f := here()
	buf := make([]byte, 0, 128)
	errNotFound := errors.New("not found")
	logger := slog.New(slog.NewJSONHandler(io.Discard, nil))

	tests := []struct {
		name string
		call func()
		want float64
	}{
		{"Here", func() { frameSink = trace.Here() }, 0},
		{"Caller(1)", func() { frameSink = trace.Caller(1) }, 0},
		{"String", func() { stringSink = f.String() }, 0},
		{"AppendText, into a buffer with room", func() { buf, _ = f.AppendText(buf[:0]) }, 0},
		{"Function, File and Line", func() { stringSink, stringSink, intSink = f.Function(), f.File(), f.Line() }, 0},
		{"Here().String()", func() { stringSink = trace.Here().String() }, 0},
		{"LogValue", func() { stringSink = f.LogValue().String() }, 0},
		// Never kept, so the warm-up run doesn't pay for it: it has to be free every time.
		{"String of the zero Frame", func() { stringSink = trace.Frame(0).String() }, 0},
		// The slice it returns.
		{"MarshalText", func() { _, _ = f.MarshalText() }, 1},
		// fmt's: the Frame in an interface, the message, and the error that wraps it.
		{"fmt.Errorf with Here", func() { _ = fmt.Errorf("%s: %w", trace.Here(), errNotFound) }, 3},
		// The Frame in the interface of slog's ...any argument. LogValue and the JSON handler allocate nothing.
		{"slog attribute", func() { logger.Info("saved", "at", trace.Here()) }, 1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := testing.AllocsPerRun(1000, tt.call); got != tt.want {
				t.Errorf("allocates %v times per call, want %v", got, tt.want)
			}
		})
	}
}
