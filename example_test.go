package trace_test

import (
	"errors"
	"fmt"
	"log/slog"
	"os"

	"github.com/vardius/trace/v2"
)

func ExampleHere() {
	fmt.Println(trace.Here())
	// Output: example_test.go:13 trace/v2_test.ExampleHere
}

func ExampleHere_error() {
	errNotFound := errors.New("user not found")

	err := fmt.Errorf("%s: %w", trace.Here(), errNotFound)

	fmt.Println(err)
	fmt.Println(errors.Is(err, errNotFound))
	// Output:
	// example_test.go:20 trace/v2_test.ExampleHere_error: user not found
	// true
}

func ExampleCaller() {
	// wrap tags err with the line that called it, not with its own.
	wrap := func(err error) error {
		return fmt.Errorf("%s: %w", trace.Caller(1), err)
	}

	err := wrap(errors.New("connection refused"))

	fmt.Println(err)
	// Output: example_test.go:35 trace/v2_test.ExampleCaller: connection refused
}

func ExampleFrame_accessors() {
	f := trace.Here()

	fmt.Println(f.Function())
	fmt.Println(f.Line())
	// Output:
	// github.com/vardius/trace/v2_test.ExampleFrame_accessors
	// 42
}

func ExampleFrame_AppendText() {
	buf := make([]byte, 0, 128)

	buf, _ = trace.Here().AppendText(buf)

	fmt.Println(string(buf))
	// Output: example_test.go:54 trace/v2_test.ExampleFrame_AppendText
}

func ExampleFrame_LogValue() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{ReplaceAttr: withoutTime}))

	logger.Info("cache miss", "at", trace.Here())
	// Output: level=INFO msg="cache miss" at="example_test.go:63 trace/v2_test.ExampleFrame_LogValue"
}
