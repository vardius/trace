package trace

import "testing"

// BenchmarkFirstUse measures a frame's first use: resolving it with runtime.CallersFrames and keeping the result.
// It empties the cache every time, which costs a little too.
func BenchmarkFirstUse(b *testing.B) {
	f := Here()
	b.ReportAllocs()
	for b.Loop() {
		locations.Clear()
		_ = f.String()
	}
}
