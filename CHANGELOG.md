# Changelog

## v2.0.0 (2026-10-24)

`Here` and `Caller` return a `Frame`, which becomes text when it is printed, in place of a string built from flags.
See [MIGRATION.md](MIGRATION.md) for the before and after of every call.

### Added

- `Frame`, with `Function`, `File`, `Line` and `String`. Frames are comparable, so they can be map keys.
- `Caller(skip)`, in place of `FromParent`.
- `Frame` implements `slog.LogValuer`, `encoding.TextAppender` and `encoding.TextMarshaler`, and converts to and from
  the program counter in `slog.Record.PC`.
- Runnable examples in `examples/`, each with a test: tagging errors, logging for the caller with `slog`, and warning
  once per call site.

### Changed

- Module path is `github.com/vardius/trace/v2`; Go 1.26 or newer is required.
- One text form, `file:line function`, with the file's base name and the function's package path cut to its last element.
- Frames are resolved with `runtime.CallersFrames`, the runtime's documented way to account for inlining,
  in place of `runtime.FuncForPC`.
- Faster, with no allocations: `Here` as a string went from 437 ns and 9 allocations to 72.7 ns and none, and from
  295 ns to 6.3 ns with 16 callers at once. Allocation counts are guarded by tests.

### Fixed

- A negative depth reported a line inside trace itself; it now counts as 0.
- A depth near `math.MaxInt` overflowed and reported `runtime.Caller`; it now gives the zero `Frame`.
- Single flags left a separator in front: `":42"` for `Lline`, `" pkg.Func"` for `Lfunction`.

### Removed

- `FromParent` and the flags `Lfile`, `Lline`, `Lfunction` and `LstdFlags`.
