package main

import (
	"bytes"
	"log/slog"
	"path/filepath"
	"strconv"
	"testing"
)

// stable drops what changes from run to run and machine to machine: the time, the duration and the source's directory.
func stable(groups []string, a slog.Attr) slog.Attr {
	switch a.Key {
	case slog.TimeKey, "took":
		return slog.Attr{}
	case slog.SourceKey:
		src := a.Value.Any().(*slog.Source)
		return slog.String(a.Key, filepath.Base(src.File)+":"+strconv.Itoa(src.Line))
	}
	return a
}

func TestLoadUser(t *testing.T) {
	var out bytes.Buffer
	logger := slog.New(slog.NewTextHandler(&out, &slog.HandlerOptions{AddSource: true, ReplaceAttr: stable}))

	loadUser(t.Context(), logger, 42)

	// audit's record has its caller's line, 38. timed logs from its own line, 31, with the line that deferred it, 36,
	// as the at attribute. A test binary names package main by its import path, so loadUser is slog.loadUser here.
	want := `level=INFO source=main.go:38 msg="user loaded" id=42` + "\n" +
		`level=INFO source=main.go:31 msg="load user" at="main.go:36 slog.loadUser"` + "\n"
	if out.String() != want {
		t.Errorf("logged\n%s\nwant\n%s", out.String(), want)
	}
}
