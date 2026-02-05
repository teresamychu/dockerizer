package generators

import (
	"strings"
	"testing"
)

func TestGoGenerator_Generate_BinaryName(t *testing.T) {
	g := &GoGenerator{
		AppInfo: AppInfo{
			BinaryName: "myapp",
			GoVersion:  "1.22",
			HasGoSum:   true,
		},
	}

	output, err := g.GenerateDockerfile(".")
	if err != nil {
		t.Fatalf("GenerateDockerfile failed: %v", err)
	}

	if !strings.Contains(output, "myapp") {
		t.Errorf("expected binary name 'myapp' in output, got:\n%s", output)
	}
}

func TestGoGenerator_Generate_GoVersion(t *testing.T) {
	g := &GoGenerator{
		AppInfo: AppInfo{
			BinaryName: "app",
			GoVersion:  "1.22",
			HasGoSum:   true,
		},
	}

	output, err := g.GenerateDockerfile(".")
	if err != nil {
		t.Fatalf("GenerateDockerfile failed: %v", err)
	}

	if !strings.Contains(output, "golang:1.22-alpine") {
		t.Errorf("expected 'golang:1.22-alpine' in output, got:\n%s", output)
	}
}

func TestGoGenerator_Generate_WithGoSum(t *testing.T) {
	g := &GoGenerator{
		AppInfo: AppInfo{
			BinaryName: "app",
			GoVersion:  "1.22",
			HasGoSum:   true,
		},
	}

	output, err := g.GenerateDockerfile(".")
	if err != nil {
		t.Fatalf("GenerateDockerfile failed: %v", err)
	}

	if !strings.Contains(output, "COPY go.mod go.sum ./") {
		t.Errorf("expected 'COPY go.mod go.sum ./' when HasGoSum=true, got:\n%s", output)
	}
}

func TestGoGenerator_Generate_WithoutGoSum(t *testing.T) {
	g := &GoGenerator{
		AppInfo: AppInfo{
			BinaryName: "app",
			GoVersion:  "1.22",
			HasGoSum:   false,
		},
	}

	output, err := g.GenerateDockerfile(".")
	if err != nil {
		t.Fatalf("GenerateDockerfile failed: %v", err)
	}

	if !strings.Contains(output, "COPY go.mod ./") {
		t.Errorf("expected 'COPY go.mod ./' when HasGoSum=false, got:\n%s", output)
	}

	if strings.Contains(output, "go.sum") {
		t.Errorf("should not contain 'go.sum' when HasGoSum=false, got:\n%s", output)
	}
}

func TestGoGenerator_Generate_MultiStage(t *testing.T) {
	g := &GoGenerator{
		AppInfo: AppInfo{
			BinaryName: "myapp",
			GoVersion:  "1.22",
			HasGoSum:   true,
		},
	}

	output, err := g.GenerateDockerfile(".")
	if err != nil {
		t.Fatalf("GenerateDockerfile failed: %v", err)
	}

	// Check it's a multi-stage build
	if !strings.Contains(output, "AS builder") {
		t.Error("expected multi-stage build with 'AS builder'")
	}

	if !strings.Contains(output, "FROM alpine:latest") {
		t.Error("expected final stage 'FROM alpine:latest'")
	}

	if !strings.Contains(output, "COPY --from=builder") {
		t.Error("expected 'COPY --from=builder' in final stage")
	}
}

func TestGoGenerator_Generate_EntryPoint(t *testing.T) {
	g := &GoGenerator{
		AppInfo: AppInfo{
			BinaryName: "myapp",
			GoVersion:  "1.22",
			HasGoSum:   true,
		},
	}

	output, err := g.GenerateDockerfile(".")
	if err != nil {
		t.Fatalf("GenerateDockerfile failed: %v", err)
	}

	if !strings.Contains(output, `ENTRYPOINT ["/myapp"]`) {
		t.Errorf("expected ENTRYPOINT with binary name, got:\n%s", output)
	}
}
