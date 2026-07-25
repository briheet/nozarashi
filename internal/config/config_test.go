package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestParseTOMLConfigPreservesEnvironmentKeyCase(t *testing.T) {
	configPath := filepath.Join(t.TempDir(), "nozarashi.toml")
	configData := []byte(`
[project]
name = "example"
version = "0.1.0"
description = "Environment key test"
profile = "test"

[inputs.nixpkgs]
type = "flake"
source = "github:NixOS/nixpkgs"

[services.postgres]
type = "oci"
reference = "postgres:17-alpine"
environment = { POSTGRES_USER = "nozarashi" }

[volumes.data]
driver = "local"

[networks.default]
driver = "default"
`)

	if err := os.WriteFile(configPath, configData, 0o600); err != nil {
		t.Fatalf("write configuration: %v", err)
	}

	project, err := ParseTOMLConfig(t.Context(), configPath)
	if err != nil {
		t.Fatalf("parse configuration: %v", err)
	}

	if project.Services["postgres"].Environment["POSTGRES_USER"] != "nozarashi" {
		t.Fatalf("environment key casing was not preserved: %#v", project.Services["postgres"].Environment)
	}
}

func TestParseSimpleDockerExample(t *testing.T) {
	configPath := filepath.Join("..", "..", "examples", "simple_docker", "nozarashi.toml")

	project, err := ParseTOMLConfig(t.Context(), configPath)
	if err != nil {
		t.Fatalf("parse simple Docker example: %v", err)
	}

	if project.Services["postgres"].Environment["POSTGRES_USER"] != "nozarashi" {
		t.Fatalf("PostgreSQL environment is incomplete")
	}
	if project.Services["backend"].Environment["REDIS_ADDRESS"] == "" {
		t.Fatalf("backend Redis address is missing")
	}
}
