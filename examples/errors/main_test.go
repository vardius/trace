package main

import (
	"encoding/json"
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	path := filepath.Join(t.TempDir(), "config.json")
	if err := os.WriteFile(path, []byte(`{"addr": ":8080"}`), 0o600); err != nil {
		t.Fatal(err)
	}

	cfg, err := loadConfig(path)
	if cfg.Addr != ":8080" || err != nil {
		t.Errorf("loadConfig() = %+v, %v, want :8080, nil", cfg, err)
	}
}

// A test binary names package main by its import path, so the functions here are errors.loadConfig, not main.loadConfig.

func TestLoadConfigMissing(t *testing.T) {
	_, err := loadConfig(filepath.Join(t.TempDir(), "missing.json"))
	if !errors.Is(err, fs.ErrNotExist) || !strings.HasPrefix(err.Error(), "main.go:23 errors.loadConfig: open ") {
		t.Errorf("loadConfig() error = %v, want it to say main.go:23 and wrap fs.ErrNotExist", err)
	}
}

func TestLoadConfigMalformed(t *testing.T) {
	path := filepath.Join(t.TempDir(), "config.json")
	if err := os.WriteFile(path, []byte(`{"addr":`), 0o600); err != nil {
		t.Fatal(err)
	}

	_, err := loadConfig(path)
	var syntax *json.SyntaxError
	if !errors.As(err, &syntax) || !strings.HasPrefix(err.Error(), "main.go:28 errors.loadConfig: decode ") {
		t.Errorf("loadConfig() error = %v, want it to say main.go:28 and wrap a *json.SyntaxError", err)
	}
}
