package trace

import (
	"testing"
)

func TestTrace(t *testing.T) {
	ref := Here(Lfile | Lfunction)
	if ref == "" {
		t.Fatalf("Did not trace the call (%s)", ref)
	}
}
