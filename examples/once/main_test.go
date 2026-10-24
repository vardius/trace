package main

import (
	"bytes"
	"testing"
)

// A test binary names package main by its import path, so this test is once.TestWarnOnce, not main.TestWarnOnce.
func TestWarnOnce(t *testing.T) {
	warned.Clear() // left over from an earlier run with -count
	var out bytes.Buffer
	deprecated := func() { warnOnce(&out, "deprecated") }

	for range 3 {
		deprecated() // line 15
	}
	deprecated() // line 17

	want := "main_test.go:15 once.TestWarnOnce: deprecated\nmain_test.go:17 once.TestWarnOnce: deprecated\n"
	if out.String() != want {
		t.Errorf("warned %q, want %q", out.String(), want)
	}
}
