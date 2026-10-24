package trace_test

import (
	"errors"
	"fmt"
	"io"
	"log/slog"
	"testing"

	"github.com/vardius/trace/v2"
)

// Every benchmark formats frames it has formatted before, as a program does after a call site's first use.
// BenchmarkFirstUse, in cache_test.go, measures that first use.

func BenchmarkHere(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		frameSink = trace.Here()
	}
}

func BenchmarkCaller(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		frameSink = trace.Caller(1)
	}
}

func BenchmarkHereString(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		stringSink = trace.Here().String()
	}
}

func BenchmarkHereAppendText(b *testing.B) {
	buf := make([]byte, 0, 128)
	b.ReportAllocs()
	for b.Loop() {
		buf, _ = trace.Here().AppendText(buf[:0])
	}
}

func BenchmarkHereStringParallel(b *testing.B) {
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		var s string
		for pb.Next() {
			s = trace.Here().String()
		}
		stringSink = s
	})
}

func BenchmarkString(b *testing.B) {
	f := trace.Here()
	b.ReportAllocs()
	for b.Loop() {
		stringSink = f.String()
	}
}

func BenchmarkAppendText(b *testing.B) {
	f := trace.Here()
	buf := make([]byte, 0, 128)
	b.ReportAllocs()
	for b.Loop() {
		buf, _ = f.AppendText(buf[:0])
	}
}

func BenchmarkErrorf(b *testing.B) {
	errNotFound := errors.New("not found")
	b.ReportAllocs()
	for b.Loop() {
		_ = fmt.Errorf("%s: %w", trace.Here(), errNotFound)
	}
}

func BenchmarkSlog(b *testing.B) {
	for _, h := range []struct {
		name    string
		handler slog.Handler
	}{
		{"text", slog.NewTextHandler(io.Discard, nil)},
		{"JSON", slog.NewJSONHandler(io.Discard, nil)},
	} {
		b.Run(h.name, func(b *testing.B) {
			logger := slog.New(h.handler)
			b.ReportAllocs()
			for b.Loop() {
				logger.Info("saved", "at", trace.Here())
			}
		})
	}
}
