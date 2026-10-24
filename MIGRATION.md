# Migrating from v1 to v2

v1's `Here` and `FromParent` returned a string, built from flags. v2's return a `Frame`: the program counter of the
call, which becomes text when you print it. The flags are gone; there is one text form, and methods for each part.

## 1. Update the import path

```sh
go get github.com/vardius/trace/v2
```

```go
import "github.com/vardius/trace/v2"
```

v2 needs Go 1.26 or newer.

## 2. Before and after

```go
// v1
return fmt.Errorf("%s: %w", trace.Here(trace.Lfile|trace.Lline|trace.Lfunction), err)
```

```go
// v2
return fmt.Errorf("%s: %w", trace.Here(), err)
```

`%s` and `%v` print a `Frame` with its `String` method, so a frame passed to `fmt` needs no change beyond the flags.
Where you need a `string`, call `String`:

```go
// v1
at := trace.FromParent(1, trace.LstdFlags)

// v2
at := trace.Caller(1).String()
```

## 3. Every call

| v1 | v2 |
| --- | --- |
| `Here(flags int) string` | `Here() Frame` |
| `FromParent(calldepth, flags int) string` | `Caller(skip int) Frame`, with the same number: `FromParent(1, …)` is `Caller(1)` |
| `Here(Lfile \| Lline \| Lfunction)` | `Here().String()`, in v2's format (below) |
| `Here(LstdFlags)`, `Here(0)` | `Here().String()`, which has the line too |
| `Here(Lfile)` | `Here().File()` |
| `Here(Lline)` | `Here().Line()`, an `int` |
| `Here(Lfunction)` | `Here().Function()` |
| Any other set of flags | Build it from `File`, `Line` and `Function` |
| `Lfile`, `Lline`, `Lfunction`, `LstdFlags` | Removed |

Call a method on a frame you keep, rather than calling `Here` once per part: two `Here` calls on different lines
are different frames.

## 4. Behaviour changes

- **One format.** `String` is `file:line function`, with the file's base name and the function's package path cut to its
  last element: `service.go:42 users.(*Service).Create`. v1 printed the full path and the full function name:
  `/home/me/src/users/service.go:42 github.com/acme/users.(*Service).Create`. `File` and `Function` still return those.
- **The default has the line.** v1's `LstdFlags`, which a flags value of 0 also meant, left the line out.
- **No stray separators.** v1's `Here(Lline)` was `":42"` and `Here(Lfunction)` was `" pkg.Func"`, with the separator
  meant for the part before them. `Line` is the number and `Function` the name, with nothing around them.
- **Beyond the stack.** v1 returned `""` when the depth was more than the stack. v2's `Caller` returns the zero `Frame`,
  which prints `???`; compare with `0` to check for it.
- **A negative depth counts as 0.** v1's `FromParent(-1, …)` reported a line inside trace's own `FromParent`, and
  `FromParent(math.MaxInt, …)` overflowed and reported `runtime.Caller`. v2's `Caller(-1)` is `Caller(0)`, and
  `Caller(math.MaxInt)` is the zero `Frame`.
- **Encoders write text.** A `Frame` implements `encoding.TextMarshaler`, so `encoding/json` writes it as its `String`.
  A `Frame` is an integer underneath; converting one to `uintptr` gives the program counter, which `slog.NewRecord` takes.
- **Resolved once per call site.** The first time a frame is used, v2 resolves it and keeps the result for as long as
  the program runs, so later uses cost no allocation. That memory is one entry per call site used.
