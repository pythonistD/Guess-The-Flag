package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadComposeEnv(t *testing.T) {
	path := filepath.Join("..", "..", "..", "config.yml")
	if _, err := os.Stat(path); err != nil {
		t.Skip("config.yml not in repo root", err)
	}
	env, err := LoadComposeEnv(path)
	if err != nil {
		t.Fatal(err)
	}
	if env.PostgresUser != "postgres" {
		t.Fatalf("user: got %q", env.PostgresUser)
	}
	if env.Domain != "flags.example.com" {
		t.Fatalf("domain: got %q", env.Domain)
	}
}
