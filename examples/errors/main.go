// Command errors loads a config file and, when that fails, says where each error on the way up was wrapped.
//
//	go run ./examples/errors missing.json
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/vardius/trace/v2"
)

type Config struct {
	Addr string `json:"addr"`
}

// loadConfig reads and decodes the config at path. Its errors say which line they came from.
func loadConfig(path string) (Config, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return Config{}, fmt.Errorf("%s: %w", trace.Here(), err)
	}

	var cfg Config
	if err := json.Unmarshal(b, &cfg); err != nil {
		return Config{}, fmt.Errorf("%s: decode %s: %w", trace.Here(), path, err)
	}
	return cfg, nil
}

func main() {
	log.SetFlags(0)
	if len(os.Args) != 2 {
		log.Fatal("usage: errors config.json")
	}

	cfg, err := loadConfig(os.Args[1])
	if err != nil {
		log.Fatal(fmt.Errorf("%s: %w", trace.Here(), err))
	}
	fmt.Println("listening on", cfg.Addr)
}
