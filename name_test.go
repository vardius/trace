package trace

import "testing"

// Frames can only come from code in this module, so the names of code elsewhere are checked here, directly.
func TestShortName(t *testing.T) {
	tests := []struct{ function, want string }{
		{"main.main", "main.main"},
		{"github.com/acme/users.(*Service).Create", "users.(*Service).Create"},
		{"github.com/acme/users.(*Service).Create.func1", "users.(*Service).Create.func1"},
		{"github.com/acme/users.Map[...]", "users.Map[...]"},
		{"github.com/vardius/trace/v2.Here", "trace/v2.Here"},
		{"github.com/vardius/gorouter/v5.(*Router).ServeHTTP", "gorouter/v5.(*Router).ServeHTTP"},
		{"github.com/acme/foo/v10.F", "foo/v10.F"},
		{"github.com/vardius/trace/v2_test.TestHere", "trace/v2_test.TestHere"},
		{"k8s.io/api/core/v1.(*Pod).Reset", "core/v1.(*Pod).Reset"},
		{"example.com/v2.F", "example.com/v2.F"},
		// Versions that aren't a path element of their own.
		{"gopkg.in/yaml.v3.Unmarshal", "yaml.v3.Unmarshal"},
		{"github.com/rogchap/v8go.NewIsolate", "v8go.NewIsolate"},
		{"github.com/acme/v.F", "v.F"},
		{"github.com/acme/vendor.F", "vendor.F"},
	}
	for _, tt := range tests {
		if got := shortName(tt.function); got != tt.want {
			t.Errorf("shortName(%q) = %q, want %q", tt.function, got, tt.want)
		}
	}
}
