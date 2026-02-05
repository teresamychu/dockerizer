package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestParseGoMod_ModuleName(t *testing.T) {
	dir := t.TempDir()
	gomod := `module github.com/user/myapp

go 1.22
`
	os.WriteFile(filepath.Join(dir, "go.mod"), []byte(gomod), 0644)

	module, _, _ := ParseGoMod(filepath.Join(dir, "go.mod"))

	if module != "myapp" {
		t.Errorf("expected module name 'myapp', got '%s'", module)
	}
}

func TestParseGoMod_GoVersion(t *testing.T) {
	dir := t.TempDir()
	gomod := `module github.com/user/myapp

go 1.21
`
	os.WriteFile(filepath.Join(dir, "go.mod"), []byte(gomod), 0644)

	_, version, _ := ParseGoMod(filepath.Join(dir, "go.mod"))

	if version != "1.21" {
		t.Errorf("expected version '1.21', got '%s'", version)
	}
}

func TestParseGoMod_DetectsPostgres(t *testing.T) {
	dir := t.TempDir()
	gomod := `module github.com/user/myapp

go 1.22

require github.com/lib/pq v1.10.9
`
	os.WriteFile(filepath.Join(dir, "go.mod"), []byte(gomod), 0644)

	_, _, services := ParseGoMod(filepath.Join(dir, "go.mod"))

	found := false
	for _, svc := range services {
		if svc.Name == "postgres" {
			found = true
			break
		}
	}
	if !found {
		t.Error("expected to detect postgres service")
	}
}

func TestParseGoMod_DetectsRedis(t *testing.T) {
	dir := t.TempDir()
	gomod := `module github.com/user/myapp

go 1.22

require github.com/go-redis/redis/v8 v8.11.5
`
	os.WriteFile(filepath.Join(dir, "go.mod"), []byte(gomod), 0644)

	_, _, services := ParseGoMod(filepath.Join(dir, "go.mod"))

	found := false
	for _, svc := range services {
		if svc.Name == "redis" {
			found = true
			break
		}
	}
	if !found {
		t.Error("expected to detect redis service")
	}
}

func TestParseGoMod_DetectsMultipleServices(t *testing.T) {
	dir := t.TempDir()
	gomod := `module github.com/user/myapp

go 1.25.6

require (
	github.com/lib/pq v1.10.9
	github.com/go-redis/redis/v8 v8.11.5
)
`
	os.WriteFile(filepath.Join(dir, "go.mod"), []byte(gomod), 0644)

	_, _, services := ParseGoMod(filepath.Join(dir, "go.mod"))

	if len(services) != 2 {
		t.Errorf("expected 2 services, got %d", len(services))
	}
}

func TestParseGoMod_NoServices(t *testing.T) {
	dir := t.TempDir()
	gomod := `module github.com/user/myapp

go 1.22

require github.com/gin-gonic/gin v1.9.1
`
	os.WriteFile(filepath.Join(dir, "go.mod"), []byte(gomod), 0644)

	_, _, services := ParseGoMod(filepath.Join(dir, "go.mod"))

	if len(services) != 0 {
		t.Errorf("expected 0 services, got %d", len(services))
	}
}

func TestScanEnvFiles(t *testing.T) {
	dir := t.TempDir()
	os.WriteFile(filepath.Join(dir, ".env"), []byte("PORT=3000\nDB_HOST=localhost"), 0644)

	envVars := ScanEnvFiles(dir)

	if envVars["PORT"] != "3000" {
		t.Errorf("expected PORT=3000, got %s", envVars["PORT"])
	}
}

func TestScanConfig_ExtractsPort(t *testing.T) {
	dir := t.TempDir()
	os.WriteFile(filepath.Join(dir, ".env"), []byte("PORT=9000"), 0644)

	config := ScanConfig(dir)

	if config.Port != "9000" {
		t.Errorf("expected port 9000, got %s", config.Port)
	}
}
