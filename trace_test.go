package trace

import (
	"testing"
)

func TestFlags(t *testing.T) {
	ref := Here(Lfile | Lfunction)
	if ref == "" {
		t.Fatalf("Did not trace the call (%s)", ref)
	}

	parentRef := FromParent(1, Lfile|Lfunction)
	if parentRef == "" {
		t.Fatalf("Did not trace the call (%s)", parentRef)
	}
}

func TestFromParent(t *testing.T) {
	type args struct {
		calldepth int
		flags     int
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"trace with one flag", args{0, Lline}, ":34"},
		{"trace with multiple flags", args{0, Lline | Lfunction}, ":34 testing.tRunner"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := FromParent(tt.args.calldepth, tt.args.flags); got != tt.want {
				t.Errorf("FromParent() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHere(t *testing.T) {
	type args struct {
		flags int
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"trace with one flag", args{Lline}, ":55"},
		{"trace with multiple flags", args{Lline | Lfunction}, ":55 github.com/vardius/trace.TestHere.func1"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Here(tt.args.flags); got != tt.want {
				t.Errorf("Here() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTraceDefaultFlags(t *testing.T) {
	ref := Here(0)
	if ref == "" {
		t.Fatalf("Did not trace the call (%s)", ref)
	}
}
