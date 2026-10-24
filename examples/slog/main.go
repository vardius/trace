// Command slog logs through two helpers, and its log lines still point at the code that called them.
//
//	go run ./examples/slog
package main

import (
	"context"
	"log/slog"
	"os"
	"time"

	"github.com/vardius/trace/v2"
)

// audit logs msg for the code that calls it: the record's source is the caller's line, not audit's.
func audit(ctx context.Context, logger *slog.Logger, msg string, args ...any) {
	if !logger.Enabled(ctx, slog.LevelInfo) {
		return
	}
	r := slog.NewRecord(time.Now(), slog.LevelInfo, msg, uintptr(trace.Caller(1)))
	r.Add(args...)
	_ = logger.Handler().Handle(ctx, r)
}

// timed logs how long the caller took, and where it was, when the function it returns is called:
//
//	defer timed(logger, "load user")()
func timed(logger *slog.Logger, name string) func() {
	at, start := trace.Caller(1), time.Now()
	return func() {
		logger.Info(name, "took", time.Since(start), "at", at)
	}
}

func loadUser(ctx context.Context, logger *slog.Logger, id int) {
	defer timed(logger, "load user")()

	audit(ctx, logger, "user loaded", "id", id)
}

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{AddSource: true}))
	loadUser(context.Background(), logger, 42)
}
