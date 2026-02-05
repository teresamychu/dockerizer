// Dockerize CLI - generates Dockerfiles for Go Codebases
package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/teresamychu/dockerizer/generators"
)

func main() {
	dockerfileOutput := flag.String("dockerfile", "Dockerfile", "Dockerfile output filename")
	composeOutput := flag.String("compose", "docker-compose.yml", "docker-compose output filename")
	flag.Parse()

	path := "."
	if flag.NArg() > 0 {
		path = flag.Arg(0)
	}

	language := Detect(path)
	if language != LangGo {
		fmt.Println("Not a Golang application")
		os.Exit(1)
	}

	module, goversion, services := ParseGoMod(filepath.Join(path, "go.mod"))
	hasGoSum := fileExists(filepath.Join(path, "go.sum"))

	// Scan config files for port and env vars
	appConfig := ScanConfig(path)
	if appConfig.Port == "" {
		appConfig.Port = "8080"
	}

	gen := generators.GoGenerator{
		AppInfo: generators.AppInfo{
			BinaryName: module,
			GoVersion:  goversion,
			HasGoSum:   hasGoSum,
			Services:   services,
		},
		AppConfig: appConfig,
	}
	dockerFile, err := gen.GenerateDockerfile(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	outputPath := filepath.Join(path, *dockerfileOutput)
	err = os.WriteFile(outputPath, []byte(dockerFile), 0644)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error writing file: %v\n", err)
		os.Exit(1)
	}

	composeFile, err := gen.GenerateComposeFile(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	composeOutputPath := filepath.Join(path, *composeOutput)
	err = os.WriteFile(composeOutputPath, []byte(composeFile), 0644)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error writing file: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Generated %s for %s (Go %s)\n", outputPath, module, goversion)
}
