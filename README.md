# 📍 trace

[![Test](https://github.com/vardius/trace/actions/workflows/test.yml/badge.svg)](https://github.com/vardius/trace/actions/workflows/test.yml)
[![Go Reference](https://pkg.go.dev/badge/github.com/vardius/trace/v2.svg)](https://pkg.go.dev/github.com/vardius/trace/v2)
[![Go Report Card](https://goreportcard.com/badge/github.com/vardius/trace/v2)](https://goreportcard.com/report/github.com/vardius/trace/v2)
[![License: MIT](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE.md)

Tell where a Go call came from, the function, file and line, in error messages and log lines.

- **Here** captures the line that calls it, and **Caller** a line further up the stack, for helpers.
- **Lazy.** A `Frame` is a program counter. It becomes `service.go:42 users.(*Service).Create` only when it is printed.
- **No allocations.** Capturing a frame and turning it into text, in a buffer, a string or a `slog` value,
  allocate nothing once its call site has been used, and tests guard every count.
- **Right with inlining.** Frames are resolved with `runtime.CallersFrames`, so a function the compiler inlined is still the one named.
- **Fits the standard library.** `%s` in `fmt`, attributes and `Record.PC` in `log/slog`, text in `encoding/json`.

📖 **[Documentation](https://vardius.github.io/trace/)** · [API reference](https://pkg.go.dev/github.com/vardius/trace/v2) · [Examples](examples) · [Migrating from v1](MIGRATION.md)

## Install

```sh
go get github.com/vardius/trace/v2
```

Requires Go 1.26 or newer.

## Tag errors

`Here` returns the frame of the line that calls it. Put it in front of an error as you wrap it:

```go
func (s *Service) Create(ctx context.Context, u User) error {
	if err := s.repo.Insert(ctx, u); err != nil {
		return fmt.Errorf("%s: %w", trace.Here(), err)
	}
	return nil
}
```

```txt
service.go:42 users.(*Service).Create: duplicate key
```

The text is `file:line function`: the file's base name and the function's package path cut to its last element.
`Function`, `File` and `Line` return the full name, the full path and the line on their own.

## Report the caller

`Caller(1)` is the line that called the function calling it, so a helper can name its caller instead of itself:

```go
// wrap tags err with the line that called wrap.
func wrap(err error) error {
	return fmt.Errorf("%s: %w", trace.Caller(1), err)
}
```

`Caller(0)` is the same as `Here()`. When the stack is less deep than asked, `Caller` returns the zero `Frame`, which prints `???`.

## Log with slog

A `Frame` is a `slog.LogValuer`: it logs as its text, which costs no allocation, and only if the line is written.

```go
logger.Warn("slow query", "at", trace.Caller(1))
```

It holds the same program counter as `slog.Record.PC`. A helper that logs for its caller hands it to `slog.NewRecord`,
so `AddSource` reports the caller's line, not the helper's:

```go
r := slog.NewRecord(time.Now(), slog.LevelInfo, msg, uintptr(trace.Caller(1)))
```

## Frames as keys

Frames are comparable: frames captured by the same call are equal, so a map keyed by them does something once per call site
(each copy of an inlined function, and each shape a generic function is compiled for, has call sites of its own).
[examples/once](examples/once/main.go) warns once about each line that calls a deprecated function.

## Benchmarks

Time and allocations per call, on an Apple M4 Max with Go 1.27.0, median of 5 runs. Both versions ran on the same harness;
v1 with all three flags, its output closest to v2's `String`. With 16 callers at once, the time is wall time per call
across all of them, so it measures throughput. Run v2's own suite with `go test -bench .`.

| Scenario | v1 | v2 |
| --- | ---: | ---: |
| `Here`, as a string | 437 ns · 9 allocs | **72.7 ns · 0 allocs** |
| `Caller(1)`, as a string | 450 ns · 9 allocs | **93.5 ns · 0 allocs** |
| `Here`, into a buffer | 588 ns · 9 allocs | **76.0 ns · 0 allocs** |
| `Here`, as a string, 16 callers at once | 295 ns · 9 allocs | **6.30 ns · 0 allocs** |
| `fmt.Errorf` with `Here` | 542 ns · 12 allocs | **166 ns · 3 allocs** |
| `slog` text line with `Here` | 1.68 µs · 10 allocs | **562 ns · 1 alloc** |
| `slog` JSON line with `Here` | 1.11 µs · 10 allocs | **431 ns · 1 alloc** |

What v2 still spends is `runtime.Callers` walking up the stack. The first time a frame is used, it is resolved with
`runtime.CallersFrames` and the result kept, one per call site, so every later use is a lookup. The allocations left in the
last three rows are `fmt`'s message and error, and the one that puts a `Frame` in an `any` argument.

## Upgrading from v1

`Here` and `FromParent` returned a string built from flags. v2 returns a `Frame`, with one text form:

```go
// v1
return fmt.Errorf("%s: %w", trace.Here(trace.Lfile|trace.Lline|trace.Lfunction), err)

// v2
return fmt.Errorf("%s: %w", trace.Here(), err)
```

`FromParent(n, flags)` is `Caller(n)`. [MIGRATION.md](MIGRATION.md) maps every v1 call and lists the behaviour changes.

## Contributing

Bug reports and pull requests are welcome: see [CONTRIBUTING.md](.github/CONTRIBUTING.md).

## License

MIT, see [LICENSE.md](LICENSE.md).
