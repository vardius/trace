package trace_test

import (
	"log/slog"

	"github.com/vardius/trace/v2"
)

// The tests expect these functions at these lines: add new ones at the end of the file.

func here() trace.Frame { return trace.Here() }

func inClosure() trace.Frame {
	return func() trace.Frame { return trace.Here() }()
}

type service struct{}

func (*service) method() trace.Frame { return trace.Here() }

func generic[T any]() trace.Frame { return trace.Here() }

// helper is small enough to inline. Caller(0) is the line in it, Caller(1) the line that called it.
func helper(skip int) trace.Frame { return trace.Caller(skip) }

func helped() (self, caller trace.Frame) { return helper(0), helper(1) }

//go:noinline
func helperNoInline(skip int) trace.Frame { return trace.Caller(skip) }

func helpedNoInline() (self, caller trace.Frame) { return helperNoInline(0), helperNoInline(1) }

// skips calls Caller(skip) two calls below the line that calls skips.
func skips(skip int) trace.Frame { return skipsInner(skip) }

//go:noinline
func skipsInner(skip int) trace.Frame { return trace.Caller(skip) }

func skipsFromHere(skip int) trace.Frame { return skips(skip) }

func logHere(l *slog.Logger) { l.Info("logged here") }

// Only TestConcurrentFirstUse calls these, so their frames are resolved there first.
func firstUseA() trace.Frame { return trace.Here() }
func firstUseB() trace.Frame { return trace.Here() }
func firstUseC() trace.Frame { return trace.Here() }
